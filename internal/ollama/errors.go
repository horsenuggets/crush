// Package ollama provides Ollama server lifecycle management for Crush.
// It also provides generic error handling for local model servers (Ollama, LM Studio, llama.cpp, etc.).
package ollama

import (
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
)

// debugLog logs debug messages for error wrapping.
func debugLog(msg string, args ...any) {
	slog.Debug(msg, args...)
}

// Error types for local model operations.
var (
	ErrNotInstalled      = errors.New("ollama is not installed")
	ErrNotRunning        = errors.New("local model server is not running")
	ErrConnectionRefused = errors.New("cannot connect to local model server")
	ErrModelNotFound     = errors.New("model not found")
	ErrStartFailed       = errors.New("failed to start Ollama server")
)

// Common local model server ports for detection.
var localServerPorts = map[string]string{
	"11434": "Ollama",
	"1234":  "LM Studio",
	"8080":  "llama.cpp/LocalAI",
	"5000":  "text-generation-webui",
}

// OllamaError wraps errors with user-friendly messages.
type OllamaError struct {
	Err     error
	Title   string
	Message string
	Action  string
	Model   string // Model name, if applicable (e.g., for ErrModelNotFound)
}

func (e *OllamaError) Error() string {
	if e.Action != "" {
		return e.Message + "\n\n" + e.Action
	}
	return e.Message
}

func (e *OllamaError) Unwrap() error {
	return e.Err
}

// NewNotInstalledError creates an error for when Ollama is not installed.
func NewNotInstalledError() *OllamaError {
	return &OllamaError{
		Err:     ErrNotInstalled,
		Title:   "Ollama Not Installed",
		Message: "Ollama is not installed on this system.",
		Action:  InstallGuidance(),
	}
}

// NewNotRunningError creates an error for when Ollama server is not running.
func NewNotRunningError(autoStartFailed bool) *OllamaError {
	action := "Start Ollama with: ollama serve"
	if autoStartFailed {
		action = "Failed to auto-start Ollama. Start manually with: ollama serve"
	}
	return &OllamaError{
		Err:     ErrNotRunning,
		Title:   "Ollama Not Running",
		Message: "Ollama server is not running.",
		Action:  action,
	}
}

// NewConnectionError creates an error for connection failures.
func NewConnectionError(baseURL string) *OllamaError {
	serverName := detectServerName(baseURL)
	action := "Ensure the local model server is running."

	switch serverName {
	case "Ollama":
		action = "Start Ollama with: ollama serve"
	case "LM Studio":
		action = "Start LM Studio and enable the local server in settings."
	case "llama.cpp/LocalAI":
		action = "Start the llama.cpp server or LocalAI container."
	}

	return &OllamaError{
		Err:     ErrConnectionRefused,
		Title:   serverName + " Connection Failed",
		Message: fmt.Sprintf("Cannot connect to %s at %s", serverName, baseURL),
		Action:  action,
	}
}

// detectServerName guesses the server type from the URL.
func detectServerName(baseURL string) string {
	for port, name := range localServerPorts {
		if strings.Contains(baseURL, ":"+port) {
			return name
		}
	}
	return "Local Model Server"
}

// NewModelNotFoundError creates an error for missing models.
func NewModelNotFoundError(model string) *OllamaError {
	return &OllamaError{
		Err:     ErrModelNotFound,
		Title:   "Model Not Found",
		Message: fmt.Sprintf("Model %q is not available in Ollama.", model),
		Action:  fmt.Sprintf("Pull the model with: ollama pull %s", model),
		Model:   model,
	}
}

