package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"gopkg.in/asaskevich/govalidator.v6"
)

type DashboardSearchResult struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Uri   string `json:"uri"`
}

func httpGet(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving dashboards from URL %v: %v", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code returned from Grafana API (got: %d, expected: 200, msg:%s)", resp.StatusCode, resp.Status)
	}
	return resp, nil
}

func main() {
	grafanaURL := flag.String("grafanaURL", "", "The URL of the grafana server to be backed up")
	s3Bucket := flag.String("s3Bucket", "", "The name of the S3 bucket where the backup should be stored")
	useSharedConfig := flag.Bool("useSharedConfig", false, "Controls whether to use the ~/.aws shared config")
	flag.Parse()

	if *grafanaURL == "" || *s3Bucket == "" {
		flag.Usage()
		os.Exit(1)
	}

	if !govalidator.IsURL(*grafanaURL) {
		log.Fatalf("Invalid grafanaURL: %v", *grafanaURL)
	}

	// S3 client
	var sess *session.Session
	if *useSharedConfig {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	} else {
		sess = session.Must(session.NewSession())
	}
	uploader := s3manager.NewUploader(sess)

	getAllDashboardsURL := *grafanaURL + "/api/search"
	searchResults := make([]DashboardSearchResult, 0)
	resp, err := httpGet(getAllDashboardsURL)
	if err != nil {
		log.Fatalf("Error retrieving dashboard list from URL %v: %v", getAllDashboardsURL, err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&searchResults)
	if err != nil {
		log.Fatalf("Error retrieving decoding response body: %v", err)
	}

	getDashBaseURL := *grafanaURL + "/api/dashboards/"

	timeStr := time.Now().UTC().Format("2006-01-02T15:04:05-0700")
	backupDir := "grafana-backup_" + timeStr + "/dashboards/"
	for _, searchResult := range searchResults {
		// retrieve dashboard
		getDashURL := getDashBaseURL + searchResult.Uri
		resp, err := httpGet(getDashURL)
		if err != nil {
			log.Fatalf("Error retrieving dashboard from URL %v: %v", getDashURL, err)
		}

		// we must awkwardly prune out the 'id' key
		// as it interferes with restoring
		dashBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response from URL %v: %v", getDashURL, err)
		}
		dashJSON := make(map[string]interface{}, 0)
		err = json.Unmarshal(dashBytes, &dashJSON)
		if err != nil {
			log.Fatalf("Error decoding JSON from URL %v: %v", getDashURL, err)
		}
		if dashboardValue, ok1 := dashJSON["dashboard"]; ok1 {
			if dashboardValueAsMap, ok2 := dashboardValue.(map[string]interface{}); ok2 {
				dashboardValueAsMap["id"] = nil
			}
		}
		revisedDashBytes, err := json.Marshal(dashJSON)

		// upload to S3
		slug := strings.TrimPrefix(searchResult.Uri, "db/")
		filename := backupDir + slug
		ui := &s3manager.UploadInput{
			Bucket: s3Bucket,
			Key:    &filename,
			Body:   bytes.NewReader(revisedDashBytes),
		}

		_, err = uploader.Upload(ui)
		if err != nil {
			log.Fatalf("Error uploading dashboard json for dashboard from URL %v: %v", getDashURL, err)
		}
		resp.Body.Close()
	}

	log.Printf("Backup to directory %v in bucket %v completed at %v", backupDir, *s3Bucket, time.Now())
}
