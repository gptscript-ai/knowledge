package assemblyai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var (
	// ErrSessionClosed is returned when attempting to write to a closed session.
	ErrSessionClosed = errors.New("session closed")

	// ErrDisconnected is returned when attempting to write to a disconnected client.
	ErrDisconnected = errors.New("client is disconnected")
)

type MessageType string

const (
	MessageTypeSessionBegins     MessageType = "SessionBegins"
	MessageTypeSessionTerminated MessageType = "SessionTerminated"
	MessageTypePartialTranscript MessageType = "PartialTranscript"
	MessageTypeFinalTranscript   MessageType = "FinalTranscript"
)

type AudioData struct {
	// Base64 encoded raw audio data
	AudioData string `json:"audio_data,omitempty"`
}

type TerminateSession struct {
	// Set to true to end your real-time session forever
	TerminateSession bool `json:"terminate_session"`
}

type endUtteranceSilenceThreshold struct {
	// Set to true to configure the silence threshold for ending utterances.
	EndUtteranceSilenceThreshold int64 `json:"end_utterance_silence_threshold"`
}

type forceEndUtterance struct {
	// Set to true to manually end the current utterance.
	ForceEndUtterance bool `json:"force_end_utterance"`
}

type RealTimeBaseMessage struct {
	// Describes the type of the message
	MessageType MessageType `json:"message_type"`
}

type RealTimeBaseTranscript struct {
	// End time of audio sample relative to session start, in milliseconds
	AudioEnd int64 `json:"audio_end"`

	// Start time of audio sample relative to session start, in milliseconds
	AudioStart int64 `json:"audio_start"`

	// The confidence score of the entire transcription, between 0 and 1
	Confidence float64 `json:"confidence"`

	// The timestamp for the partial transcript
	Created string `json:"created"`

	// The partial transcript for your audio
	Text string `json:"text"`

	// An array of objects, with the information for each word in the transcription text.
	// Includes the start and end time of the word in milliseconds, the confidence score of the word, and the text, which is the word itself.
	Words []Word `json:"words"`
}

type FinalTranscript struct {
	RealTimeBaseTranscript

	// Describes the type of message
	MessageType MessageType `json:"message_type"`

	// Whether the text is punctuated and cased
	Punctuated bool `json:"punctuated"`

	// Whether the text is formatted, for example Dollar -> $
	TextFormatted bool `json:"text_formatted"`
}

type PartialTranscript struct {
	RealTimeBaseTranscript

	// Describes the type of message
	MessageType MessageType `json:"message_type"`
}

type Word struct {
	// Confidence score of the word
	Confidence float64 `json:"confidence"`

	// End time of the word in milliseconds
	End int64 `json:"end"`

	// Start time of the word in milliseconds
	Start int64 `json:"start"`

	// The word itself
	Text string `json:"text"`
}

var DefaultSampleRate = 16_000

type RealTimeClient struct {
	baseURL *url.URL
	apiKey  string

	conn *websocket.Conn

	mtx         sync.RWMutex
	sessionOpen bool

	// done is used to clean up resources when the client disconnects.
	done chan bool

	handler RealTimeHandler
}

func (c *RealTimeClient) isSessionOpen() bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return c.sessionOpen
}

func (c *RealTimeClient) setSessionOpen(open bool) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	c.sessionOpen = open
}

type RealTimeError struct {
	Error string `json:"error"`
}

type RealTimeClientOption func(*RealTimeClient)

// WithRealTimeBaseURL sets the API endpoint used by the client. Mainly used for testing.
func WithRealTimeBaseURL(rawurl string) RealTimeClientOption {
	return func(c *RealTimeClient) {
		if u, err := url.Parse(rawurl); err == nil {
			c.baseURL = u
		}
	}
}

func WithRealTimeAPIKey(apiKey string) RealTimeClientOption {
	return func(rtc *RealTimeClient) {
		rtc.apiKey = apiKey
	}
}

func WithHandler(handler RealTimeHandler) RealTimeClientOption {
	return func(rtc *RealTimeClient) {
		rtc.handler = handler
	}
}

func WithRealTimeSampleRate(sampleRate int) RealTimeClientOption {
	return func(rtc *RealTimeClient) {
		if sampleRate > 0 {
			vs := rtc.baseURL.Query()
			vs.Set("sample_rate", strconv.Itoa(sampleRate))
			rtc.baseURL.RawQuery = vs.Encode()
		}
	}
}

func WithRealTimeWordBoost(wordBoost []string) RealTimeClientOption {
	return func(rtc *RealTimeClient) {
		vs := rtc.baseURL.Query()

		if len(wordBoost) > 0 {
			b, _ := json.Marshal(wordBoost)
			vs.Set("word_boost", string(b))
		}

		rtc.baseURL.RawQuery = vs.Encode()
	}
}

// RealTimeEncoding is the encoding format for the audio data.
type RealTimeEncoding string

const (
	// PCM signed 16-bit little-endian (default)
	RealTimeEncodingPCMS16LE RealTimeEncoding = "pcm_s16le"

	// PCM Mu-law
	RealTimeEncodingPCMMulaw RealTimeEncoding = "pcm_mulaw"
)

func WithRealTimeEncoding(encoding RealTimeEncoding) RealTimeClientOption {
	return func(rtc *RealTimeClient) {
		vs := rtc.baseURL.Query()
		vs.Set("encoding", string(encoding))
		rtc.baseURL.RawQuery = vs.Encode()
	}
}

