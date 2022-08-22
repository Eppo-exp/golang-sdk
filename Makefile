BINARY_NAME=eppo-golang-sdk
default: help

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin doc.go

clean:
	go clean
	rm ${BINARY_NAME}-darwin

help: Makefile
	@echo "usage: make <target>"
	@sed -n 's/^##//p' $<

test:
	go test ./...

lint:
	golangci-lint run --enable-all
