TAGS=netgo osusergo sqlite_omit_load_extension
VERSION=$(shell git describe --tags --abbrev=0 | cut -c 2-)
COMMIT=$(shell git rev-parse --verify HEAD)
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LD_FLAGS=-s -w -linkmode external -extldflags "-static" -X main.BuildDate=$(DATE) -X main.BuildMode=prod -X main.BuildCommit=$(COMMIT) -X main.BuildVersion=$(VERSION)
BUILD_DIR=./build
PWD=$(shell pwd)
GOLANG_CROSS_VERSION=v1.22.0

license-dir:
	mkdir -p build/license || true

download-tools:
	go install golang.org/x/tools/cmd/goimports@v0.1.10
	go install github.com/gobuffalo/packr/v2/packr2@v2.7.1
	go install github.com/99designs/gqlgen@v0.17.44

generate-go:
	gqlgen

generate-js:
	(cd ui && yarn generate)

generate: generate-go generate-js

lint-go:
	go vet ./...
	goimports -l $(shell find . -type f -name '*.go' -not -path "./vendor/*")

lint-js:
	(cd ui && yarn format:check)
	(cd ui && yarn lint:check)

lint: lint-go lint-js

format-go:
	goimports -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

format-js:
	(cd ui && yarn format)

format: format-go format-js

test-go:
	go test --race -coverprofile=coverage.txt -covermode=atomic ./...

test-js:
	(cd ui && CI=true yarn test)

test: test-go test-js

install-go:
	go mod download

install-js:
	(cd ui && yarn)

build-js:
	(cd ui && yarn build)

packr:
	packr2

packr-clean:
	packr2 clean

pre-build: build-js packr

build-bin-local: pre-build
	CGO_ENABLED=1 go build -a -ldflags '${LD_FLAGS}' -tags '${TAGS}' -o ${BUILD_DIR}/traggo-server

.PHONY: release
release:
	docker build -t traggo:build -f docker/Dockerfile.build docker
	docker run \
		--rm \
		-e CGO_ENABLED=1 \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $$PWD:/work \
		-w /work \
		traggo:build \
		release --skip-validate


install: install-go install-js

.PHONY: build
