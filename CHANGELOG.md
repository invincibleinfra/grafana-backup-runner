## v0.0.2 / 2017-08-17

This release improves usability by integrating with a [restore utility](https://github.com/invincibleinfra/grafana-backup-restore)
and by making the dashboard filenames in the S3 backups human-readable.

## v0.0.1 / 2017-08-15

Very first alpha release. These fundamental operations have been tested:

* Ad-hoc backup using the executable, with credentials provided via shared AWS config or environment variables.

* CronJob backup using Kubernetes secrets to store the credentials

There is no restore functionality in this release.
