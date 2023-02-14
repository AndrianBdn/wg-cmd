VERSION=0.1.0
BUILD=`git rev-parse --short=8 HEAD`
.PHONY: all fmt static precommit arm64 amd64

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"
BINARY=wg-cmd

all:
	go build ${LDFLAGS}


arm64: export GOOS=linux
arm64: export GOARCH=arm64
arm64:
	go build ${LDFLAGS} -o ${BINARY}-${GOOS}-${GOARCH}


amd64: export GOOS=linux
amd64: export GOARCH=amd64
amd64:
	go build ${LDFLAGS} -o ${BINARY}-${GOOS}-${GOARCH}


fmt:
	gofumpt -l -w .

static:
	go vet ./...
	staticcheck ./...

precommit: fmt static