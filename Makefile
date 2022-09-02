BINARY_NAME=eppo-golang-sdk
default: help

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin ./eppoclient

clean:
	go clean
	rm ${BINARY_NAME}-darwin

help: Makefile
	@echo "usage: make <target>"
	@sed -n 's/^##//p' $<

testDataDir := eppoclient/test-data/
.PHONY: test-data
test-data:
	rm -rf $(testDataDir)
	mkdir -p $(testDataDir)
	gsutil cp gs://sdk-test-data/rac-experiments.json $(testDataDir)
	gsutil cp -r gs://sdk-test-data/assignment-v2 $(testDataDir)

test: test-data
	go test ./...

lint:
	golangci-lint run
