# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

SHELL := /bin/bash # Use bash syntax

# Build variables
BUILD_DIR ?= bin
VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)
COMMIT_HASH ?= $(shell git rev-parse HEAD 2>/dev/null)
BUILD_DATE ?= $(shell date +%FT%T%z)
BUILD_BY ?= $(shell git config user.email)
LDFLAGS += -X main.tag=${VERSION} -X main.commit=${COMMIT_HASH} -X main.buildDate=${BUILD_DATE} -X main.builtBy=${BUILD_BY}

# Project variables
DOCKER_IMAGE = adrienaury/owl
DOCKER_TAG ?= $(shell echo -n ${VERSION} | sed -e 's/[^A-Za-z0-9_\\.-]/_/g')
RELEASE := $(shell [[ $(VERSION) =~ ^[0-9]*.[0-9]*.[0-9]*$$ ]] && echo 1 || echo 0 )
MAJOR := $(shell echo $(VERSION) | cut -f1 -d.)
MINOR := $(shell echo $(VERSION) | cut -f2 -d.)
PATCH := $(shell echo $(VERSION) | cut -f3 -d. | cut -f1 -d-)

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

.PHONY: info
info: ## Prints build informations
	@echo COMMIT_HASH=$(COMMIT_HASH)
	@echo VERSION=$(VERSION)
	@echo RELEASE=$(RELEASE)
ifeq (${RELEASE}, 1)
	@echo MAJOR=$(MAJOR)
	@echo MINOR=$(MINOR)
	@echo PATCH=$(PATCH)
endif
	@echo DOCKER_IMAGE=$(DOCKER_IMAGE)
	@echo DOCKER_TAG=$(DOCKER_TAG)
	@echo BUILD_BY=$(BUILD_BY)

.PHONY: clean
clean: ## Clean builds
	rm -rf ${BUILD_DIR}/

.PHONY: mkdir
mkdir:
	mkdir -p ${BUILD_DIR}

.PHONY: tidy
tidy: ## Add missing and remove unused modules
	GO111MODULE=on go mod tidy

.PHONY: lint
lint: ## Examines Go source code and reports suspicious constructs
ifeq (, $(shell which golangci-lint))
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(go env GOPATH)/bin v1.21.0
endif
	golangci-lint run -E misspell -E gocyclo -E gosec -E unparam -E goimports -E nakedret -E gocritic

#.PHONY: generate
#generate: ## Update generated code
#	GO111MODULE=on go generate configs/generate.go

.PHONY: build-%
build-%: mkdir #generate
	GO111MODULE=on go build ${GOARGS} -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/$* ./cmd/$*

.PHONY: build
build: $(patsubst cmd/%,build-%,$(wildcard cmd/*)) ## Build all binaries

.PHONY: test
test: mkdir ## Run all tests with coverage
	GO111MODULE=on go test -coverprofile=${BUILD_DIR}/coverage.txt -covermode=atomic ./...

.PHONY: run-%
run-%: build-%
	${BUILD_DIR}/$* ${ARGS}

.PHONY: run
run: $(patsubst cmd/%,run-%,$(wildcard cmd/*)) ## Build and execute a binary

.PHONY: release-%
release-%: mkdir
	GO111MODULE=on go build ${GOARGS} -ldflags "-w -s ${LDFLAGS}" -o ${BUILD_DIR}/$* ./cmd/$*

.PHONY: release
release: clean info lint $(patsubst cmd/%,release-%,$(wildcard cmd/*)) ## Build all binaries for production

.PHONY: docker
docker: info ## Build docker image locally
	docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} --build-arg VERSION=${VERSION} --build-arg BUILD_BY=${BUILD_BY} .
ifeq (${RELEASE}, 1)
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:${MAJOR}.${MINOR}
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:${MAJOR}
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
endif

.PHONY: push
push: docker ## Push docker image on DockerHub
	docker push ${DOCKER_IMAGE}:${DOCKER_TAG}
ifeq (${RELEASE}, 1)
	docker push ${DOCKER_IMAGE}:${MAJOR}.${MINOR}
	docker push ${DOCKER_IMAGE}:${MAJOR}
	docker push ${DOCKER_IMAGE}:latest
endif

.PHONY: license
license: mkdir docker ## Scan dependencies and licenses
	docker create --name owl-license ${DOCKER_IMAGE}:${DOCKER_TAG}
	docker cp owl-license:/owl - > ${BUILD_DIR}/owl.tar
	docker rm -v owl-license
	mkdir -p ${BUILD_DIR}/owl-license
	tar xvf ${BUILD_DIR}/owl.tar -C ${BUILD_DIR}/owl-license
	golicense ${BUILD_DIR}/owl-license/owl
