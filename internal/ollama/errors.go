// Package ollama provides Ollama server lifecycle management for Crush.
// It also provides generic error handling for local model servers (Ollama, LM Studio, llama.cpp, etc.).
package ollama

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

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
}

func (e *OllamaError) Error() string {
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
	}
}

// WrapError wraps a raw error with user-friendly context if it's a local model server error.
func WrapError(err error, baseURL string) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	// Connection refused - check if it's a local server
	if strings.Contains(errStr, "connection refused") {
		if isLocalServer(baseURL) {
			return NewConnectionError(baseURL)
		}
	}

	// Model not found
	if strings.Contains(errStr, "model") && strings.Contains(errStr, "not found") {
		serverName := detectServerName(baseURL)
		action := "Ensure the model is available on your local server."
		if serverName == "Ollama" {
			action = "Pull the model with: ollama pull <model-name>"
		}
		return &OllamaError{
			Err:     ErrModelNotFound,
			Title:   "Model Not Found",
			Message: "The requested model is not available.",
			Action:  action,
		}
	}

	return err
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
