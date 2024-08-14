package retrievers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/store"
	"github.com/gptscript-ai/knowledge/pkg/llm"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"github.com/philippgille/chromem-go"
	"log/slog"
	"strings"
)

const SubqueryRetrieverName = "subquery"

type SubqueryRetriever struct {
	Model llm.LLMConfig
	Limit int
	TopK  int
}

func (s SubqueryRetriever) Name() string {
	return SubqueryRetrieverName
}

var subqueryPrompt = `The following query will be used for a vector similarity search.
If it is too complex or covering multiple topics or entities, please split it into multiple subqueries.
I.e. a comparative query like "What are the differences between cats and dogs?" could be split into subqueries concerning cats and dogs separately.
The resulting subqueries will then be used for separate vector similarity searches.
Just changing the phrasing of the input question often won't change the semantic meaning, so those may not be good candidates.
Limit the number of subqueries to a maximum of {{.limit}} (less is ok).
Query: "{{.query}}"
Reply with all subqueries in a json list like the following and don't reply with anything else (also don't use any markdown syntax).
Response schema: {"results": ["<subquery-1>", "<subquery-2>"]}`

type subqueryResp struct {
	Results []string `json:"results"`
}

func (s SubqueryRetriever) Retrieve(ctx context.Context, store store.Store, query string, datasetIDs []string, where map[string]string, whereDocument []chromem.WhereDocument) ([]vs.Document, error) {

	if len(datasetIDs) > 1 {
		return nil, fmt.Errorf("basic retriever does not support querying multiple datasets")
	}

	var datasetID string
	if len(datasetIDs) == 0 {
		datasetID = "default"
	} else {
		datasetID = datasetIDs[0]
	}

	m, err := llm.NewFromConfig(s.Model)
	if err != nil {
		return nil, err
	}

	if s.TopK <= 0 {
		s.TopK = 3
	}

	if s.Limit < 1 {
		return nil, fmt.Errorf("limit must be at least 1")
	}

	if s.Limit == 0 {
		s.Limit = 3
	}

	result, err := m.Prompt(context.Background(), subqueryPrompt, map[string]interface{}{"query": query, "limit": s.Limit})
	if err != nil {
		return nil, err
	}
	var resp subqueryResp
	err = json.Unmarshal([]byte(result), &resp)
	if err != nil {
		slog.Debug("llm response", "response", result)
		return nil, fmt.Errorf("[retrievers/subquery] failed to unmarshal llm response: %w", err)
	}

	queries := resp.Results

	slog.Debug("SubqueryQueryRetriever generated subqueries", "queries", strings.Join(queries, " | "))

	var resultDocs []vs.Document
	for _, q := range queries {
		docs, err := store.SimilaritySearch(ctx, q, s.TopK, datasetID, where, whereDocument)
		if err != nil {
			return nil, err
		}
		resultDocs = append(resultDocs, docs...)
	}

	return resultDocs, nil
}
