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
branchName := main
githubRepoLink := https://github.com/Eppo-exp/sdk-test-data.git
.PHONY: test-data
test-data:
	rm -rf $(testDataDir)
	git clone -b ${branchName} --depth 1 --single-branch ${githubRepoLink} ${testDataDir}

test: test-data
	go test -v ./...

lint:
	golangci-lint run

## profile-memory - Run test and generate memory profile
profile-memory: test-data
	@cd eppoclient && \
	{ \
	echo "Using OUTFILE_SUFFIX: $$OUTFILE_SUFFIX"; \
	go test -run Test_e2e -memprofile ../memprofile$$OUTFILE_SUFFIX.out ./...; \
	go tool pprof -text -nodecount=50 ../memprofile$$OUTFILE_SUFFIX.out > ../memprofile$$OUTFILE_SUFFIX.text; \
	}

## profile-memory-compare - Compare two memory profiles
## example: make profile-memory-compare BASE_FILE=memprofile1.out FEAT_FILE=memprofile2.out
profile-memory-compare:
	go tool pprof -base $$BASE_FILE -text $$FEAT_FILE
