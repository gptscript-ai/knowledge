package embeddings

import (
	"github.com/gptscript-ai/knowledge/pkg/config"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/vertex"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoadConfOpenAI(t *testing.T) {
	dotenv := "test_assets/openai.env"
	require.NoError(t, godotenv.Load(dotenv))

	// Load the configuration
	p, err := GetEmbeddingsModelProvider("openai", "test_assets/testcfg.yaml")
	require.NoError(t, err)
	require.Equal(t, "openai", p.Name())

	t.Logf("Provider: %s", p.Name())
	t.Logf("Config: %#v", p.Config())

	cfg := p.Config().(config.OpenAIConfig)

	require.Equal(t, "https://foo.bar.com", cfg.APIBase)
	assert.Equal(t, "sk-1234567890abcdef", cfg.APIKey)
}

func TestLoadConfGoogleVertexAI(t *testing.T) {
	dotenv := "test_assets/google_vertex_ai.env"
	require.NoError(t, godotenv.Load(dotenv))

	// Load the configuration
	p, err := GetEmbeddingsModelProvider("google_vertex_ai", "")
	require.NoError(t, err)
	require.Equal(t, "google_vertex_ai", p.Name())

	t.Logf("Config: %#v", p.Config())

	cfg := p.Config().(*vertex.EmbeddingProviderGoogleVertexAI)

	require.Equal(t, "foo-embedding-001", cfg.Model)
	require.Equal(t, "foo-project", cfg.Project)
}
