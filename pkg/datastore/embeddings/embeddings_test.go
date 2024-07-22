package embeddings

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoadConf(t *testing.T) {
	dotenv := "test_assets/openai.env"
	require.NoError(t, godotenv.Load(dotenv))

	// Load the configuration
	p, err := GetEmbeddingsModelProvider("openai", "")
	require.NoError(t, err)

	t.Logf("Provider: %s", p.Name())
	t.Logf("Config: %#v", p.Config())

}
