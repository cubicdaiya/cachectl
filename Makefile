VERSION=0.3.0

all: bin/cachectl bin/cachectld

build-cross: cachectl.go cachectld.go cachectl/*.go
	GOOS=linux GOARCH=amd64 gom build -o bin/linux/amd64/cachectl-${VERSION}/cachectl cachectl.go
	GOOS=linux GOARCH=amd64 gom build -o bin/linux/amd64/cachectl-${VERSION}/cachectld cachectld.go

dist: build-cross
	cd bin/linux/amd64 && tar zcvf cachectl-linux-amd64-${VERSION}.tar.gz cachectl-${VERSION}

gom:
	go get -u github.com/mattn/gom

bundle:
	gom install

bin/cachectl: cachectl.go cachectl/*.go
	gom build -o bin/cachectl cachectl.go

bin/cachectld: cachectld.go cachectl/*.go
	gom build -o bin/cachectld cachectld.go

fmt:
	go fmt ./...

clean:
	rm -rf bin/*
