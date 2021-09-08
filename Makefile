export GO111MODULE=on
export SHELL=/bin/bash
export GOVCS=*:git

GO = go

LOGPREFIX = $(shell printf "::")

MODULE   = $(shell env GO111MODULE=on $(GO) list -m)
CMDDIR   = $(CURDIR)/cmd
MAIN     = $(CMDDIR)/main.go
PKGS     = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))
TESTPKGS = $(shell $(GO) list -f \
			'{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' \
			$(PKGS))
ARGS     = -cover
TIMEOUT  = 15
COVEROUT = cover.out

# To echo commands executed in targets, set VERBOSE to 1
VERBOSE ?= 0
ECHOSKIP = $(if $(filter 1,$(VERBOSE)),,@)

.PHONY: lint
lint: ; $(info $(LOGPREFIX) running linter...) @ ## Run lint
	$(ECHOSKIP) golint $(PKGS)

.PHONY: metalint
metalint: ; $(info $(LOGPREFIX) running golangci-lint...) @ ## Run metalint
	$(ECHOSKIP) golangci-lint run

.PHONY: fmt
fmt: ; $(ECHOSKIP) $(info $(LOGPREFIX) running gofmt...) @ ## Run gofmt on all source files
	$(ECHOSKIP) $(GO) fmt $(PKGS)

cover: ARGS=-coverprofile=$(COVEROUT) ## Run tests with coverage
cover: fmt lint test; $(info $(LOGPREFIX) running $(NAME:%=% )cover report...) @
	$(ECHOSKIP) $(GO) tool cover -html=$(COVEROUT)
	-rm -f $(COVEROUT)

TEST_TARGETS := test check test-race test-verbose
.PHONY: $(TEST_TARGETS)
test-race:      ARGS=-race  ## Run tests with race detection on.
test-verbose:   ARGS=-v     ## Run tests with verbose output.
test:                       ## Run tests
$(TEST_TARGETS): fmt lint ; $(info $(LOGPREFIX) running $(NAME:%=% )tests...) @
	$(ECHOSKIP) $(GO) test $(TESTCOUNT) -timeout $(TIMEOUT)s $(ARGS) $(TESTPKGS)

.PHONY: make-bin
make-bin:
	$(ECHOSKIP) mkdir -p ./bin

.PHONY: build
build: make-bin ; $(info $(LOGPREFIX) building...) @  ## Build a binary
	$(ECHOSKIP) $(GO) build -o ./bin/main ./cmd/...

.PHONY: build-linux
build-linux: make-bin ; $(info $(LOGPREFIX) building linux version...) @  ## Build a linux binary
	$(ECHOSKIP) CGO_ENABLED=0 GOOS=linux $(GO) build -a -installsuffix cgo -o ./bin/main ./cmd/...

.PHONY: help
help: ## Show this help
	$(ECHOSKIP) grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: build ; $(info $(LOGPREFIX) installing...) @ ## Install poleno binary to $HOME/.bin
	$(ECHOSKIP) install -m 755 bin/main ~/.bin/poleno

clean: ; $(info $(LOGPREFIX) cleaning up...) @ ## Clean product files
	$(ECHOSKIP) -rm -rf ./bin
	$(ECHOSKIP) -rm -f $(COVEROUT)
