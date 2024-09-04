package gptscript

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGSRead(t *testing.T) {
	ctx := context.Background()
	reader := strings.NewReader("test")
	docs, err := GSRead(ctx, reader)
	require.NoError(t, err)
	assert.Len(t, docs, 1)
	assert.Equal(t, "test", docs[0].Content)
}
