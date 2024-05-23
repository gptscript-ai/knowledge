package transformers

import "github.com/gptscript-ai/knowledge/pkg/datastore/types"

var TransformerMap = map[string]types.DocumentTransformer{
	"extra_metadata": &ExtraMetadata{},
}
