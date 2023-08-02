GO ?= go
GO_TOOLS= $(GO) run -modfile ./tools/go.mod 


lint:
	@$(GO_TOOLS) github.com/golangci/golangci-lint/cmd/golangci-lint run -c .golangci.yml

lint-fix:
	@$(GO_TOOLS) github.com/golangci/golangci-lint/cmd/golangci-lint run -c .golangci.yml --fix