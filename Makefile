GO_TAGS ?= netgo
build: openapi
	CGO_ENABLED=0 go build -o bin/knowledge -tags "${GO_TAGS}" -ldflags "-s -w" .

run: openapi build
	bin/knowledge server

clean-dev:
	rm knowledge.db
	rm -r vector.db

openapi:
	swag init -g pkg/server/server.go -o pkg/docs