package bm25

import (
	"log/slog"
	"strings"

	"github.com/iwilltry42/bm25-go/bm25"
	"github.com/jmcarbo/stopwords"
)

func CleanStopwords(content string, docID string, languages []string) string {
	if len(languages) > 0 {
		langCodes := languages
		if languages[0] == "auto" {
			langCodes = []string{}
		}
		cleanedContent, langs, removed, total := stopwords.GetLanguage([]byte(content), langCodes)
		slog.Debug("Removed stopwords", "langs", langs, "document", docID, "removed", removed, "total", total, "lenContent", len(content), "lenCleaned", len(cleanedContent))
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
