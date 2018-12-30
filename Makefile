download-tools:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports

generate-go:
# this is a workaround as gqlgen currently does not support go modules
# https://github.com/99designs/gqlgen/issues/226
	GO111MODULE=off go run hack/gqlgen.go

generate: generate-go

lint-go:
	go vet ./...
	golint -set_exit_status $(shell go list ./...)
	goimports -l $(shell find . -type f -name '*.go' -not -path "./vendor/*")

lint: lint-go

format-go:
	goimports -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

format: format-go

test-go:
	go test --race -v -coverprofile=coverage.txt -covermode=atomic ./...

test: test-go

install-go:
	go mod vendor

install: install-go
