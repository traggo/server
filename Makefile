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
	echo "" > coverage.txt
	for d in $(shell go list ./... | grep -v vendor); do \
		go test -v -coverprofile=profile.out -covermode=atomic $$d ; \
		if [ -f profile.out ]; then  \
			cat profile.out >> coverage.txt ; \
			rm profile.out ; \
		fi \
	done

test: test-go

install-go:
	go mod vendor

install: install-go
