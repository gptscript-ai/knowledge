package textractor

import (
	"sort"

	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

// Query represents a query with associated information, including an identifier,
// text, alias, query pages, results, a page, and raw block data.
type Query struct {
	id         string         // Identifier for the query
	text       string         // Text associated with the query
	alias      string         // Alias for the query
	queryPages []string       // Pages to which the query is applied
	results    []*QueryResult // Results associated with the query
	page       *Page          // Page information
	raw        types.Block    // Raw block data
}

// Text returns the text associated with the query.
func (q *Query) Text() string {
	return q.text
}

// Alias returns the alias for the query.
func (q *Query) Alias() string {
	return q.alias
}

func (q *Query) HasResult() bool {
	return len(q.results) > 0
}

// TopResult retrieves the top result by confidence score, if any are available.
func (q *Query) TopResult() *QueryResult {
	r := q.ResultsByConfidence()
	if len(r) > 0 {
		return r[0]
	}

	return nil
}

// ResultsByConfidence lists this query instance's results, sorted from most to least confident.
func (q *Query) ResultsByConfidence() []*QueryResult {
	sortedResults := make([]*QueryResult, len(q.results))
	copy(sortedResults, q.results)
	sort.Slice(sortedResults, func(i, j int) bool {
		return sortedResults[j].Confidence() < sortedResults[i].Confidence()
	})

	return sortedResults
}

// QueryResult represents the result of a parsed query.
type QueryResult struct {
	base
	text string
}

// Text returns the extracted text from the query result.
func (qr *QueryResult) Text() string {
	return qr.text
}
