package transformers

import (
	"context"
	"fmt"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"github.com/hupe1980/golc/schema"
	"log/slog"
	"strings"
)

func NewKeyWordExtractor(numKeywords int, llm schema.ChatModel) *KeywordExtractor {
	return &KeywordExtractor{
		NumKeywords: numKeywords,
		LLM:         llm,
	}
}

type KeywordExtractor struct {
	NumKeywords int
	LLM         schema.ChatModel
}

func (k *KeywordExtractor) extractKeywords(ctx context.Context, doc vs.Document) ([]string, error) {
	// Implement keyword extraction here
	result, err := k.LLM.Generate(ctx, []schema.ChatMessage{schema.NewHumanChatMessage(fmt.Sprintf(tpl, k.NumKeywords, doc.Content))})
	if err != nil {
		return nil, err
	}
	keywords := strings.Split(result.Generations[0].Message.Content(), ",")
	return keywords, nil
}

var tpl = `Extract %d keywords from the following document and return them as a comma-separated list:
%s
`

func (k *KeywordExtractor) Transform(ctx context.Context, docs []vs.Document) ([]vs.Document, error) {
	slog.Debug("Extracting keywords from documents")
	for i, doc := range docs {
		keywords, err := k.extractKeywords(ctx, doc)
		if err != nil {
			return nil, err
		}
		slog.Debug("Extracted keywords", "keywords", keywords)
		docs[i].Metadata["keywords"] = strings.Join(keywords, ",")
	}
	return docs, nil
}
