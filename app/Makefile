SHELL=/bin/bash -o pipefail
GIT_ROOT_DIR=$(shell git rev-parse --show-toplevel)

.PHONY: bundle_api
bundle_api:
	@echo "Bundling API Spec:"
	swagger-cli bundle -t yaml docs/open-api-spec/index.yaml -o docs/open-api-spec/index_bundled_openapi.yaml

.PHONY: lint_api
lint_api: bundle_api
	@echo "Executing linters on API Spec"
	@spectral lint ${GIT_ROOT_DIR}/app/docs/open-api-spec/index_bundled_openapi.yaml --ruleset ${GIT_ROOT_DIR}/app/docs/config/spectral.yaml --format pretty --fail-severity=error --verbose

.PHONY: serve_api
serve_api:
	@echo "Serving API Spec:"
	swagger-ui-watcher ${GIT_ROOT_DIR}/app/docs/open-api-spec/index_bundled_openapi.yaml
