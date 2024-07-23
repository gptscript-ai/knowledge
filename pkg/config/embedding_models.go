package config

type EmbeddingModelProviderConfig struct {
	EmbeddingModelProvider string `usage:"Embedding model provider" default:"openai" env:"KNOW_EMBEDDING_MODEL_PROVIDER"`
}

type OpenAIConfig struct {
	APIBase           string `usage:"OpenAI API base" default:"https://api.openai.com/v1" env:"OPENAI_BASE_URL" koanf:"baseURL"` // clicky-chats
	APIKey            string `usage:"OpenAI API key (not required if used with clicky-chats)" default:"sk-foo" env:"OPENAI_API_KEY" koanf:"apiKey"`
	Model             string `usage:"OpenAI model" default:"gpt-4" env:"OPENAI_MODEL" koanf:"openai-model"`
	EmbeddingModel    string `usage:"OpenAI Embedding model" default:"text-embedding-ada-002" env:"OPENAI_EMBEDDING_MODEL" koanf:"embeddingModel"`
	EmbeddingEndpoint string `usage:"OpenAI Embedding endpoint" default:"/embeddings" env:"OPENAI_EMBEDDING_ENDPOINT" koanf:"embeddingEndpoint"`
	APIVersion        string `usage:"OpenAI API version (for Azure)" default:"2024-02-01" env:"OPENAI_API_VERSION" koanf:"apiVersion"`
	APIType           string `usage:"OpenAI API type (OPEN_AI, AZURE, AZURE_AD)" default:"OPEN_AI" env:"OPENAI_API_TYPE" koanf:"apiType"`
	AzureOpenAIConfig
}

type AzureOpenAIConfig struct {
	Deployment string `usage:"Azure OpenAI deployment name (overrides openai-embedding-model, if set)" default:"" env:"OPENAI_AZURE_DEPLOYMENT" koanf:"azureDeployment"`
}
