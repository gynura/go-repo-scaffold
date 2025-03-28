SHELL=/bin/bash -o pipefail

GIT_ROOT_DIR=$(shell git rev-parse --show-toplevel)
API_SPEC_BUNDLE ?= "${GIT_ROOT_DIR}/app/docs/open-api-spec/index_bundled_openapi.yaml"
GENERATED_API_CODE_OUTPUT ?= "$(GIT_ROOT_DIR)/app/internal/infra/http/generated"
API_SPEC_ROOT ?= "${GIT_ROOT_DIR}/app/docs/open-api-spec"
GENERATED_HTTP_INFRA_OUTPUT ?= "$(GIT_ROOT_DIR)/app/internal/infra/http/generatedhttpinfra"

.PHONY: bundle_api
bundle_api:
	@echo "Bundling API Spec:"
	swagger-cli bundle -t yaml docs/open-api-spec/index.yaml -o docs/open-api-spec/index_bundled_openapi.yaml

.PHONY: lint_api
lint_api: bundle_api
	@echo "Executing linters on API Spec"
	@spectral lint ${API_SPEC_BUNDLE} --ruleset ${GIT_ROOT_DIR}/app/docs/config/spectral.yaml --format pretty --fail-severity=error --verbose

.PHONY: serve_api
serve_api:
	@echo "Serving API Spec:"
	swagger-ui-watcher ${API_SPEC_BUNDLE}

.PHONY: generate_http_infra
generate_http_infra: bundle generate-api

.PHONY: bundle
bundle:
	@echo "Generating the API spec bundle"
	swagger-cli bundle -t yaml $(API_SPEC_ROOT)/index.yaml -o $(API_SPEC_BUNDLE)

.PHONY: install-ogen
install-ogen:
	@go install github.com/ogen-go/ogen/cmd/ogen@latest

.PHONY: generate-api
generate-api: install-ogen
	@echo "Generating HTTP infra code"
	@ogen \
		-target $(GENERATED_HTTP_INFRA_OUTPUT) \
		-package generatedhttpinfra \
		-clean \
		$(API_SPEC_BUNDLE)
	@echo "Code generated successfully"