TAGS=netgo osusergo sqlite_omit_load_extension
VERSION=$(shell git describe --tags --abbrev=0 | cut -c 2-)
COMMIT=$(shell git rev-parse --verify HEAD)
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LD_FLAGS=-s -w -linkmode external -extldflags "-static" -X main.BuildDate=$(DATE) -X main.BuildMode=prod -X main.BuildCommit=$(COMMIT) -X main.BuildVersion=$(VERSION)
BUILD_DIR=./build
BUILD_LICENSE=${BUILD_DIR}/license/3RD_PARTY_LICENSES
UI_BUILD_LICENSE=${BUILD_DIR}/license/UI_3RD_PARTY_LICENSES
PWD=$(shell pwd)
GO_VERSION=1.13.1
DOCKER_BUILD_IMAGE=traggo/build
DOCKER_WORKDIR=/proj
GOPATH_VOLUME=-v "`go env GOPATH`/pkg/mod/:/go/pkg/mod"
WORKDIR_VOLUME=-v "${PWD}/:${DOCKER_WORKDIR}"
DOCKER_GO_BUILD=go build -a -installsuffix cgo -ldflags '${LD_FLAGS}' -tags '${TAGS}'
DOCKER_RUN=docker run --rm ${WORKDIR_VOLUME} ${GOPATH_VOLUME} -w ${DOCKER_WORKDIR}
NEW_IMAGE_NAME=traggo/server
DOCKER_MANIFEST=DOCKER_CLI_EXPERIMENTAL=enabled docker manifest
BIN_PREFIX=traggo-server

license-dir:
	mkdir -p build/license || true

download-tools:
	go install golang.org/x/tools/cmd/goimports@v0.1.10
	go install github.com/gobuffalo/packr/v2/packr2@v2.7.1


generate-go:
	go run hack/gqlgen/gqlgen.go

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

build-bin-linux-amd64: pre-build
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-linux-amd64   ${DOCKER_GO_BUILD} -o ${BUILD_DIR}/${BIN_PREFIX}-${VERSION}-linux-amd64 ${DOCKER_WORKDIR}

build-docker-linux-amd64:
	cp ${BUILD_DIR}/${BIN_PREFIX}-${VERSION}-linux-amd64 docker/traggo && docker build -t ${NEW_IMAGE_NAME}:amd64-latest -t ${NEW_IMAGE_NAME}:amd64-${VERSION} docker/ && rm docker/traggo

build-bin-linux-386: pre-build
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-linux-386     ${DOCKER_GO_BUILD} -o ${BUILD_DIR}/${BIN_PREFIX}-${VERSION}-linux-386 ${DOCKER_WORKDIR}

build-docker-linux-386:
	cp ${BUILD_DIR}/${BIN_PREFIX}-${VERSION}-linux-386 docker/traggo && docker build -t ${NEW_IMAGE_NAME}:386-latest -t ${NEW_IMAGE_NAME}:386-${VERSION} docker/ && rm docker/traggo

build-bin-linux-arm-7: pre-build
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-linux-arm-7   ${DOCKER_GO_BUILD} -o ${BUILD_DIR}/${BIN_PREFIX}-${VERSION}-linux-arm-7 ${DOCKER_WORKDIR}

build-docker-linux-arm-7:
	cp ${BUILD_DIR}/${BIN_PREFIX}-${VERSION}-linux-arm-7 docker/traggo && docker build -t ${NEW_IMAGE_NAME}:arm-7-latest -t ${NEW_IMAGE_NAME}:arm-7-${VERSION} docker/ && rm docker/traggo

build-bin-linux-arm64: pre-build
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-linux-arm64   ${DOCKER_GO_BUILD} -o ${BUILD_DIR}/${BIN_PREFIX}-${VERSION}-linux-arm64 ${DOCKER_WORKDIR}

build-docker-linux-arm64:
	cp ${BUILD_DIR}/${BIN_PREFIX}-${VERSION}-linux-arm64 docker/traggo && docker build -t ${NEW_IMAGE_NAME}:arm64-latest -t ${NEW_IMAGE_NAME}:arm64-${VERSION} docker/ && rm docker/traggo

build-bin-windows-amd64: pre-build
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-windows-amd64 ${DOCKER_GO_BUILD} -o ${BUILD_DIR}/${BIN_PREFIX}-${VERSION}-windows-amd64.exe ${DOCKER_WORKDIR}

build-bin-windows-386: pre-build
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-windows-386   ${DOCKER_GO_BUILD} -o ${BUILD_DIR}/${BIN_PREFIX}-${VERSION}-windows-386.exe ${DOCKER_WORKDIR}

build-bin-linux-ppc64le: pre-build
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-linux-ppc64le   ${DOCKER_GO_BUILD} -o ${BUILD_DIR}/${BIN_PREFIX}-${VERSION}-linux-ppc64le ${DOCKER_WORKDIR}

build-docker-linux-ppc64le:
	cp ${BUILD_DIR}/${BIN_PREFIX}-${VERSION}-linux-ppc64le docker/traggo && docker build -t ${NEW_IMAGE_NAME}:ppc64le-latest -t ${NEW_IMAGE_NAME}:ppc64le-${VERSION} docker/ && rm docker/traggo

build-bin-linux-s390x: pre-build
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-linux-s390x   ${DOCKER_GO_BUILD} -o ${BUILD_DIR}/${BIN_PREFIX}-${VERSION}-linux-s390x ${DOCKER_WORKDIR}

build-docker-linux-s390x:
	cp ${BUILD_DIR}/${BIN_PREFIX}-${VERSION}-linux-s390x docker/traggo && docker build -t ${NEW_IMAGE_NAME}:s390x-latest -t ${NEW_IMAGE_NAME}:s390x-${VERSION} docker/ && rm docker/traggo

