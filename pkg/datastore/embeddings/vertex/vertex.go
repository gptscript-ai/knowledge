package vertex

import (
	"context"
	"dario.cat/mergo"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/load"
	cg "github.com/philippgille/chromem-go"
	"regexp"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"

	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

type EmbeddingProviderGoogleVertexAI struct {
	APIEndpoint          string `koanf:"apiEndpoint" env:"GOOGLE_VERTEX_AI_API_ENDPOINT"`
	Project              string `koanf:"project" env:"GOOGLE_VERTEX_AI_PROJECT"`
	Model                string `koanf:"model" env:"GOOGLE_VERTEX_AI_MODEL"`
	Task                 string `koanf:"task" env:"GOOGLE_VERTEX_AI_TASK"`
	OutputDimensionality *int   `koanf:"outputDimensionality" env:"GOOGLE_VERTEX_AI_OUTPUT_DIMENSIONALITY"`
}

const EmbeddingProviderGoogleVertexAIName = "google_vertex_ai"

func (p *EmbeddingProviderGoogleVertexAI) Name() string {
	return EmbeddingProviderGoogleVertexAIName
}

func New(configFile string) (*EmbeddingProviderGoogleVertexAI, error) {
	p := &EmbeddingProviderGoogleVertexAI{}

	err := load.FillConfig(configFile, "GOOGLE_VERTEX_AI_", &p)
	if err != nil {
		return nil, fmt.Errorf("failed to fill GoogleVertexAI config")
	}

	if err := p.fillDefaults(); err != nil {
		return nil, fmt.Errorf("failed to fill GoogleVertexAI defaults: %w", err)
	}

	return p, nil
}

func (p *EmbeddingProviderGoogleVertexAI) fillDefaults() error {
	defaultCfg := EmbeddingProviderGoogleVertexAI{
		APIEndpoint:          "",
		Project:              "",
		Task:                 "SEMANTIC_SIMILARITY",
		Model:                "text-embedding-004",
		OutputDimensionality: nil,
	}

	if err := mergo.Merge(p, defaultCfg); err != nil {
		return fmt.Errorf("failed to merge GoogleVertexAI config: %w", err)
	}

	return nil
}

// embedTexts is mostly taken from the official Google Cloud AI Platform Vertex AI documentation
func embedTexts(ctx context.Context,
	apiEndpoint, project, model string, texts []string,
	task string, customOutputDimensionality *int) ([]float32, error) {

	var opts []option.ClientOption
	if apiEndpoint != "" {
		opts = append(opts, option.WithEndpoint(apiEndpoint))
	}

	if task == "" {
		task = "SEMANTIC_SIMILARITY"
	}

	client, err := aiplatform.NewPredictionClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	match := regexp.MustCompile(`^(\w+-\w+)`).FindStringSubmatch(apiEndpoint)
	location := "us-central1"
	if match != nil {
		location = match[1]
	}
	endpoint := fmt.Sprintf("projects/%s/locations/%s/publishers/google/models/%s", project, location, model)
	instances := make([]*structpb.Value, len(texts))
	for i, text := range texts {
		instances[i] = structpb.NewStructValue(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"content":   structpb.NewStringValue(text),
				"task_type": structpb.NewStringValue(task),
			},
		})
	}
	outputDimensionality := structpb.NewNullValue()
	if customOutputDimensionality != nil {
		outputDimensionality = structpb.NewNumberValue(float64(*customOutputDimensionality))
	}
	params := structpb.NewStructValue(&structpb.Struct{
		Fields: map[string]*structpb.Value{"outputDimensionality": outputDimensionality},
	})

	req := &aiplatformpb.PredictRequest{
		Endpoint:   endpoint,
		Instances:  instances,
		Parameters: params,
	}
	resp, err := client.Predict(ctx, req)
	if err != nil {
		return nil, err
	}
	prediction := resp.Predictions[0]
	values := prediction.GetStructValue().Fields["embeddings"].GetStructValue().Fields["values"].GetListValue().Values
	embeddings := make([]float32, len(values))
	for j, value := range values {
		embeddings[j] = float32(value.GetNumberValue())
	}

	return embeddings, nil
}

func (p *EmbeddingProviderGoogleVertexAI) EmbeddingFunc() (cg.EmbeddingFunc, error) {
	return func(ctx context.Context, text string) ([]float32, error) {
		return embedTexts(ctx, p.APIEndpoint, p.Project, p.Model, []string{text}, p.Task, p.OutputDimensionality)
	}, nil
}

func (p *EmbeddingProviderGoogleVertexAI) Config() any {
	return p
}
