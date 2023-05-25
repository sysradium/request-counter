.PHONY: test clean all
GO ?= go

run::
	$(GO) run ./cmd/main.go


load::
	$(GO) run ./e2e/main.go


test::
	$(GO) test -v ./...
