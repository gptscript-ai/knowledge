GO_TAGS ?= netgo
build:
	CGO_ENABLED=0 go build -o bin/knowledge -tags "${GO_TAGS}" -ldflags "-s -w" .

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



# cross-compilation for all targets
build-cross: LDFLAGS += -extldflags "-static"
build-cross:
	CGO_ENABLED=0 gox -parallel=3 -output="_dist/$(BINARIES)-{{.OS}}-{{.Arch}}" -osarch='$(TARGETS)' $(GOFLAGS) $(if $(TAGS),-tags '$(TAGS)',) -ldflags '$(LDFLAGS)'
gen-checksum:	build-cross
	$(eval ARTIFACTS_TO_PUBLISH := $(shell ls _dist/*))
	$$(sha256sum $(ARTIFACTS_TO_PUBLISH) > _dist/checksums.txt)

.PHONY: install-tools
install-tools:
ifndef HAS_GOX
	($(GO) install $(PKG_GOX))
endif
ifndef HAS_GOLANGCI
	(curl -sfL $(PKG_GOLANGCI_LINT_SCRIPT) | sh -s -- -b $(GOENVPATH)/bin v${PKG_GOLANGCI_LINT_VERSION})
endif
ifdef HAS_GOLANGCI
ifeq ($(HAS_GOLANGCI_VERSION),)
ifdef INTERACTIVE
	@echo "Warning: Your installed version of golangci-lint (interactive: ${INTERACTIVE}) differs from what we'd like to use. Switch to v${PKG_GOLANGCI_LINT_VERSION}? [Y/n]"
	@read line; if [ $$line == "y" ]; then (curl -sfL $(PKG_GOLANGCI_LINT_SCRIPT) | sh -s -- -b $(GOENVPATH)/bin v${PKG_GOLANGCI_LINT_VERSION}); fi
else
	@echo "Warning: you're not using the same version of golangci-lint as us (v${PKG_GOLANGCI_LINT_VERSION})"
endif
endif
endif