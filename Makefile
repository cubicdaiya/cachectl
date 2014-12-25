
all: bin/cachectl bin/cachectld

bundle:
	gom install

bin/cachectl:
	gom build -o bin/cachectl cachectl.go

bin/cachectld:
	gom build -o bin/cachectld cachectld.go

fmt:
	go fmt ./...

clean:
	rm -rf bin/cachectl bin/cachectld
