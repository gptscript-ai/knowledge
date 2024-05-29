// Package postprocessors is basically the same as package transformers, but used at a different stage of the RAG pipeline
package postprocessors

import (
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/transformers"
	"github.com/gptscript-ai/knowledge/pkg/datastore/types"
)

// Postprocessor may be a "normal"
type Postprocessor types.DocumentTransformer

var PostprocessorMap = map[string]Postprocessor{}

func GetPostprocessor(name string) (Postprocessor, error) {
	var postprocessor Postprocessor
	var ok bool
	postprocessor, ok = PostprocessorMap[name]
	if !ok {
		var err error
		postprocessor, err = transformers.GetTransformer(name)
		if err != nil {
			return nil, fmt.Errorf("unknown postprocessor %q", name)
		}
	}
	return postprocessor, nil
}
