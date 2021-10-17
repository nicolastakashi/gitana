RELEASE_VERSION ?=$(shell cat VERSION)
RELEASE=1
REVISION ?= $(shell git rev-parse HEAD)
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
BINARY_FOLDER=bin
BINARY_NAME=gitana
ARTIFACT_NAME=ntakashi/$(BINARY_NAME)
GOCMD=go
GOMAIN=main.go
GOBUILD=$(GOCMD) build
ENVVARS=GOOS=linux CGO_ENABLED=0
GOCMD=go
LDFLAGS=-w -extldflags "-static" \
		-X github.com/prometheus/common/version.Version=$(RELEASE_VERSION) \
		-X github.com/prometheus/common/version.Revision=$(REVISION) \
		-X github.com/prometheus/common/version.Branch=$(BRANCH) \
		-X github.com/prometheus/common/version.BuildUser=$(shell whoami) \
		-X "github.com/prometheus/common/version.BuildDate=$(shell date -u)"

docker-build:
	@DOCKER_BUILDKIT=1 docker build -t ${ARTIFACT_NAME}:${RELEASE_VERSION} -f ./Dockerfile --progress=plain .

docker-push:
	@DOCKER_BUILDKIT=1 docker push $(ARTIFACT_NAME):${RELEASE_VERSION}
		
deps:
	$(ENVVARS) $(GOCMD) mod download

fmt:
	$(ENVVARS) $(GOCMD) fmt -x ./...

vet:
	$(ENVVARS) $(GOCMD) vet ./...

build:
	$(ENVVARS) $(GOCMD) build -ldflags '$(LDFLAGS)' -o $(BINARY_FOLDER)/$(BINARY_NAME) -v $(GOMAIN)

all: build fmt vet deps