// WrapError wraps a raw error with user-friendly context if it's a local model server error.
// It extracts the URL from the error message if baseURL is empty or not a valid URL.
func WrapError(err error, baseURL string) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	// Try to extract URL from error message if not provided or not a valid URL
	if baseURL == "" || !strings.HasPrefix(baseURL, "http") {
		baseURL = extractURLFromError(errStr)
	}

	// Debug logging
	debugLog("WrapError called", "errStr", errStr, "extractedURL", baseURL, "containsConnRefused", strings.Contains(errStr, "connection refused"))

	// Connection refused - check if it's a local server
	if strings.Contains(errStr, "connection refused") {
		if baseURL != "" && isLocalServer(baseURL) {
			return NewConnectionError(baseURL)
		}
		// Also check if the error mentions localhost directly
		if strings.Contains(errStr, "localhost") || strings.Contains(errStr, "127.0.0.1") {
			// Extract port to determine server type
			url := extractURLFromError(errStr)
			if url != "" {
				return NewConnectionError(url)
			}
			return NewConnectionError("localhost")
		}
	}

	// Model not found
	if strings.Contains(errStr, "model") && strings.Contains(errStr, "not found") {
		// Try to extract model name from error message
		modelName := extractModelName(errStr)
		serverName := detectServerName(baseURL)
		// Default to Ollama for unknown local servers or when baseURL is empty
		if serverName == "Local Model Server" || serverName == "" {
			serverName = "Ollama"
		}

		action := "Ensure the model is available on your local server."
		if serverName == "Ollama" {
			if modelName != "" {
				action = fmt.Sprintf("Pull the model with: ollama pull %s", modelName)
			} else {
				action = "Pull the model with: ollama pull <model-name>"
			}
		}
		return &OllamaError{
			Err:     ErrModelNotFound,
			Title:   "Model Not Found",
			Message: "The requested model is not available.",
			Action:  action,
			Model:   modelName,
		}
	}

	return err
}

// extractModelName tries to extract a model name from an error message.
func extractModelName(errStr string) string {
	// Look for patterns like: model "llama3.2" not found
	// or: model 'llama3.2' not found
	patterns := []string{`"`, `'`}
	for _, quote := range patterns {
		start := strings.Index(errStr, "model "+quote)
		if start >= 0 {
			start += len("model " + quote)
			end := strings.Index(errStr[start:], quote)
			if end > 0 {
				return errStr[start : start+end]
			}
		}
	}
	return ""
}

// extractURLFromError tries to extract a URL from an error message.
func extractURLFromError(errStr string) string {
	// Look for patterns like: Post "http://localhost:11434/..."
	// or: dial tcp localhost:11434

	// Check for quoted URL
	if idx := strings.Index(errStr, "\"http"); idx >= 0 {
		end := strings.Index(errStr[idx+1:], "\"")
		if end > 0 {
			return errStr[idx+1 : idx+1+end]
		}
	}

	// Check for localhost with port in dial error
	for port := range localServerPorts {
		if strings.Contains(errStr, "localhost:"+port) || strings.Contains(errStr, "127.0.0.1:"+port) {
			return "http://localhost:" + port
		}
	}

	return ""
}

// isLocalServer checks if the URL points to a local server.
func isLocalServer(baseURL string) bool {
	return strings.Contains(baseURL, "localhost") ||
		strings.Contains(baseURL, "127.0.0.1") ||
		strings.Contains(baseURL, "0.0.0.0")
}

// IsOllamaURL checks if the URL is for an Ollama server.
func IsOllamaURL(baseURL string) bool {
	return isLocalServer(baseURL) && strings.Contains(baseURL, ":11434")
}

// InstallGuidance returns platform-specific installation instructions.
func InstallGuidance() string {
	switch runtime.GOOS {
	case "darwin":
		return `Install Ollama on macOS:
  brew install ollama

Or download from: https://ollama.ai/download/mac`

	case "linux":
		return `Install Ollama on Linux:
  curl -fsSL https://ollama.ai/install.sh | sh

Or visit: https://ollama.ai/download/linux`

	case "windows":
		return `Install Ollama on Windows:
  Download from: https://ollama.ai/download/windows`

	default:
		return "Install Ollama from: https://ollama.ai/download"
	}
}
