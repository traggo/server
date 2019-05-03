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

build-js:
	(cd ui && yarn build)

install: install-go install-js
