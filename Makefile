TAGS=netgo osusergo sqlite_omit_load_extension
VERSION=$(shell git describe --tags)
COMMIT=$(shell git rev-parse --verify HEAD)
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LD_FLAGS=-s -w -linkmode external -extldflags "-static" -X main.BuildDate="$(DATE)" -X main.BuildMode="prod" -X main.BuildCommit="$(COMMIT)" -X main.BuildVersion="$(VERSION)"
BIN_NAME=build/traggo
BUILD_LICENSE=build/license/3RD_PARTY_LICENSES
UI_BUILD_LICENSE=build/license/UI_3RD_PARTY_LICENSES

license-dir:
	mkdir -p build/license || true

download-tools:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports
	GO111MODULE=off go get -u github.com/gobuffalo/packr/v2/packr2

generate-go:
	go run hack/gqlgen.go

generate-js:
	(cd ui && yarn generate)

generate: generate-go generate-js

lint-go:
	go vet ./...
	golint -set_exit_status $(shell go list ./...)
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
	go test --race -v -coverprofile=coverage.txt -covermode=atomic ./...

test: test-go

install-go:
	go mod download

install-js:
	(cd ui && yarn)

build-bin-js:
	(cd ui && yarn build)

packr:
	packr2

packr-clean:
	packr2 clean

build-bin-go:
	CGO_ENABLED=1 go build -a -ldflags '${LD_FLAGS}' -tags '${TAGS}' -o ${BIN_NAME}
	upx ${BIN_NAME} || true

build-bin: build-bin-js packr build-bin-go packr-clean

build-docker:
	cp ${BIN_NAME} docker/traggo
	(cd docker && docker build -t traggo/server:amd64-${VERSION} -t traggo/server:amd64-latest .)

docker-push:
	docker push traggo/server:amd64-${VERSION} traggo/server:amd64-latest

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
	find build/* -maxdepth 0 -type f -exec zip -9 -j {}.zip {} build/license/3RD_PARTY_LICENSES build/license/UI_3RD_PARTY_LICENSES \;

install: install-go install-js
