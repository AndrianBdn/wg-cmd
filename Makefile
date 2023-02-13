VERSION=0.1.0
# UTC timestamp in ISO 8601 format
BUILD_TIME=`date -u +"%Y-%m-%dT%H:%M:%SZ"`
.PHONY: all fmt static precommit

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"


all:
	go build ${LDFLAGS}

fmt:
	gofumpt -l -w .

static:
	go vet ./...
	staticcheck ./...

precommit: fmt static