package assemblyai

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/google/go-querystring/query"
)

const (
	TranscriptStatusQueued     TranscriptStatus = "queued"
	TranscriptStatusProcessing TranscriptStatus = "processing"
	TranscriptStatusCompleted  TranscriptStatus = "completed"
	TranscriptStatusError      TranscriptStatus = "error"
)

const (
	SpeechModelNano SpeechModel = "nano"
)

// TranscriptService groups the operations related to transcribing audio.
type TranscriptService struct {
	client *Client
}

// SubmitFromURL submits an audio file for transcription without waiting for it
// to finish.
//
// https://www.assemblyai.com/docs/API%20reference/transcript#create-a-transcript
func (s *TranscriptService) SubmitFromURL(ctx context.Context, audioURL string, opts *TranscriptOptionalParams) (Transcript, error) {
	var transcript Transcript

	params := TranscriptParams{
		AudioURL: String(audioURL),
	}

	if opts != nil {
		params.TranscriptOptionalParams = *opts
	}

	req, err := s.client.newJSONRequest("POST", "/v2/transcript", params)
	if err != nil {
		return Transcript{}, err
	}

	resp, err := s.client.do(ctx, req, &transcript)
	if err != nil {
		return Transcript{}, err
	}
	defer resp.Body.Close()

	return transcript, nil
}

// SubmitFromReader submits audio for transcription without waiting for it to
// finish.
func (s *TranscriptService) SubmitFromReader(ctx context.Context, reader io.Reader, params *TranscriptOptionalParams) (Transcript, error) {
	u, err := s.client.Upload(ctx, reader)
	if err != nil {
		return Transcript{}, err
	}
	return s.SubmitFromURL(ctx, u, params)
}

// Delete permanently deletes a transcript.
//
// https://www.assemblyai.com/docs/API%20reference/listing_and_deleting#deleting-transcripts-from-the-api
func (s *TranscriptService) Delete(ctx context.Context, transcriptID string) (Transcript, error) {
	req, err := s.client.newJSONRequest("DELETE", fmt.Sprint("/v2/transcript/", transcriptID), nil)
	if err != nil {
		return Transcript{}, err
	}

	var transcript Transcript

	resp, err := s.client.do(ctx, req, &transcript)
	if err != nil {
		return Transcript{}, err
	}
	defer resp.Body.Close()

	return transcript, nil
}

// Get returns a transcript.
//
// https://www.assemblyai.com/docs/API%20reference/transcript
func (s *TranscriptService) Get(ctx context.Context, transcriptID string) (Transcript, error) {
	req, err := s.client.newJSONRequest("GET", fmt.Sprint("/v2/transcript/", transcriptID), nil)
	if err != nil {
		return Transcript{}, err
	}

	var transcript Transcript

	resp, err := s.client.do(ctx, req, &transcript)
	if err != nil {
		return Transcript{}, err
	}
	defer resp.Body.Close()

	return transcript, nil
}

// GetSentences returns the sentences for a transcript.
func (s *TranscriptService) GetSentences(ctx context.Context, transcriptID string) (SentencesResponse, error) {
	req, err := s.client.newJSONRequest("GET", fmt.Sprint("/v2/transcript/", transcriptID, "/sentences"), nil)
	if err != nil {
		return SentencesResponse{}, err
	}

	var results SentencesResponse

	resp, err := s.client.do(ctx, req, &results)
	if err != nil {
		return SentencesResponse{}, err
	}
	defer resp.Body.Close()

	return results, nil
}

// GetParagraphs returns the paragraphs for a transcript.
func (s *TranscriptService) GetParagraphs(ctx context.Context, transcriptID string) (ParagraphsResponse, error) {
	req, err := s.client.newJSONRequest("GET", fmt.Sprint("/v2/transcript/", transcriptID, "/paragraphs"), nil)
	if err != nil {
		return ParagraphsResponse{}, err
	}

	var results ParagraphsResponse

	resp, err := s.client.do(ctx, req, &results)
	if err != nil {
		return ParagraphsResponse{}, err
	}
	defer resp.Body.Close()

	return results, nil
}

