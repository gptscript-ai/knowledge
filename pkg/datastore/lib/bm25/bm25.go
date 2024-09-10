package bm25

import (
	"strings"

	"github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"github.com/iwilltry42/bm25-go/bm25"
	"github.com/jmcarbo/stopwords"
)

const (
	DefaultK1 = 1.2
	DefaultB  = 0.75
)

func BM25Run(docs []vectorstore.Document, query string, k1, b float64, cleanStopwords []string) ([]float64, error) {
	return Score(BuildCorpus(docs, cleanStopwords), query, k1, b)
}

func BuildCorpus(docs []vectorstore.Document, cleanStopwords []string) []string {
	corpus := make([]string, len(docs))
	for i, doc := range docs {
		content := doc.Content
		corpus[i] = CleanStopwords(content, doc.ID, cleanStopwords)
	}
	return corpus
}

func CleanStopwords(content string, docID string, languages []string) string {
	if len(languages) > 0 {
		langCodes := languages
		if languages[0] == "auto" {
			langCodes = []string{}
		}
		cleanedContent, _, _, _ := stopwords.GetLanguage([]byte(content), langCodes)
		return string(cleanedContent)
	}
	return content
}

var whiteSpaceTokenizer = func(s string) []string { return strings.Split(s, " ") }

func Score(corpus []string, query string, k1, b float64) ([]float64, error) {
	okapi, err := bm25.NewBM25Okapi(corpus, whiteSpaceTokenizer, k1, b, nil)
	if err != nil {
		return nil, err
	}

	return okapi.GetScores(whiteSpaceTokenizer(query))
}
