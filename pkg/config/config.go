package config

type DatabaseConfig struct {
	DSN         string `usage:"Server database connection string (default \"sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db\")" default:"" env:"KNOW_DB_DSN"`
	AutoMigrate string `usage:"Auto migrate database" default:"true" env:"KNOW_DB_AUTO_MIGRATE"`
}

type VectorDBConfig struct {
	VectorDBPath string `usage:"VectorDBPath to the vector database (default \"$XDG_DATA_HOME/gptscript/knowledge/vector.db\")" default:"" env:"KNOW_VECTOR_DB_PATH"`
}
