# On Windows OS is set. This makefile requires Linux or MacOS
ifeq ($(OS),)
SHELL := /bin/bash
MAKE ?= make

checktool = $(shell command -v $1 2>/dev/null)
tool = $(if $(call checktool, $(firstword $1)), $1, @echo "$(firstword $1) was not found on the system. Please install it")

GO ?= $(call checktool, go)
GOTEST ?= $(GO) test
GOTEST_ARGS ?= -timeout 2m -count 1 -cover

GO_TOOLS= $(GO) run -modfile ./tools/go.mod 

.PHONY: test
test:
	@$(GO) clean -testcache
	$(GOTEST) $(GOTEST_ARGS) ./... -v

.PHONY: lint
lint:
	@$(GO_TOOLS) github.com/golangci/golangci-lint/cmd/golangci-lint run -c .golangci.yml

.PHONY: lint-fix
lint-fix:
	@$(GO_TOOLS) github.com/golangci/golangci-lint/cmd/golangci-lint run -c .golangci.yml --fix
endif