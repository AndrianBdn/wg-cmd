VERSION=0.1.10
BUILD=`git rev-parse --short=8 HEAD`
.PHONY: all fmt static precommit arm64 amd64 fmt static test release release_dir

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"
BINARY=wg-cmd

all:
	go build ${LDFLAGS}


arm64: export GOOS=linux
arm64: export GOARCH=arm64
arm64: release_dir
	go build ${LDFLAGS} -o release/${BINARY}-${GOOS}-${GOARCH}


amd64: export GOOS=linux
amd64: export GOARCH=amd64
amd64: release_dir
	go build ${LDFLAGS} -o release/${BINARY}-${GOOS}-${GOARCH}


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
	cd release && shasum -a 256 ${BINARY}-linux-amd64 ${BINARY}-linux-arm64 > SHA256SUMS

precommit: fmt static test