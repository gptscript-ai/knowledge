package transformers

import (
	"fmt"

	"github.com/gptscript-ai/knowledge/pkg/datastore/types"
)

var TransformerMap = map[string]types.DocumentTransformer{
	ExtraMetadataName:               &ExtraMetadata{},
	FilterMarkdownDocsNoContentName: &FilterMarkdownDocsNoContent{},
	KeywordExtractorName:            &KeywordExtractor{},
	MetadataManipulatorName:         &MetadataManipulator{},
}

func GetTransformer(name string) (types.DocumentTransformer, error) {
	transformer, ok := TransformerMap[name]
	if !ok {
		return nil, fmt.Errorf("unknown transformer %q", name)
	}
	return transformer, nil
}
