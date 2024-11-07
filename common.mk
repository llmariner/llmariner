# ------------------------------------------------------------------------------
# GIT
# ------------------------------------------------------------------------------

.PHONY: git-clean-check
git-clean-check:
	git diff --stat --exit-code

# ------------------------------------------------------------------------------
# GO - TESTS
# ------------------------------------------------------------------------------

.PHONY: go-test-all
go-test-all: go-test go-test-race

.PHONY: go-test
go-test:
	@go test ./...

.PHONY: go-test-race
go-test-race:
	@go test -race ./...

# ------------------------------------------------------------------------------
# GO - LINTERS
# ------------------------------------------------------------------------------

.PHONY: go-lint-all
go-lint-all: go-fmt go-vet go-mod-tidy go-errcheck go-lint

.PHONY: go-fmt
go-fmt:
	@output=$$(go fmt ./...); \
	echo $$output; \
	test -z "$${output}"

.PHONY: go-vet
go-vet:
	@go vet ./...

.PHONY: go-mod-tidy
go-mod-tidy:
	@go mod tidy

.PHONY: go-errcheck
go-errcheck:
	@errcheck ./...

.PHONY: go-lint
go-lint:
	@golint -set_exit_status ./...

# ------------------------------------------------------------------------------
# Buf rules
# ------------------------------------------------------------------------------

buf-generate-all: buf-generate-base buf-generate-ts

.PHONY: buf-mod-update
buf-mod-update:
	buf mod update

.PHONY: buf-generate-base
buf-generate-base: buf-mod-update
	buf generate

.PHONY: buf-generate-ts
buf-generate-ts: buf-mod-update
ifneq (,$(wildcard ./buf.gen.ts.yaml))
	buf generate --template buf.gen.ts.yaml --include-imports --include-wkt
endif

.PHONY: typescript-compile
typescript-compile: buf-generate-ts
ifneq (,$(wildcard ./buf.gen.ts.yaml))
	tsc --skipLibCheck
endif

buf-lint:
	buf lint
