package openai

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/load"
	cg "github.com/philippgille/chromem-go"
	"log/slog"
	"net/url"
	"strings"
)

const EmbeddingModelProviderOpenAIName string = "openai"

type EmbeddingModelProviderOpenAI struct {
	OpenAIConfig
}

type OpenAIConfig struct {
	BaseURL           string            `usage:"OpenAI API base" default:"https://api.openai.com/v1" env:"OPENAI_BASE_URL" koanf:"baseURL"`
	APIKey            string            `usage:"OpenAI API key (not required if used with clicky-chats)" default:"sk-foo" env:"OPENAI_API_KEY" koanf:"apiKey"`
	Model             string            `usage:"OpenAI model" default:"gpt-4" env:"OPENAI_MODEL" koanf:"openai-model"`
	EmbeddingModel    string            `usage:"OpenAI Embedding model" default:"text-embedding-ada-002" env:"OPENAI_EMBEDDING_MODEL" koanf:"embeddingModel"`
	EmbeddingEndpoint string            `usage:"OpenAI Embedding endpoint" default:"/embeddings" env:"OPENAI_EMBEDDING_ENDPOINT" koanf:"embeddingEndpoint"`
	APIVersion        string            `usage:"OpenAI API version (for Azure)" default:"2024-02-01" env:"OPENAI_API_VERSION" koanf:"apiVersion"`
	APIType           string            `usage:"OpenAI API type (OPEN_AI, AZURE, AZURE_AD, ...)" default:"OPEN_AI" env:"OPENAI_API_TYPE" koanf:"apiType"`
	AzureOpenAIConfig AzureOpenAIConfig `koanf:"azure"`
}

type AzureOpenAIConfig struct {
	Deployment string `usage:"Azure OpenAI deployment name (overrides openai-embedding-model, if set)" default:"" env:"OPENAI_AZURE_DEPLOYMENT" koanf:"deployment"`
}

func (p *EmbeddingModelProviderOpenAI) Name() string {
	return EmbeddingModelProviderOpenAIName
}

func New(c EmbeddingModelProviderOpenAI) (*EmbeddingModelProviderOpenAI, error) {

	if err := load.FillConfigEnv("OPENAI_", &c.OpenAIConfig); err != nil {
		return nil, fmt.Errorf("failed to fill OpenAI config from environment: %w", err)
	}

	if err := c.fillDefaults(); err != nil {
		return nil, fmt.Errorf("failed to fill OpenAI defaults: %w", err)
	}

	return &c, nil
}

func (p *EmbeddingModelProviderOpenAI) fillDefaults() error {
	defaultAzureOpenAIConfig := AzureOpenAIConfig{
		Deployment: "",
	}

	defaultConfig := OpenAIConfig{
		BaseURL:           "https://api.openai.com/v1",
		APIKey:            "sk-foo",
		Model:             "gpt-4",
		EmbeddingModel:    "text-embedding-ada-002",
		EmbeddingEndpoint: "/embeddings",
		APIVersion:        "2024-02-01",
		APIType:           "OPEN_AI",
		AzureOpenAIConfig: defaultAzureOpenAIConfig,
	}

	err := mergo.Merge(&defaultConfig, p.OpenAIConfig, mergo.WithOverride)
	if err != nil {
		return fmt.Errorf("failed to merge OpenAI config: %w", err)
	}

	p.OpenAIConfig = defaultConfig

	return nil
}

func (p *EmbeddingModelProviderOpenAI) EmbeddingFunc() (cg.EmbeddingFunc, error) {
	var embeddingFunc cg.EmbeddingFunc

	switch strings.ToLower(p.OpenAIConfig.APIType) {

	// except for Azure, most other OpenAI API compatible providers only differ in the normalization of output vectors (apart from the obvious API endpoint, etc.)
	case "azure", "azure_ad":
		// TODO: clean this up to support inputting the full deployment URL
		deployment := p.OpenAIConfig.AzureOpenAIConfig.Deployment
		if deployment == "" {
			deployment = p.OpenAIConfig.EmbeddingModel
		}

		deploymentURL, err := url.Parse(p.OpenAIConfig.BaseURL)
		if err != nil || deploymentURL == nil {
			return nil, fmt.Errorf("failed to parse OpenAI Base URL %q: %w", p.OpenAIConfig.BaseURL, err)
		}

		deploymentURL = deploymentURL.JoinPath("openai", "deployments", deployment)

		slog.Debug("Using Azure OpenAI API", "deploymentURL", deploymentURL.String(), "APIVersion", p.OpenAIConfig.APIVersion)

		embeddingFunc = cg.NewEmbeddingFuncAzureOpenAI(
			p.OpenAIConfig.APIKey,
			deploymentURL.String(),
			p.OpenAIConfig.APIVersion,
			"",
		)
	case "open_ai":
		cfg := cg.NewOpenAICompatConfig(
			p.OpenAIConfig.BaseURL,
			p.OpenAIConfig.APIKey,
			p.OpenAIConfig.EmbeddingModel,
		).
			WithNormalized(true).
			WithEmbeddingsEndpoint(p.OpenAIConfig.EmbeddingEndpoint)
		embeddingFunc = cg.NewEmbeddingFuncOpenAICompat(cfg)
	default:
		return nil, fmt.Errorf("unknown OpenAI API type: %q", p.OpenAIConfig.APIType)
	}

	return embeddingFunc, nil
}

func (p *EmbeddingModelProviderOpenAI) Config() any {
	return p.OpenAIConfig
}
