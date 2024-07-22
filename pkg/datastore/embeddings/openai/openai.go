package openai

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/acorn-io/z"
	"github.com/gptscript-ai/knowledge/pkg/config"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/load"
	cg "github.com/philippgille/chromem-go"
	"log/slog"
	"net/url"
)

const EmbeddingModelProviderOpenAIName string = "openai"

type EmbeddingModelProviderOpenAI struct {
	OpenAIConfig config.OpenAIConfig `json:"openai"`
}

func (p *EmbeddingModelProviderOpenAI) Name() string {
	return EmbeddingModelProviderOpenAIName
}

func New(configFile string) (*EmbeddingModelProviderOpenAI, error) {
	p := &EmbeddingModelProviderOpenAI{}

	p.OpenAIConfig = config.OpenAIConfig{}

	err := load.FillConfig(configFile, "OPENAI_", &p.OpenAIConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to fill OpenAI config: %w", err)
	}

	if err := p.fillDefaults(); err != nil {
		return nil, fmt.Errorf("failed to fill OpenAI defaults: %w", err)
	}

	return p, nil
}

func (p *EmbeddingModelProviderOpenAI) fillDefaults() error {
	defaultAzureOpenAIConfig := config.AzureOpenAIConfig{
		Deployment: "",
	}

	defaultConfig := config.OpenAIConfig{
		APIBase:           "https://api.openai.com/v1",
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

	if p.OpenAIConfig.APIType == "Azure" {
		// TODO: clean this up to support inputting the full deployment URL
		deployment := p.OpenAIConfig.AzureOpenAIConfig.Deployment
		if deployment == "" {
			deployment = p.OpenAIConfig.EmbeddingModel
		}

		deploymentURL, err := url.Parse(p.OpenAIConfig.APIBase)
		if err != nil || deploymentURL == nil {
			return nil, fmt.Errorf("failed to parse OpenAI Base URL %q: %w", p.OpenAIConfig.APIBase, err)
		}

		deploymentURL = deploymentURL.JoinPath("openai", "deployments", deployment)

		slog.Debug("Using Azure OpenAI API", "deploymentURL", deploymentURL.String(), "APIVersion", p.OpenAIConfig.APIVersion)

		embeddingFunc = cg.NewEmbeddingFuncAzureOpenAI(
			p.OpenAIConfig.APIKey,
			deploymentURL.String(),
			p.OpenAIConfig.APIVersion,
			"",
		)
	} else {
		embeddingFunc = cg.NewEmbeddingFuncOpenAICompat(
			p.OpenAIConfig.APIBase,
			p.OpenAIConfig.APIKey,
			p.OpenAIConfig.EmbeddingModel,
			z.Pointer(true),
			cg.WithOpenAICompatEmbeddingsEndpointOverride(p.OpenAIConfig.EmbeddingEndpoint),
		)
	}

	return embeddingFunc, nil
}

func (p *EmbeddingModelProviderOpenAI) Config() any {
	return p.OpenAIConfig
}
