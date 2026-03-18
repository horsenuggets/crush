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
	"strings"
	"time"

	"github.com/charmbracelet/crush/internal/config"
)

const whisperAPIURL = "https://api.openai.com/v1/audio/transcriptions"

// BasePrompt contains common technical vocabulary for better transcription accuracy.
const BasePrompt = "Claude, Claude Code, Anthropic, " +
	"GitHub, git, commit, push, pull request, PR, repository, codebase, " +
	"API, CLI, terminal"

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

// MinAudioDuration is the minimum recording duration to avoid hallucinations.
const MinAudioDuration = 100 * time.Millisecond

// MinAudioFileSize is the minimum file size for a valid recording (16kHz mono 16-bit = 32KB/sec).
// 100ms of audio = ~3.2KB minimum.
const MinAudioFileSize = 3200

// TranscribeFile transcribes an audio file using the Whisper API.
// dynamicContext is optional session-specific terms to improve transcription accuracy.
func (c *WhisperClient) TranscribeFile(ctx context.Context, audioPath string, dynamicContext string) (string, error) {
	// Check file size to reject very short/empty recordings
	fileInfo, err := os.Stat(audioPath)
	if err != nil {
		return "", fmt.Errorf("failed to stat audio file: %w", err)
	}
	if fileInfo.Size() < MinAudioFileSize {
		return "", fmt.Errorf("recording too short (less than 100ms)")
	}

	file, err := os.Open(audioPath)
	if err != nil {
		return "", fmt.Errorf("failed to open audio file: %w", err)
	}
	defer file.Close()

	return c.TranscribeReader(ctx, file, audioPath, dynamicContext)
}

// loadContextVocabulary loads additional vocabulary from the configured context file.
func loadContextVocabulary() string {
	cfg := config.Get()
	if cfg.Options == nil || cfg.Options.TUI == nil || cfg.Options.TUI.Voice == nil {
		return ""
	}
	contextFile := cfg.Options.TUI.Voice.ContextFile
	if contextFile == "" {
		return ""
	}
	// Expand ~ if present
	if strings.HasPrefix(contextFile, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			contextFile = strings.Replace(contextFile, "~", home, 1)
		}
	}
	content, err := os.ReadFile(contextFile)
	if err != nil {
		return ""
	}
	// Read file line by line and join with commas
	lines := strings.Split(string(content), "\n")
	var terms []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		terms = append(terms, line)
	}
	return strings.Join(terms, ", ")
}

// TranscribeReader transcribes audio from an io.Reader using the Whisper API.
func (c *WhisperClient) TranscribeReader(ctx context.Context, audio io.Reader, filename string, dynamicContext string) (string, error) {
	// Build the prompt with base vocabulary
	prompt := BasePrompt

	// Add vocabulary from context file if configured
	contextVocab := loadContextVocabulary()
	if contextVocab != "" {
		prompt = fmt.Sprintf("%s, %s", prompt, contextVocab)
	}

	// Add dynamic context from the session
	if dynamicContext != "" {
		prompt = fmt.Sprintf("%s, %s", prompt, dynamicContext)
	}

	// Truncate if too long (Whisper prompt limit is ~224 tokens)
	if len(prompt) > 400 {
		prompt = prompt[:400]
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

	// Add temperature=0 for deterministic output (reduces hallucination)
	if err := writer.WriteField("temperature", "0"); err != nil {
		return "", fmt.Errorf("failed to write temperature field: %w", err)
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

	// Check for common hallucination patterns
	if isLikelyHallucination(result.Text) {
		return "", nil // Return empty string instead of hallucinated text
	}

	return result.Text, nil
}

// commonHallucinations contains phrases commonly hallucinated by Whisper
// when processing silence, background noise, or very short audio.
var commonHallucinations = []string{
	"thank you for watching",
	"thanks for watching",
	"please subscribe",
	"like and subscribe",
	"see you in the next",
	"don't forget to",
	"comment below",
	"hit the bell",
	"notification",
	"subtitles by",
	"captions by",
	"transcribed by",
	"music playing",
	"[music]",
	"[applause]",
	"[laughter]",
	"you",
	"bye",
	"bye.",
	"goodbye",
	"goodbye.",
}

// isLikelyHallucination checks if the transcription is likely a hallucination.
func isLikelyHallucination(text string) bool {
	text = strings.ToLower(strings.TrimSpace(text))

	// Empty or very short text
	if len(text) < 2 {
		return true
	}

	// Check for common hallucination patterns
	for _, pattern := range commonHallucinations {
		if strings.Contains(text, pattern) {
			return true
		}
	}

	// Single-word responses that are common hallucinations
	words := strings.Fields(text)
	if len(words) == 1 {
		singleWord := strings.ToLower(words[0])
		// Common single-word hallucinations
		if singleWord == "you" || singleWord == "bye" || singleWord == "goodbye" ||
			singleWord == "thanks" || singleWord == "okay" || singleWord == "yes" ||
			singleWord == "no" || singleWord == "um" || singleWord == "uh" ||
			singleWord == "hmm" || singleWord == "ah" {
			return true
		}
	}

	return false
}
