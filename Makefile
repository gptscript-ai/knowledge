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
