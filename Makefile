# get git tag
ifneq ($(GIT_TAG_OVERRIDE),)
$(info GIT_TAG set from env override!)
GIT_TAG := $(GIT_TAG_OVERRIDE)
endif

GIT_TAG   ?= $(shell git describe --tags)
ifeq ($(GIT_TAG),)
GIT_TAG   := $(shell git describe --always)
endif

GO_TAGS := netgo
LD_FLAGS := -s -w -X github.com/gptscript-ai/knowledge/version.Version=${GIT_TAG}
build:
	go build -o bin/knowledge -tags "${GO_TAGS}" -ldflags '$(LD_FLAGS) ' .

run: build
	bin/knowledge server

run-dev: generate run

clean-dev:
	rm -rf knowledge.db vector.db

generate: tools openapi

openapi:
	swag init --parseDependency -g pkg/server/server.go -o pkg/docs

tools:
	if ! command -v swag &> /dev/null; then go install github.com/swaggo/swag/cmd/swag@latest; fi

lint:
	golangci-lint run ./...

test:
	go test -v ./...

build-cross:
	GIT_TAG=${GIT_TAG} ./scripts/cross-build.sh

ci-setup:
	@echo "### Installing Go tools..."
	@echo "### -> Installing golangci-lint..."
	curl -sfL $(PKG_GOLANGCI_LINT_SCRIPT) | sh -s -- -b $(GOENVPATH)/bin v$(PKG_GOLANGCI_LINT_VERSION)

