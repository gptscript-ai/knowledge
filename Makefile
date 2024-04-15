GO_TAGS ?= netgo
build:
	CGO_ENABLED=0 go build -o bin/knowledge -tags "${GO_TAGS}" -ldflags "-s -w" .

run-dev:
	go run -tags "${GO_TAGS}" -ldflags "-s -w" . server

clean-dev:
	rm knowledge.db

openapi:
	swag init -g pkg/server/server.go -o pkg/docs