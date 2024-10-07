.PHONY: default
default: test

include common.mk

.PHONY: test
test: go-test-all

.PHONY: lint
lint: go-lint-all git-clean-check

.PHONY: build
build:
	go build $(BUILD_FLAGS) -o ./bin/llma ./cli/cmd/

.PHONY: gendoc
gendoc:
	@go run ./cli/cmd/gendoc
