package types

type OpenAIConfig struct {
	APIBase        string `usage:"OpenAI API base" default:"http://localhost:8080/v1" env:"KNOW_OPENAI_API_BASE"` // clicky-chats
	APIKey         string `usage:"OpenAI API key (not required if used with clicky-chats)" default:"sk-foo" env:"KNOW_OPENAI_API_KEY"`
	EmbeddingModel string `usage:"OpenAI Embedding model" default:"text-embedding-ada-002" env:"KNOW_OPENAI_EMBEDDING_MODEL"`
}

type DatabaseConfig struct {
	DSN         string `usage:"Server database connection string" default:"sqlite://knowledge.db" env:"KNOW_DSN"`
	AutoMigrate string `usage:"Auto migrate database" default:"true" env:"KNOW_AUTO_MIGRATE"`
}
