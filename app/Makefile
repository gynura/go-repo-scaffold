SHELL=/bin/bash -o pipefail
GIT_ROOT_DIR=$(shell git rev-parse --show-toplevel)
API_LINTER=tools/spectral-openapi-linter
CURDIR := $(shell pwd)

.PHONY: bundle_api
bundle_api:
	@echo "Bundling API:"
	swagger-cli bundle -t yaml docs/open-api-spec/index.yaml -o docs/open-api-spec/index_bundled_openapi.yaml
