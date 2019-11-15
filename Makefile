ROOT := $(shell realpath $(dir $(lastword $(MAKEFILE_LIST))))
APP := saas
APP_IMPORT_PATH := github.com/2hamed/saas
TIMESTAMP?=$(shell TZ="Asia/Tehran" date +'%Y-%m-%dT%H:%M:%S%z')
GIT_HEAD_REF := $(shell cat .git/HEAD | cut -d' ' -f2)
CI_COMMIT_REF_SLUG?=$(shell cat .git/HEAD | cut -d'/' -f3)
CI_COMMIT_SHORT_SHA?=$(shell cat .git/$(GIT_HEAD_REF) | head -c 8)
LDFLAGS := "-w -s \
	-X $(APP_IMPORT_PATH)/cmd.BuildDate=$(TIMESTAMP)\
	-X $(APP_IMPORT_PATH)/cmd.GitCommit=$(CI_COMMIT_SHORT_SHA) \
	-X $(APP_IMPORT_PATH)/cmd.GitRef=$(CI_COMMIT_REF_SLUG)"
DC_FILE="docker-compose.yml"

all: format lint-ci build-static-vendor

############################################################
# Build & Run
############################################################
build:
	go build -race .

build-static:
	CGO_ENABLED=0 go build -v -o $(APP) -a -installsuffix cgo -ldflags $(LDFLAGS) .

build-static-vendor:
	CGO_ENABLED=0 go build -mod vendor -v -o $(APP) -a -installsuffix cgo -ldflags $(LDFLAGS) .

docker:
	docker build \
          --build-arg http_proxy=$(PROXY) \
          --build-arg https_proxy=$(PROXY) \
          --build-arg=GIT_BRANCH=$(CI_COMMIT_REF_SLUG) \
          --build-arg=GIT_SHA=$(CI_COMMIT_SHORT_SHA) \
          --build-arg=GIT_TAG=$(CI_COMMIT_TAG) \
          --build-arg=BUILD_TIMESTAMP=$(TIMESTAMP) \
          -t $(APP):$(CI_COMMIT_REF_SLUG) .

install:
	cp $(APP) $(GOPATH)/bin

run:
	go run -race .

############################################################
# Test & Coverage
############################################################
test:
	go test -race -short -gcflags=-l -mod vendor -v ./...

test-integration:
	docker-compose -f ${ROOT}/docker-compose.testing.yml up

coverage:
	go test -race -short -gcflags=-l -mod vendor -v -coverprofile=.testCoverage.txt ./...
	GOFLAGS=-mod=vendor go tool cover -func=.testCoverage.txt

coverage-report: coverage
	GOFLAGS=-mod=vendor go tool cover -html=.testCoverage.txt -o testCoverageReport.html

############################################################
# Format & Lint
############################################################
check-goimport:
	which goimports || GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports

format: check-goimport
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R goimports -w R
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R gofmt -s -w R

check-golint:
	which golint || (go get -u golang.org/x/lint/golint)

lint: check-golint
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R golint -set_exit_status R

############################################################
# Development Environment
############################################################
up:
	docker-compose -f $(DC_FILE) up -d

.PHONY: build get install run watch start stop restart clean up ci-test