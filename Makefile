TEST ?= $(shell $(GO) list ./... | grep -v vendor)
VERSION = $(shell cat version)
REVISION = $(shell git describe --always)
BRANCH = $(shell git branch --show-current)

INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m
ifeq ("$(shell uname)","Darwin")
GO ?= GO111MODULE=on go
else
GO ?= GO111MODULE=on /usr/local/go/bin/go
endif

default: build

depsdev: ## Installing dependencies for development
	$(GO) get golang.org/x/lint/golint
	$(GO) get github.com/tcnksm/ghr
	$(GO) get github.com/mitchellh/gox

lint: ## Exec golint
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Linting$(RESET)"
	golint -min_confidence 1.1 -set_exit_status $(TEST)

build: ## Build for release (default)
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Building$(RESET)"
	./misc/build $(VERSION) $(REVISION)

ghr: ## Upload to Github releases without github token check
ifeq 'master' '$(BRANCH)'
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Releasing for Github$(RESET)"
	ghr -u heat1024 v$(VERSION)-$(REVISION) pkg
else
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Releasing for Github$(RESET)"
	ghr -u heat1024 -prerelease -recreate v$(VERSION)-manual-latest pkg
endif

dist: build ## Upload to Github releases
	@test -z $(GITHUB_TOKEN) || $(MAKE) ghr

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(INFO_COLOR)%-30s$(RESET) %s\n", $$1, $$2}'

.PHONY: default dist test
