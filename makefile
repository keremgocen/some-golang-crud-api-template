GO ?= go

## Run unit tests
test-unit:
	@echo "Running unit tests."
	@CGO_ENABLED=1 $(GO) test -short -coverprofile=unit.coverprofile -covermode=atomic -race ./...

.PHONY: test-unit