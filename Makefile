default: build

.PHONY: build, test

git_revision=git rev-parse --short HEAD
BUILD_TIME=`date +%FT%T%z`
BUILD_DATE=`date +%F`
GIT_REVISION=`${git_revision}`
GIT_COMMIT=`git rev-parse HEAD`
GIT_BRANCH=`git rev-parse --symbolic-full-name --abbrev-ref HEAD`
GIT_DIRTY=`test -z "$$(git status --porcelain)" && echo "clean" || echo "dirty"`
VERSION=`git describe --tag --abbrev=0 --exact-match HEAD 2> /dev/null || (echo 'GitVersion: Git tag not found, fallback to commit id' >&2; ${git_revision})`

METADATA_PATH=github.com/XSAM/go-hybrid/metadata
INJECT_VARIABLE=-X ${METADATA_PATH}.gitVersion=${VERSION} -X ${METADATA_PATH}.gitCommit=${GIT_COMMIT} -X ${METADATA_PATH}.gitBranch=${GIT_BRANCH} -X ${METADATA_PATH}.gitTreeState=${GIT_DIRTY} -X ${METADATA_PATH}.buildTime=${BUILD_TIME}

FLAGS=-trimpath -ldflags "-s ${INJECT_VARIABLE}"
DEBUG_FLAGS=-gcflags "all=-N -l" ${FLAGS}

export CGO_ENABLED=0

build:
	@echo "> Building example"
	go build -o bin/example ${FLAGS} _example/main.go

test:
	@echo "> Running unit test"
	@go test -gcflags=-l -coverprofile coverage.out ./...
	@echo "> Total coverage"
	@go tool cover -func coverage.out | grep total

install-tool:
	@echo "> Installing tools"
	GO111MODULE=off go get github.com/golangci/golangci-lint/cmd/golangci-lint

check:
	@echo "> Examining Go source code and reports suspicious constructs"
	golangci-lint run