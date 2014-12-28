
all: bin/cachectl bin/cachectld

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
	rm -rf bin/cachectl bin/cachectld
