SHELL=/bin/bash -o pipefail

GIT_ROOT_DIR=$(shell git rev-parse --show-toplevel)
API_SPEC_BUNDLE ?= "${GIT_ROOT_DIR}/app/docs/open-api-spec/index_bundled_openapi.yaml"
GENERATED_API_CODE_OUTPUT ?= "$(GIT_ROOT_DIR)/app/internal/infra/http/generated"
API_SPEC_ROOT ?= "${GIT_ROOT_DIR}/app/docs/open-api-spec"
GENERATED_HTTP_INFRA_OUTPUT ?= "$(GIT_ROOT_DIR)/app/internal/infra/http/generatedhttpinfra"

.PHONY: default
default: generate_http_infra

.PHONY: bundle_api
bundle_api:
	@echo "Generating the API spec bundle"
	swagger-cli bundle -t yaml $(API_SPEC_ROOT)/index.yaml -o $(API_SPEC_BUNDLE)

.PHONY: lint_api
lint_api: bundle_api
	@echo "Executing linters on API Spec"
	@spectral lint ${API_SPEC_BUNDLE} --ruleset ${GIT_ROOT_DIR}/app/docs/config/spectral.yaml --format pretty --fail-severity=error --verbose

.PHONY: serve_api
serve_api:
	@echo "Serving API Spec:"
	swagger-ui-watcher ${API_SPEC_BUNDLE}

.PHONY: generate_http_infra
generate_http_infra: bundle_api generate-ogen-http-infra

.PHONY: install-ogen
install-ogen:
	@go install github.com/ogen-go/ogen/cmd/ogen@latest


.PHONY: install-mockgen
install-mockgen:
	@go install go.uber.org/mock/mockgen@latest

.PHONY: generate-ogen-http-infra
generate-ogen-http-infra: install-ogen
	@echo "Generating HTTP infra code with ogen"
	@ogen \
		-target $(GENERATED_HTTP_INFRA_OUTPUT) \
		-package generatedhttpinfra \
		-clean \
		$(API_SPEC_BUNDLE)
	@echo "✅ Code generated successfully"

# TODO: this could be done with a script or a go process that reads a file where we define all the needed interfaces to be mocked
# as it is right now, it is highly coupled to the only usecase there is
.PHONY: generate_mocks
generate_mocks: install-mockgen
	@echo "Generating mocks"
	@mkdir -p app/internal/test/generatedmocks
	@go run github.com/golang/mock/mockgen -source=internal/core/usecases/user/user_storage.go -destination=internal/test/generatedmocks/user_storage_mock.go -package=generatedmocks
	@go run github.com/golang/mock/mockgen -source=internal/core/usecases/user/user_queue.go -destination=internal/test/generatedmocks/user_queue_mock.go -package=generatedmocks
	@echo "✅ Mocks generated successfully inside internal/test/generatedmocks"