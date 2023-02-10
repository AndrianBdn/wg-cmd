VERSION=0.0

all:
	go build -ldflags "-X main.version=0.10"
