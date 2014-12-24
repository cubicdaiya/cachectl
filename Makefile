
all: bin/cachectl bin/cachectld

bin/cachectl:
	go build -o bin/cachectl cachectl.go

bin/cachectld:
	go build -o bin/cachectld cachectld.go

fmt:
	go fmt ./...

clean:
	rm -rf bin/cachectl bin/cachectld
