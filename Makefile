help: ## show this help message
	@grep -E '^[a-zA-Z_-]+.*:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'

build: ## build the tinkerbell-generate binary
	CGO_ENABLED=0 go build .

.PHONY: release-local
release-local: ## Build and release all binaries locally
	goreleaser build --clean --snapshot

.PHONY: release
release: ## Build and release all binaries
	goreleaser release --clean --auto-snapshot