func NewRealTimeClientWithOptions(options ...RealTimeClientOption) *RealTimeClient {
	client := &RealTimeClient{
		baseURL: &url.URL{
			Scheme:   "wss",
			Host:     "api.assemblyai.com",
			Path:     "/v2/realtime/ws",
			RawQuery: fmt.Sprintf("sample_rate=%v", DefaultSampleRate),
		},
	}

	for _, option := range options {
		option(client)
	}

	return client
}

type SessionBegins struct {
	RealTimeBaseMessage
	// Timestamp when this session will expire
	ExpiresAt string `json:"expires_at"`

	// Describes the type of the message
	MessageType string `json:"message_type"`

	// Unique identifier for the established session
	SessionID string `json:"session_id"`
}

type SessionTerminated struct {
	// Describes the type of the message
	MessageType MessageType `json:"message_type"`
}

type RealTimeHandler interface {
	SessionBegins(ev SessionBegins)
	SessionTerminated(ev SessionTerminated)
	FinalTranscript(transcript FinalTranscript)
	PartialTranscript(transcript PartialTranscript)
	Error(err error)
}

func NewRealTimeClient(apiKey string, handler RealTimeHandler) *RealTimeClient {
	return NewRealTimeClientWithOptions(WithRealTimeAPIKey(apiKey), WithHandler(handler))
}

// Connects opens a WebSocket connection and waits for a session to begin.
// Closes the any open WebSocket connection in case of errors.
func (c *RealTimeClient) Connect(ctx context.Context) error {
	header := make(http.Header)
	header.Set("Authorization", c.apiKey)

	opts := &websocket.DialOptions{
		HTTPHeader: header,
		HTTPClient: &http.Client{},
	}

	conn, _, err := websocket.Dial(ctx, c.baseURL.String(), opts)
	if err != nil {
		return err
	}

	c.conn = conn

	var msg json.RawMessage
	if err := wsjson.Read(ctx, c.conn, &msg); err != nil {
		return err
	}

	var realtimeError RealTimeError
	if err := json.Unmarshal(msg, &realtimeError); err != nil {
		return err
	}
	if realtimeError.Error != "" {
		return errors.New(realtimeError.Error)
	}

	var session SessionBegins
	if err := json.Unmarshal(msg, &session); err != nil {
		return err
	}

	c.setSessionOpen(true)
	c.handler.SessionBegins(session)

	c.done = make(chan bool)

	go func() {
		for {
			if !c.isSessionOpen() {
				return
			}

			var msg json.RawMessage

			if err := wsjson.Read(ctx, c.conn, &msg); err != nil {
				c.handler.Error(err)
				return
			}

			var messageType struct {
				MessageType MessageType `json:"message_type"`
			}

			if err := json.Unmarshal(msg, &messageType); err != nil {
				c.handler.Error(err)
				return
			}

			switch messageType.MessageType {
			case MessageTypeFinalTranscript:
				var transcript FinalTranscript
				if err := json.Unmarshal(msg, &transcript); err != nil {
					c.handler.Error(err)
					continue
				}

				if transcript.Text != "" {
					c.handler.FinalTranscript(transcript)
				}
			case MessageTypePartialTranscript:
				var transcript PartialTranscript
				if err := json.Unmarshal(msg, &transcript); err != nil {
					c.handler.Error(err)
					continue
				}

				if transcript.Text != "" {
					c.handler.PartialTranscript(transcript)
				}
			case MessageTypeSessionTerminated:
				var session SessionTerminated
				if err := json.Unmarshal(msg, &session); err != nil {
					c.handler.Error(err)
					continue
				}

				c.setSessionOpen(false)

				c.handler.SessionTerminated(session)

				c.done <- true
			}
		}
	}()

	return nil
}

// Disconnect sends the terminate_session message and waits for the server to
// send a SessionTerminated message before closing the connection.
func (c *RealTimeClient) Disconnect(ctx context.Context, waitForSessionTermination bool) error {
	terminate := TerminateSession{TerminateSession: true}

	if err := wsjson.Write(ctx, c.conn, terminate); err != nil {
		return err
	}

	if waitForSessionTermination {
		<-c.done
	}

	return c.conn.Close(websocket.StatusNormalClosure, "")
}

// Send sends audio samples to be transcribed.
//
// Expected audio format:
//
// - 16-bit signed integers
// - PCM-encoded
// - Single-channel
func (c *RealTimeClient) Send(ctx context.Context, samples []byte) error {
	if c.conn == nil || !c.isSessionOpen() {
		return ErrSessionClosed
	}

	return c.conn.Write(ctx, websocket.MessageBinary, samples)
}

// ForceEndUtterance manually ends an utterance.
func (c *RealTimeClient) ForceEndUtterance(ctx context.Context) error {
	return wsjson.Write(ctx, c.conn, forceEndUtterance{
		ForceEndUtterance: true,
	})
}

// SetEndUtteranceSilenceThreshold configures the threshold for how long to wait
// before ending an utterance. Default is 700ms.
func (c *RealTimeClient) SetEndUtteranceSilenceThreshold(ctx context.Context, threshold int64) error {
	return wsjson.Write(ctx, c.conn, endUtteranceSilenceThreshold{
		EndUtteranceSilenceThreshold: threshold,
	})
}
