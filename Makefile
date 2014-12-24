
bin/cachectl: *.go
	go build -o bin/cachectl

fmt:
	go fmt ./...

clean:
	rm -rf bin/cachectl