// GetRedactedAudio returns the redacted audio for a transcript.
//
// https://www.assemblyai.com/docs/Models/pii_redaction#create-a-redacted-audio-file
func (s *TranscriptService) GetRedactedAudio(ctx context.Context, transcriptID string) (RedactedAudioResponse, error) {
	req, err := s.client.newJSONRequest("GET", fmt.Sprint("/v2/transcript/", transcriptID, "/redacted-audio"), nil)
	if err != nil {
		return RedactedAudioResponse{}, err
	}

	var audio RedactedAudioResponse

	resp, err := s.client.do(ctx, req, &audio)
	if err != nil {
		return RedactedAudioResponse{}, err
	}
	defer resp.Body.Close()

	return audio, nil
}

type TranscriptGetSubtitlesOptions struct {
	CharsPerCaption int64 `json:"chars_per_caption"`
}

func (s *TranscriptService) GetSubtitles(ctx context.Context, transcriptID string, format SubtitleFormat, opts *TranscriptGetSubtitlesOptions) ([]byte, error) {
	req, err := s.client.newRequest("GET", fmt.Sprintf("/v2/transcript/%s/%s", transcriptID, format), nil)
	if err != nil {
		return nil, err
	}

	if opts != nil {
		values := make(url.Values)
		values.Set("chars_per_caption", fmt.Sprint(opts.CharsPerCaption))
		req.URL.RawQuery = values.Encode()
	}

	resp, err := s.client.do(ctx, req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// List returns a collection of transcripts based on a filter.
//
// https://www.assemblyai.com/docs/API%20reference/listing_and_deleting#listing-historical-transcripts
func (s *TranscriptService) List(ctx context.Context, options ListTranscriptParams) (TranscriptList, error) {
	req, err := s.client.newJSONRequest("GET", "/v2/transcript", options)
	if err != nil {
		return TranscriptList{}, err
	}

	vs, err := query.Values(options)
	if err != nil {
		return TranscriptList{}, err
	}
	req.URL.RawQuery = vs.Encode()

	var results TranscriptList

	resp, err := s.client.do(ctx, req, &results)
	if err != nil {
		return TranscriptList{}, err
	}
	defer resp.Body.Close()

	return results, nil
}

// Wait returns once a transcript has completed or failed.
func (s *TranscriptService) Wait(ctx context.Context, transcriptID string) (Transcript, error) {
	b := backoff.NewExponentialBackOff()

	b.InitialInterval = 3 * time.Second

	ticker := backoff.NewTicker(b)

	for {
		select {
		case <-ticker.C:
			ts, err := s.Get(ctx, transcriptID)
			if err != nil {
				return ts, err
			}

			if ts.Status == "completed" || ts.Status == "error" {
				return ts, err
			}
		case <-ctx.Done():
			return Transcript{}, ctx.Err()
		}
	}
}

// TranscribeFromURL submits a URL to an audio file for transcription and waits for it to finish.
func (s *TranscriptService) TranscribeFromURL(ctx context.Context, audioURL string, opts *TranscriptOptionalParams) (Transcript, error) {
	transcript, err := s.SubmitFromURL(ctx, audioURL, opts)
	if err != nil {
		return transcript, err
	}
	return s.Wait(ctx, *transcript.ID)
}

// TranscribeFromReader submits audio for transcription and waits for it to finish.
func (s *TranscriptService) TranscribeFromReader(ctx context.Context, reader io.Reader, opts *TranscriptOptionalParams) (Transcript, error) {
	transcript, err := s.SubmitFromReader(ctx, reader, opts)
	if err != nil {
		return transcript, err
	}
	return s.Wait(ctx, *transcript.ID)
}

// WordSearch searches a transcript for any occurrences of the provided words.
func (s *TranscriptService) WordSearch(ctx context.Context, transcriptID string, words []string) (WordSearchResponse, error) {
	values := url.Values{}
	values.Set("words", strings.Join(words, ","))

	req, err := s.client.newJSONRequest("GET", fmt.Sprint("/v2/transcript/", transcriptID, "/word-search?", values.Encode()), nil)
	if err != nil {
		return WordSearchResponse{}, err
	}

	var results WordSearchResponse

	resp, err := s.client.do(ctx, req, &results)
	if err != nil {
		return WordSearchResponse{}, err
	}
	defer resp.Body.Close()

	return results, nil

}
