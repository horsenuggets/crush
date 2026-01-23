package voice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const whisperAPIURL = "https://api.openai.com/v1/audio/transcriptions"

// BasePrompt contains common technical vocabulary for better transcription accuracy.
const BasePrompt = "Claude, Claude Code, Anthropic, Bubbletea, Lipgloss, Charm, " +
	"GitHub, git, commit, push, pull request, PR, repository, codebase, " +
	"TypeScript, JavaScript, Go, Golang, Rust, Python, API, CLI, terminal"

// TranscriptionResponse represents the response from the Whisper API.
type TranscriptionResponse struct {
	Text string `json:"text"`
}

// APIError represents an error response from the API.
type APIError struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// WhisperClient handles transcription using OpenAI's Whisper API.
type WhisperClient struct {
	apiKey     string
	httpClient *http.Client
}

// NewWhisperClient creates a new Whisper API client.
// If apiKey is empty, it will try to read from OPENAI_API_KEY environment variable.
func NewWhisperClient(apiKey string) (*WhisperClient, error) {
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}

	return &WhisperClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// TranscribeFile transcribes an audio file using the Whisper API.
// dynamicContext is optional session-specific terms to improve transcription accuracy.
func (c *WhisperClient) TranscribeFile(ctx context.Context, audioPath string, dynamicContext string) (string, error) {
	file, err := os.Open(audioPath)
	if err != nil {
		return "", fmt.Errorf("failed to open audio file: %w", err)
	}
	defer file.Close()

	return c.TranscribeReader(ctx, file, audioPath, dynamicContext)
}

// TranscribeReader transcribes audio from an io.Reader using the Whisper API.
func (c *WhisperClient) TranscribeReader(ctx context.Context, audio io.Reader, filename string, dynamicContext string) (string, error) {
	// Build the prompt
	prompt := BasePrompt
	if dynamicContext != "" {
		prompt = fmt.Sprintf("%s, %s", BasePrompt, dynamicContext)
		// Truncate if too long (Whisper prompt limit is ~224 tokens)
		if len(prompt) > 400 {
			prompt = prompt[:400]
		}
	}

	// Create multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add file
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, audio); err != nil {
		return "", fmt.Errorf("failed to copy audio data: %w", err)
	}

	// Add model
	if err := writer.WriteField("model", "whisper-1"); err != nil {
		return "", fmt.Errorf("failed to write model field: %w", err)
	}

	// Add language
	if err := writer.WriteField("language", "en"); err != nil {
		return "", fmt.Errorf("failed to write language field: %w", err)
	}

	// Add prompt
	if err := writer.WriteField("prompt", prompt); err != nil {
		return "", fmt.Errorf("failed to write prompt field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", whisperAPIURL, &buf)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err == nil {
			return "", fmt.Errorf("%s", apiErr.Error.Message)
		}
		return "", fmt.Errorf("API error: %s", string(body))
	}

	var result TranscriptionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Text, nil
}
