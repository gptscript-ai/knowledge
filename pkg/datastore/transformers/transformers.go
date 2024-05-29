package transformers

import (
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/types"
)

var TransformerMap = map[string]types.DocumentTransformer{
	"extra_metadata":                  &ExtraMetadata{},
	"filter_markdown_docs_no_content": &FilterMarkdownDocsNoContent{},
	"keywords":                        &KeywordExtractor{},
}

func GetTransformer(name string) (types.DocumentTransformer, error) {
	transformer, ok := TransformerMap[name]
	if !ok {
		return nil, fmt.Errorf("unknown transformer %q", name)
	}
	return transformer, nil
}
