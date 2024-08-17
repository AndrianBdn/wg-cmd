VERSION=0.1.6
BUILD=`git rev-parse --short=8 HEAD`
.PHONY: all fmt static precommit arm64 amd64 fmt static test release release_dir

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"
BINARY=wg-cmd

all:
	go build ${LDFLAGS}


arm64: export GOOS=linux
arm64: export GOARCH=arm64
arm64: release_dir
	go build ${LDFLAGS} -o release/${BINARY}-${VERSION}-${GOOS}-${GOARCH}


amd64: export GOOS=linux
amd64: export GOARCH=amd64
amd64: release_dir
	go build ${LDFLAGS} -o release/${BINARY}-${VERSION}-${GOOS}-${GOARCH}


fmt:
	gofumpt -l -w .

static:
	go vet ./...
	staticcheck ./...

test:
	go test ./...

release_dir:
	mkdir -p release

release: arm64 amd64

precommit: fmt static test