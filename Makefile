VERSION=0.1.0
BUILD=`git rev-parse --short=8 HEAD`
.PHONY: all fmt static precommit

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"


all:
	go build ${LDFLAGS}

fmt:
	gofumpt -l -w .

static:
	go vet ./...
	staticcheck ./...

precommit: fmt static