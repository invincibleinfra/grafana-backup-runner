package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/asaskevich/govalidator.v6"
	"log"
	"net/http"
	"os"
	"strings"
)

type DashboardSearchResult struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Uri   string `json:"uri"`
}

func httpGet(url string, result interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error retrieving dashboards from URL %v: %v", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status code returned from Grafana API (got: %d, expected: 200, msg:%s)", resp.StatusCode, resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return fmt.Errorf("Error retrieving dashboards from URL %v: %v", url, err)
	}
	return nil
}

func main() {
	grafanaURL := flag.String("grafanaURL", "", "The URL of the grafana server to be backed up")
	s3BucketURL := flag.String("s3BucketURL", "", "The URL of the S3 bucket where the backup should be stored")
	flag.Parse()

	if *grafanaURL == "" || *s3BucketURL == "" {
		flag.Usage()
		os.Exit(1)
	}

	if !govalidator.IsURL(*grafanaURL) {
		log.Fatalf("Invalid grafanaURL: %v", *grafanaURL)
	}

	// validator cannot handle 's3://' type URLs
	tmpS3BucketURL := strings.Replace(*s3BucketURL, "s3:", "http:", 1)
	if !govalidator.IsURL(tmpS3BucketURL) {
		log.Fatalf("Invalid s3BucketURL: %v", *s3BucketURL)
	}

	getAllDashboardsURL := *grafanaURL + "/api/search"
	searchResults := make([]DashboardSearchResult, 0)
	err := httpGet(getAllDashboardsURL, &searchResults)
	if err != nil {
		log.Fatalf("Error in http request: %v", err)
	}

	getDashBaseURL := *grafanaURL + "/api/dashboards/"
	for _, searchResult := range searchResults {
		getDashURL := getDashBaseURL + searchResult.Uri
		dashboardJSON := map[string]interface{}{}
		err := httpGet(getDashURL, &dashboardJSON)
		if err != nil {
			log.Fatalf("Error in http request: %v", err)
		}
		log.Printf("dashboard %v: %v", searchResult.Id, dashboardJSON)
	}
}
