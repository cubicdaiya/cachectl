VERSION=0.3.1

all: bin/cachectl bin/cachectld

build-cross: cmd/cachectl/cachectl.go cmd/cachectld/cachectld.go cachectl/*.go
	GOOS=linux GOARCH=amd64 gom build -o bin/linux/amd64/cachectl-${VERSION}/cachectl cmd/cachectl/cachectl.go
	GOOS=linux GOARCH=amd64 gom build -o bin/linux/amd64/cachectl-${VERSION}/cachectld cmd/cachectld/cachectld.go

dist: build-cross
	cd bin/linux/amd64 && tar zcvf cachectl-linux-amd64-${VERSION}.tar.gz cachectl-${VERSION}

gom:
	go get -u github.com/mattn/gom

bundle:
	gom install

bin/cachectl: cmd/cachectl/cachectl.go cachectl/*.go
	gom build $(GOFLAGS) -o bin/cachectl github.com/cubicdaiya/cachectl/cmd/cachectl

bin/cachectld: cmd/cachectld/cachectld.go cachectl/*.go
	gom build $(GOFLAGS) -o bin/cachectld github.com/cubicdaiya/cachectl/cmd/cachectld

fmt:
	go fmt ./...

clean:
	rm -rf bin/*
