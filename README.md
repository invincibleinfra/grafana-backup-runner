# Description

This program backs up all grafana dashboards from a given server to S3.
The user provides the Grafana server URL and the destination S3 Bucket.

# Setup

To build the executable, use `make build`. To build a docker image from this
executable, use `make image`. Note: `make image` does not automatically build the executable.

# Usage

The easiest way to test this program (on a computer with AWS CLI access configured in `~/.aws`) is to run the command:

## Basic Testing

```
./grafana-backup-runner -grafanaURL ${GRAFANA_URL} -s3Bucket ${S3_BUCKET_NAME}
```

## Running on Kubernetes

To use this program on a production Kubernetes cluster, the user must define an
appropriate set of secrets with the relevant AWS credentials. It is recommended to
place the secrets in an env file like so:

```
AWS_ACCESS_KEY_ID=SUPER_SECRET
AWS_SECRET_ACCESS_KEY=SUPER_SECRET
AWS_REGION=us-east-1

GRAFANA_BACKUP_GRAFANA_URL=YOUR_BACKUP_URL
GRAFANA_BACKUP_S3_BUCKET_NAME=YOUR_BUCKET_NAME
```

This env file can be used to create the required Kubernetes secret via the command:


```
kubectl create secret generic grafana-backup-runner-config --from-env-file grafana-backup-runner-config.env
```

To run a single backup operation on a Kubernetes cluster, use the following manifest:

```
kubectl create -f manifests/grafana-backup-runner-job.yaml
```

Once this job completes successfully, you should see a directory with a name of the form "grafana_backup-$TIMESTAMP" appear in your S3 bucket.

## Recurring Backups

In order to use Kubernetes' `CronJob` support to implement recurring backups, the `CronJob` alpha feature must be enabled on the Kubernetes API server. To test
if your API server supports `CronJob`, run `kubectl api-versions`. If you do not see `batch/v2alpha1` in the output, you will need to enable it yourself. For a cluster deployed using Kubespray, this can be accomplished by editing the appropriate manifest file:

```
sudo vi /etc/kubernetes/manifests/kube-apiserver.manifest
```

Add the following line to the `apiserver` invocation:

```
- --runtime-config=batch/v2alpha1
```

You may need to restart the API server for this setting take effect.

Once the `CronJob`-compatible API version is successfully enabled, you can deploy a recurring backup
`CronJob` using our provided manifest:

```
kubectl create -f manifests/grafana-backup-runner-cronjob.yaml
```
