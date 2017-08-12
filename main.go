package main

import "flag"
import "gopkg.in/asaskevich/govalidator.v6"
import "log"
import "os"
import "strings"

func main() {
	grafanaURLRaw := flag.String("grafanaURL", "", "The URL of the grafana server to be backed up")
	s3BucketURLRaw := flag.String("s3BucketURL", "", "The URL of the S3 bucket where the backup should be stored")
	flag.Parse()

	if *grafanaURLRaw == "" || *s3BucketURLRaw == "" {
		flag.Usage()
		os.Exit(1)
	}

	if !govalidator.IsURL(*grafanaURLRaw) {
		log.Fatalf("Invalid grafanaURL: %v", *grafanaURLRaw)
	}

	// validator cannot handle 's3://' type URLs
	tmpS3BucketURL := strings.Replace(*s3BucketURLRaw, "s3:", "http:", 1)
	if !govalidator.IsURL(tmpS3BucketURL) {
		log.Fatalf("Invalid s3BucketURL: %v", *s3BucketURLRaw)
	}
}
