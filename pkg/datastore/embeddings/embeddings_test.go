package embeddings

import (
	"github.com/gptscript-ai/knowledge/pkg/config"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/openai"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/vertex"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestLoadConfOpenAI(t *testing.T) {

	// Unset the OPENAI_API_KEY env var so test passes even if it's set in the system env
	originalEnv := os.Getenv("OPENAI_API_KEY")
	defer os.Setenv("OPENAI_API_KEY", originalEnv)
	_ = os.Unsetenv("OPENAI_API_KEY")

	dotenv := "test_assets/openai_env"
	require.NoError(t, godotenv.Load(dotenv))

	configFile := "test_assets/testcfg.yaml"
	cfg, err := config.LoadConfig(configFile)
	require.NoError(t, err)

	p, err := GetEmbeddingsModelProvider("openai", cfg.EmbeddingsConfig)
	require.NoError(t, err)
	require.Equal(t, "openai", p.Name())

	t.Logf("Provider: %s", p.Name())
	t.Logf("Config: %#v", p.Config())

	conf := p.Config().(openai.OpenAIConfig)

	assert.Equal(t, "https://foo.bar.spam", conf.APIBase) // this is in config and env, so env should take precedence
	assert.Equal(t, "sk-1234567890abcdef", conf.APIKey)   // this should come from config
}

func TestLoadConfVertex(t *testing.T) {
	dotenv := "test_assets/vertex_env"
	require.NoError(t, godotenv.Load(dotenv))

	cfg, err := config.LoadConfig("")
	require.NoError(t, err)

	// Load the configuration
	p, err := GetEmbeddingsModelProvider("vertex", cfg.EmbeddingsConfig)
	require.NoError(t, err)
	require.Equal(t, "vertex", p.Name())

	t.Logf("Config: %#v", p.Config())

	conf := p.Config().(*vertex.EmbeddingProviderVertex)

	require.Equal(t, "foo-embedding-001", conf.Model)
	require.Equal(t, "foo-project", conf.Project)
}
