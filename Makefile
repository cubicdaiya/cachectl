
cachectl: *.go
	go build -o cachectl

fmt:
	go fmt ./...

clean:
	rm -rf cachectl
