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

lint:
	golangci-lint run ./...

test:
	go test -v ./...

# cross-compilation for all targets
TARGETS ?= darwin/amd64 darwin/arm64 linux/amd64 linux/386 linux/arm linux/arm64 windows/amd64
build-cross: LDFLAGS += -extldflags "-static"
build-cross:
	CGO_ENABLED=0 gox -parallel=3 -output="dist/knowledge-{{.OS}}-{{.Arch}}" -osarch='$(TARGETS)' $(GOFLAGS) $(if $(GO_TAGS),-tags '$(TAGS)',) -ldflags '$(LDFLAGS)'
gen-checksum:	build-cross
	$(eval ARTIFACTS_TO_PUBLISH := $(shell ls dist/*))
	$$(sha256sum $(ARTIFACTS_TO_PUBLISH) > dist/checksums.txt)

ci-setup:
	@echo "### Installing Go tools..."
	@echo "### -> Installing golangci-lint..."
	curl -sfL $(PKG_GOLANGCI_LINT_SCRIPT) | sh -s -- -b $(GOENVPATH)/bin v$(PKG_GOLANGCI_LINT_VERSION)

	@echo "### -> Installing gox..."
	./scripts/install-tools.sh gox