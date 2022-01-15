BUILD_GOOS ?= $(or ${DOCKER_DEFAULT_GOOS},linux)
BUILD_GOARCH ?= $(or ${DOCKER_DEFAULT_GOARCH},amd64)
BUILD_GOARM ?= 7
BUILD_CGO_ENABLED ?= 0

export GOSUMDB := off
export GOFLAGS=-mod=mod
# Go 1.13 defaults to TLS 1.3 and requires an opt-out.  Opting out for now until certs can be regenerated before 1.14
# https://golang.org/doc/go1.12#tls_1_3
export GODEBUG := tls13=0
export GOPRIVATE=bitbucket.org/geniusrabbit/*

APP_TAGS := "nats"

GOLANGLINTCI_VERSION := latest
GOLANGLINTCI := $(TMP_VERSIONS)/golangci-lint/$(GOLANGLINTCI_VERSION)
$(GOLANGLINTCI):
	$(eval GOLANGLINTCI_TMP := $(shell mktemp -d))
	cd $(GOLANGLINTCI_TMP); go get github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGLINTCI_VERSION)
	@rm -rf $(GOLANGLINTCI_TMP)
	@rm -rf $(dir $(GOLANGLINTCI))
	@mkdir -p $(dir $(GOLANGLINTCI))
	@touch $(GOLANGLINTCI)


.PHONY: deps
deps: $(GOLANGLINTCI)


.PHONY: all
all: lint cover


.PHONY: lint
lint: golint


.PHONY: golint
golint: $(GOLANGLINTCI)
	# golint -set_exit_status ./...
	golangci-lint run -v ./...


.PHONY: fmt
fmt: ## Run formatting code
	@echo "Fix formatting"
	@gofmt -w ${GO_FMT_FLAGS} $$(go list -f "{{ .Dir }}" ./...); if [ "$${errors}" != "" ]; then echo "$${errors}"; fi


.PHONY: test
test: ## Run unit tests
	go test -v -tags ${APP_TAGS} -race ./...


.PHONY: tidy
tidy:
	go mod tidy


.PHONY: cover
cover:
	@mkdir -p $(TMP_ETC)
	@rm -f $(TMP_ETC)/coverage.txt $(TMP_ETC)/coverage.html
	go test -race -coverprofile=$(TMP_ETC)/coverage.txt -coverpkg=./... ./...
	@go tool cover -html=$(TMP_ETC)/coverage.txt -o $(TMP_ETC)/coverage.html
	@echo
	@go tool cover -func=$(TMP_ETC)/coverage.txt | grep total
	@echo
	@echo Open the coverage report:
	@echo open $(TMP_ETC)/coverage.html


.PHONY: generate-code
generate-code: ## Run codegeneration procedure
	@echo "Generate code"
	@go generate ./...


.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.DEFAULT_GOAL := help
