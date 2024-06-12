#!/bin/bash
set -ex

GO_TAGS="netgo"
LD_FLAGS="-s -w -X github.com/gptscript-ai/knowledge/version.Version=${GIT_TAG}"

if [ "$(go env GOOS)" = "linux" ]; then
  CGO_ENABLED=1 GOARCH=amd64 go build -o dist/knowledge-linux-amd64 -tags "${GO_TAGS}" -ldflags "${LD_FLAGS}\" -extldflags \"-static\" " .
else
  CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -o dist/knowledge-windows-amd64 -tags "${GO_TAGS}" -ldflags "${LD_FLAGS}" .
  CGO_ENABLED=1 GOARCH=amd64 go build -o dist/knowledge-darwin-amd64 -tags "${GO_TAGS}" -ldflags "${LD_FLAGS}" .
  CGO_ENABLED=1 GOARCH=arm64 go build -o dist/knowledge-darwin-arm64 -tags "${GO_TAGS}" -ldflags "${LD_FLAGS}" .
fi
