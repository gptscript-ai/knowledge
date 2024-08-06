#!/bin/bash
set -ex

GO_TAGS="netgo"
LD_FLAGS="-s -w -X github.com/gptscript-ai/knowledge/version.Version=${GIT_TAG}"

export CGO_ENABLED=1

if [ "$(go env GOOS)" = "linux" ]; then
  # Linux: amd64, arm64
  GOARCH=amd64 go build -o dist/knowledge-linux-amd64 -tags "${GO_TAGS}" -ldflags "${LD_FLAGS}\" -extldflags \"-static\" " .
else

  # Windows: amd64
  CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -o dist/knowledge-windows-amd64.exe -tags "${GO_TAGS}" -ldflags "${LD_FLAGS}" .

  # Darwin: amd64, arm64
  GOARCH=amd64 go build -o dist/knowledge-darwin-amd64 -tags "${GO_TAGS}" -ldflags "${LD_FLAGS}" .
  GOARCH=arm64 go build -o dist/knowledge-darwin-arm64 -tags "${GO_TAGS}" -ldflags "${LD_FLAGS}" .
fi
