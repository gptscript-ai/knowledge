GO_TAGS ?= netgo
build:
	CGO_ENABLED=0 go build -o bin/knowledge -tags "${GO_TAGS}" -ldflags "-s -w" .

openapi:
	swag init -g pkg/server/server.go -o pkg/docs