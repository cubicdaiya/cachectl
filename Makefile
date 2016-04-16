VERSION=0.3.4
TARGETS_NOVENDOR=$(shell glide novendor)

all: bin/cachectl bin/cachectld

build-cross: cmd/cachectl/cachectl.go cmd/cachectld/cachectld.go cachectl/*.go
	GO15VENDOREXPERIMENT=1 GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64/cachectl-${VERSION}/cachectl cmd/cachectl/cachectl.go
	GO15VENDOREXPERIMENT=1 GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64/cachectl-${VERSION}/cachectld cmd/cachectld/cachectld.go

dist: build-cross
	cd bin/linux/amd64 && tar zcvf cachectl-linux-amd64-${VERSION}.tar.gz cachectl-${VERSION}

bundle:
	glide install

bin/cachectl: cmd/cachectl/cachectl.go cachectl/*.go
	GO15VENDOREXPERIMENT=1 go build $(GOFLAGS) -o bin/cachectl github.com/cubicdaiya/cachectl/cmd/cachectl

bin/cachectld: cmd/cachectld/cachectld.go cachectl/*.go
	GO15VENDOREXPERIMENT=1 go build $(GOFLAGS) -o bin/cachectld github.com/cubicdaiya/cachectl/cmd/cachectld

fmt:
	@echo $(TARGETS_NOVENDOR) | xargs go fmt

clean:
	rm -rf bin/*
