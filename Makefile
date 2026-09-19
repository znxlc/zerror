.PHONY: test lint tag tag-major

GO ?= go
GOLANGCI_LINT ?= golangci-lint
REMOTE ?= origin

test:
	$(GO) test ./... -race -cover

lint:
	$(GOLANGCI_LINT) run ./...

# Version format: v<version>.<major>.<minor>
# Usage:
#   make tag                -> bumps minor (last) number: v0.1.3 -> v0.1.4
#   make tag VERSION=v1.2.3 -> uses the given version
#   make tag-major          -> bumps major (middle) number: v0.1.3 -> v0.2.0
tag:
	@set -e; \
	version="$(VERSION)"; \
	if [ -z "$$version" ]; then \
		last=$$(git tag --list 'v*' --sort=-v:refname | head -n 1); \
		if [ -z "$$last" ]; then \
			version="v0.1.0"; \
		else \
			version=$$(echo "$$last" | awk -F. '{sub(/^v/,"",$$1); printf "v%s.%s.%d\n", $$1, $$2, $$3+1}'); \
		fi; \
	fi; \
	case "$$version" in v*) ;; *) version="v$$version";; esac; \
	echo "tagging $$version"; \
	git tag -a "$$version" -m "$$version"; \
	git push $(REMOTE) "$$version"

tag-major:
	@set -e; \
	last=$$(git tag --list 'v*' --sort=-v:refname | head -n 1); \
	if [ -z "$$last" ]; then \
		version="v0.1.0"; \
	else \
		version=$$(echo "$$last" | awk -F. '{sub(/^v/,"",$$1); printf "v%s.%d.0\n", $$1, $$2+1}'); \
	fi; \
	echo "tagging $$version"; \
	git tag -a "$$version" -m "$$version"; \
	git push $(REMOTE) "$$version"
