

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec


.PHONY: generate
generate: mdtogo
	rm -rf internal/docs/generated
	mkdir internal/docs/generated
	GOBIN=$(LOCALBIN) go generate ./...
	go fmt ./internal/docs/generated/...

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: all
all: generate fmt vet ## Build manager binary.
	go build -ldflags "-X github.com/henderiw-kpt/kptgen/run.version=${GIT_COMMIT}" -o $(GOBIN)/kptgen -v .

##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
MDTOGO ?= $(LOCALBIN)/mdtogo

.PHONY: mdtogo
mdtogo: $(MDTOGO) ## Install mdtgo locallt
$(MDTOGO): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install ./mdtogo