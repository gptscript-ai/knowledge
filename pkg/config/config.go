package config

type OpenAIConfig struct {
	APIBase        string `usage:"OpenAI API base" default:"https://api.openai.com/v1" env:"OPENAI_BASE_URL" name:"openai-api-base"` // clicky-chats
	APIKey         string `usage:"OpenAI API key (not required if used with clicky-chats)" default:"sk-foo" env:"OPENAI_API_KEY" name:"openai-api-key"`
	EmbeddingModel string `usage:"OpenAI Embedding model" default:"text-embedding-ada-002" env:"OPENAI_EMBEDDING_MODEL" name:"openai-embedding-model"`
	APIVersion     string `usage:"OpenAI API version (for Azure)" default:"2024-02-01" env:"OPENAI_API_VERSION" name:"openai-api-version"`
	APIType        string `usage:"OpenAI API type (OPEN_AI, AZURE, AZURE_AD)" default:"OPEN_AI" env:"OPENAI_API_TYPE" name:"openai-api-type"`
}

type AzureOpenAIConfig struct {
	Deployment string `usage:"Azure OpenAI deployment name" default:"" env:"OPENAI_AZURE_DEPLOYMENT" name:"openai-azure-deployment"`
}

type DatabaseConfig struct {
	DSN         string `usage:"Server database connection string (default \"sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db\")" default:"" env:"KNOW_DB_DSN"`
	AutoMigrate string `usage:"Auto migrate database" default:"true" env:"KNOW_DB_AUTO_MIGRATE"`
}

type VectorDBConfig struct {
	VectorDBPath string `usage:"VectorDBPath to the vector database (default \"$XDG_DATA_HOME/gptscript/knowledge/vector.db\")" default:"" env:"KNOW_VECTOR_DB_PATH"`
}
