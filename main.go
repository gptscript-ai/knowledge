package main

import (
	acmd "github.com/acorn-io/cmd"
	"github.com/gptscript-ai/knowledge/pkg/cmd"
)

func main() {
	acmd.Main(cmd.New())
}
