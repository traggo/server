download-tools:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports

lint-go:
	go vet ./...
	golint -set_exit_status $(shell go list ./...)
	goimports -l $(shell find . -type f -name '*.go' -not -path "./vendor/*")

lint: lint-go

format-go:
	goimports -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

format: format-go

test-go:
	go test --race ./...

test: test-go

install-go:
	go mod vendor

install: install-go
