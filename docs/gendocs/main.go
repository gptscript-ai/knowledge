package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	knowledge "github.com/gptscript-ai/knowledge/pkg/cmd"
	"github.com/spf13/cobra/doc"
)

const fmTemplate = `---
title: "%s"
---
`

func main() {
	cmd := knowledge.New()
	cmd.DisableAutoGenTag = true

	files, err := filepath.Glob("docs/docs/04-command-line-reference/knowledge_*.md")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			log.Fatal(err)
		}
	}

	err = doc.GenMarkdownTreeCustom(cmd, "docs/docs/03-cmd", filePrepender, linkHandler)
	if err != nil {
		log.Fatal(err)
	}
}

func filePrepender(filename string) string {
	name := filepath.Base(filename)
	base := strings.TrimSuffix(name, path.Ext(name))
	return fmt.Sprintf(fmTemplate, strings.Replace(base, "_", " ", -1))
}

func linkHandler(name string) string {
	return name
}
