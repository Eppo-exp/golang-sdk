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

## test-data
testDataDir := eppoclient/test-data/
tempDir := ${testDataDir}temp/
gitDataDir := ${tempDir}sdk-test-data/
branchName := main
githubRepoLink := https://github.com/Eppo-exp/sdk-test-data.git
.PHONY: test-data
test-data:
	rm -rf $(testDataDir)
	mkdir -p $(tempDir)
	git clone -b ${branchName} --depth 1 --single-branch ${githubRepoLink} ${gitDataDir}
	cp ${gitDataDir}rac-experiments-v3.json ${testDataDir}
	cp -r ${gitDataDir}assignment-v2 ${testDataDir}
	rm -rf ${tempDir}

test: test-data
	go test ./...

lint:
	golangci-lint run
