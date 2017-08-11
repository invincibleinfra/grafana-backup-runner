REPO_URL = github.com/invincibleinfra/grafana-backup-runner

# Build using golang docker container for reproducibility
GOLANG_IMAGE=golang:1.8.3-jessie
CONTAINER_GOPATH=/go
CONTAINER_SOURCE_DIR=$(CONTAINER_GOPATH)/src/$(REPO_URL)

build:
	docker run \
		-v $(CURDIR):$(CONTAINER_SOURCE_DIR) \
		-w $(CONTAINER_SOURCE_DIR) \
		--rm $(GOLANG_IMAGE) \
		go build -o grafana-backup-runner