build-bin:    build-bin-linux-amd64    build-bin-linux-386    build-bin-linux-arm-7    build-bin-linux-arm64    build-bin-linux-ppc64le    build-bin-linux-s390x build-bin-windows-amd64 build-bin-windows-386 

build-docker: build-docker-linux-amd64 build-docker-linux-386 build-docker-linux-arm-7 build-docker-linux-arm64 build-docker-linux-ppc64le build-docker-linux-s390x

fix-build-owner:
	sudo chown -R $(shell id -u):$(shell id -g) ${BUILD_DIR}

docker-login-ci:
	docker login -u "$$DOCKER_USER" -p "$$DOCKER_PASS";

docker-push:
	docker push ${NEW_IMAGE_NAME}

docker-push-manifest:
	${DOCKER_MANIFEST} create "${NEW_IMAGE_NAME}:latest"     "${NEW_IMAGE_NAME}:amd64-latest"     "${NEW_IMAGE_NAME}:386-latest"     "${NEW_IMAGE_NAME}:arm-7-latest"     "${NEW_IMAGE_NAME}:arm64-latest"     "${NEW_IMAGE_NAME}:ppc64le-latest"     "${NEW_IMAGE_NAME}:s390x-latest"
	${DOCKER_MANIFEST} create "${NEW_IMAGE_NAME}:${VERSION}" "${NEW_IMAGE_NAME}:amd64-${VERSION}" "${NEW_IMAGE_NAME}:386-${VERSION}" "${NEW_IMAGE_NAME}:arm-7-${VERSION}" "${NEW_IMAGE_NAME}:arm64-${VERSION}" "${NEW_IMAGE_NAME}:ppc64le-${VERSION}" "${NEW_IMAGE_NAME}:s390x-${VERSION}"
	${DOCKER_MANIFEST} annotate "${NEW_IMAGE_NAME}:latest"     "${NEW_IMAGE_NAME}:amd64-latest"       --os=linux --arch=amd64
	${DOCKER_MANIFEST} annotate "${NEW_IMAGE_NAME}:${VERSION}" "${NEW_IMAGE_NAME}:amd64-${VERSION}"   --os=linux --arch=amd64
	${DOCKER_MANIFEST} annotate "${NEW_IMAGE_NAME}:latest"     "${NEW_IMAGE_NAME}:386-latest"         --os=linux --arch=386
	${DOCKER_MANIFEST} annotate "${NEW_IMAGE_NAME}:${VERSION}" "${NEW_IMAGE_NAME}:386-${VERSION}"     --os=linux --arch=386
	${DOCKER_MANIFEST} annotate "${NEW_IMAGE_NAME}:latest"     "${NEW_IMAGE_NAME}:arm-7-latest"       --os=linux --arch=arm --variant=v7
	${DOCKER_MANIFEST} annotate "${NEW_IMAGE_NAME}:${VERSION}" "${NEW_IMAGE_NAME}:arm-7-${VERSION}"   --os=linux --arch=arm --variant=v7
	${DOCKER_MANIFEST} annotate "${NEW_IMAGE_NAME}:latest"     "${NEW_IMAGE_NAME}:arm64-latest"       --os=linux --arch=arm64
	${DOCKER_MANIFEST} annotate "${NEW_IMAGE_NAME}:${VERSION}" "${NEW_IMAGE_NAME}:arm64-${VERSION}"   --os=linux --arch=arm64
	${DOCKER_MANIFEST} annotate "${NEW_IMAGE_NAME}:latest"     "${NEW_IMAGE_NAME}:ppc64le-latest"     --os=linux --arch=ppc64le
	${DOCKER_MANIFEST} annotate "${NEW_IMAGE_NAME}:${VERSION}" "${NEW_IMAGE_NAME}:ppc64le-${VERSION}" --os=linux --arch=ppc64le
	${DOCKER_MANIFEST} annotate "${NEW_IMAGE_NAME}:latest"     "${NEW_IMAGE_NAME}:s390x-latest"       --os=linux --arch=s390x
	${DOCKER_MANIFEST} annotate "${NEW_IMAGE_NAME}:${VERSION}" "${NEW_IMAGE_NAME}:s390x-${VERSION}"   --os=linux --arch=s390x
	${DOCKER_MANIFEST} push "${NEW_IMAGE_NAME}:${VERSION}"
	${DOCKER_MANIFEST} push "${NEW_IMAGE_NAME}:latest"

licenses-ui: license-dir
	(cd ui && yarn -s licenses generate-disclaimer --prod > ../${UI_BUILD_LICENSE})

licenses-go: license-dir
	go mod vendor
	echo "THE FOLLOWING SETS FORTH ATTRIBUTION NOTICES FOR THIRD PARTY SOFTWARE THAT MAY BE CONTAINED IN PORTIONS OF THE TRAGGO PRODUCT" > ${BUILD_LICENSE}
	echo >> ${BUILD_LICENSE}
	echo ------- >> ${BUILD_LICENSE}
	echo >> ${BUILD_LICENSE}
	(cd vendor && find . -type f \( -iname "LICENSE*" -o -iname "NOTICE*" \) -exec echo The following software may be included in this product {} \; -exec echo  \; -exec cat {} \; -exec echo \; -exec echo -------- \; -exec echo \;) >> ${BUILD_LICENSE}

package-zip: licenses-ui licenses-go
	find build/* -maxdepth 0 -type f -exec zip -9 -j {}.zip {} build/license/3RD_PARTY_LICENSES build/license/UI_3RD_PARTY_LICENSES LICENSE \;

build-compress:
	find build/* -maxdepth 0 -type f -exec upx {} \;

install: install-go install-js

.PHONY: build
