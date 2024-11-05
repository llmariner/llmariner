.PHONY: default
default: test

include common.mk

.PHONY: test
test: go-test-all

.PHONY: lint
lint: go-lint-all helm-lint git-clean-check

.PHONY: build
build:
	go build $(BUILD_FLAGS) -o ./bin/llma ./cli/cmd/

.PHONY: gendoc
gendoc:
	@go run ./cli/cmd/gendoc

.PHONY: check-helm-tool
check-helm-tool:
	@command -v helm-tool >/dev/null 2>&1 || $(MAKE) install-helm-tool

.PHONY: install-helm-tool
install-helm-tool:
	go install github.com/cert-manager/helm-tool@latest

.PHONY: generate-chart-schema
generate-chart-schema: check-helm-tool
	@cd ./deployments/llmariner && helm-tool schema > values.schema.json

.PHONY: helm-lint
helm-lint: generate-chart-schema
	helm lint ./deployments/llmariner

.PHONY: helm-cleanup
helm-cleanup:
	-rm deployments/llmariner/Chart.lock

.PHONY: helm-dependency
helm-dependency: helm-cleanup
	helm dependency build ./deployments/llmariner
