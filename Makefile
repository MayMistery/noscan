include common.mk

LDFLAGS += -X "$(MODULE)/version.Version=$(VERSION)" -X "$(MODULE)/version.CommitSHA=$(VERSION_HASH)"

## Build:

.PHONY: build
build: | build-noscan ## Build binary

.PHONY: build-noscan
build-noscan: ## Build noscan
	$Q $(go) build -ldflags '$(LDFLAGS)' -o .

.PHONY: test
test: | test-noscan ## Run all tests

.PHONY: test-noscan
test-noscan: ## Run noscan tests
	$Q $(go) test -v ./...

.PHONY: lint
lint: lint-noscan ## Run all linters

.PHONY: lint-noscan
lint-backend: | $(golangci-lint) ## Run noscan linters
	$Q $(golangci-lint) run -v

fmt: $(goimports) ## Format source files
	$Q $(goimports) -local $(MODULE) -w $$(find . -type f -name '*.go' -not -path "./vendor/*")

clean: clean ## Clean

## Release:

.PHONY: version

## Help:
help: ## Show this help
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target> [options]${RESET}'
	@echo ''
	@echo 'Options:'
	@$(call global_option, "V [0|1]", "enable verbose mode (default:0)")
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)