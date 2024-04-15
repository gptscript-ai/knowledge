package assemblyai

import (
	"context"
)

// LeMURService groups the operations related to LeMUR.
type LeMURService struct {
	client *Client
}

// Question returns answers to free-form questions about one or more transcripts.
//
// https://www.assemblyai.com/docs/Models/lemur#question--answer
func (s *LeMURService) Question(ctx context.Context, params LeMURQuestionAnswerParams) (LeMURQuestionAnswerResponse, error) {
	var response LeMURQuestionAnswerResponse

	req, err := s.client.newJSONRequest("POST", "/lemur/v3/generate/question-answer", params)
	if err != nil {
		return LeMURQuestionAnswerResponse{}, err
	}

	if _, err := s.client.do(ctx, req, &response); err != nil {
		return LeMURQuestionAnswerResponse{}, err
	}

	return response, nil
}

// Summarize returns a custom summary of a set of transcripts.
//
// https://www.assemblyai.com/docs/Models/lemur#action-items
func (s *LeMURService) Summarize(ctx context.Context, params LeMURSummaryParams) (LeMURSummaryResponse, error) {
	req, err := s.client.newJSONRequest("POST", "/lemur/v3/generate/summary", params)
	if err != nil {
		return LeMURSummaryResponse{}, err
	}

	var response LeMURSummaryResponse

	if _, err := s.client.do(ctx, req, &response); err != nil {
		return LeMURSummaryResponse{}, err
	}

	return response, nil
}

// ActionItems returns a set of action items based on a set of transcripts.
//
// https://www.assemblyai.com/docs/Models/lemur#action-items
func (s *LeMURService) ActionItems(ctx context.Context, params LeMURActionItemsParams) (LeMURActionItemsResponse, error) {
	req, err := s.client.newJSONRequest("POST", "/lemur/v3/generate/action-items", params)
	if err != nil {
		return LeMURActionItemsResponse{}, err
	}

	var response LeMURActionItemsResponse

	if _, err := s.client.do(ctx, req, &response); err != nil {
		return LeMURActionItemsResponse{}, err
	}

	return response, nil
}

// Task lets you submit a custom prompt to LeMUR.
//
// https://www.assemblyai.com/docs/Models/lemur#task
func (s *LeMURService) Task(ctx context.Context, params LeMURTaskParams) (LeMURTaskResponse, error) {
	req, err := s.client.newJSONRequest("POST", "/lemur/v3/generate/task", params)
	if err != nil {
		return LeMURTaskResponse{}, err
	}

	var response LeMURTaskResponse

	if _, err := s.client.do(ctx, req, &response); err != nil {
		return LeMURTaskResponse{}, err
	}

	return response, nil
}

func (s *LeMURService) PurgeRequestData(ctx context.Context, requestID string) (PurgeLeMURRequestDataResponse, error) {
	req, err := s.client.newJSONRequest("DELETE", "/lemur/v3/"+requestID, nil)
	if err != nil {
		return PurgeLeMURRequestDataResponse{}, err
	}

	var response PurgeLeMURRequestDataResponse

	if _, err := s.client.do(ctx, req, &response); err != nil {
		return PurgeLeMURRequestDataResponse{}, err
	}

	return response, nil
}
