GO_TAGS ?= netgo
build:
	CGO_ENABLED=0 go build -o bin/knowledge -tags "${GO_TAGS}" -ldflags "-s -w" .

run: build
	bin/knowledge server

run-dev: openapi run

clean-dev:
	rm knowledge.db
	rm -r vector.db

openapi:
	swag init -g pkg/server/server.go -o pkg/docs