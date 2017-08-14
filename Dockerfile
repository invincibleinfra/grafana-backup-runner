FROM scratch

# Need this for https from within container
# https://blog.codeship.com/building-minimal-docker-containers-for-go-applications/
COPY ca-certificates.crt  /etc/ssl/certs/
COPY grafana-backup-runner /grafana-backup-runner

ENTRYPOINT ["/grafana-backup-runner"]
