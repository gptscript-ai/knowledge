package vertex

import (
	"context"
	"fmt"
	"regexp"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"

	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

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
