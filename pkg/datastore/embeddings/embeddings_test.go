package embeddings

import (
	"github.com/gptscript-ai/knowledge/pkg/config"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/openai"
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
	type Conf struct {
		Embeddings map[string]any `yaml:"embeddings"`
	}

	configFile := "test_assets/testcfg.yaml"
	cfg, err := config.LoadConfig(configFile)
	require.NoError(t, err)

	p, err := GetEmbeddingsModelProvider("openai", cfg.EmbeddingsConfig)
	require.NoError(t, err)
	require.Equal(t, "openai", p.Name())

	t.Logf("Provider: %s", p.Name())
	t.Logf("Config: %#v", p.Config())

	conf := p.Config().(openai.OpenAIConfig)

	require.Equal(t, "https://foo.bar.com", conf.APIBase)
	assert.Equal(t, "sk-1234567890abcdef", conf.APIKey)
}

func TestLoadConfGoogleVertexAI(t *testing.T) {
	dotenv := "test_assets/google_vertex_ai.env"
	require.NoError(t, godotenv.Load(dotenv))

	cfg, err := config.LoadConfig("")
	require.NoError(t, err)

	// Load the configuration
	p, err := GetEmbeddingsModelProvider("google_vertex_ai", cfg.EmbeddingsConfig)
	require.NoError(t, err)
	require.Equal(t, "google_vertex_ai", p.Name())

	t.Logf("Config: %#v", p.Config())

	conf := p.Config().(*vertex.EmbeddingProviderGoogleVertexAI)

	require.Equal(t, "foo-embedding-001", conf.Model)
	require.Equal(t, "foo-project", conf.Project)
}
