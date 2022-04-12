VERSION=0.3.9

all: bin/cachectl bin/cachectld

build-cross: cmd/cachectl/cachectl.go cmd/cachectld/cachectld.go cachectl/*.go
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o bin/linux/amd64/cachectl-${VERSION}/cachectl cmd/cachectl/cachectl.go
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o bin/linux/amd64/cachectl-${VERSION}/cachectld cmd/cachectld/cachectld.go

bin/cachectl: cmd/cachectl/cachectl.go cachectl/*.go
	GO111MODULE=on go build $(GOFLAGS) -o bin/cachectl cmd/cachectl/cachectl.go

bin/cachectld: cmd/cachectld/cachectld.go cachectl/*.go
	GO111MODULE=on go build $(GOFLAGS) -o bin/cachectld cmd/cachectld/cachectld.go

check:
	go test -v ./cachectl/

fmt:
	GO111MODULE=on go fmt ./...

clean:
	rm -rf bin/*
