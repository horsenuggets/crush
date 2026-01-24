package ollama

import (
	"testing"
)

func TestIsInstalled(t *testing.T) {
	// This tests the actual system state
	installed := IsInstalled()
	t.Logf("Ollama installed: %v", installed)
}

func TestIsServerRunning(t *testing.T) {
	// This tests the actual system state
	running := IsServerRunning()
	t.Logf("Ollama server running: %v", running)
}

func TestInstallGuidance(t *testing.T) {
	guidance := InstallGuidance()
	if guidance == "" {
		t.Error("InstallGuidance should not be empty")
	}
	t.Logf("Install guidance:\n%s", guidance)
}

func TestIsOllamaURL(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"http://localhost:11434/v1", true},
		{"http://127.0.0.1:11434/v1", true},
		{"http://localhost:1234/v1", false},  // LM Studio
		{"http://localhost:8080/v1", false},  // llama.cpp
		{"https://api.openai.com/v1", false}, // Remote
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			result := IsOllamaURL(tt.url)
			if result != tt.expected {
				t.Errorf("IsOllamaURL(%q) = %v, want %v", tt.url, result, tt.expected)
			}
		})
	}
}

func TestDetectServerName(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"http://localhost:11434/v1", "Ollama"},
		{"http://localhost:1234/v1", "LM Studio"},
		{"http://localhost:8080/v1", "llama.cpp/LocalAI"},
		{"http://localhost:5000/v1", "text-generation-webui"},
		{"http://localhost:9999/v1", "Local Model Server"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			result := detectServerName(tt.url)
			if result != tt.expected {
				t.Errorf("detectServerName(%q) = %q, want %q", tt.url, result, tt.expected)
			}
		})
	}
}

func TestWrapError(t *testing.T) {
	// Test connection refused error
	err := WrapError(
		&testError{msg: "dial tcp [::1]:11434: connect: connection refused"},
		"http://localhost:11434/v1",
	)
	if err == nil {
		t.Error("Expected wrapped error, got nil")
	}
	ollamaErr, ok := err.(*OllamaError)
	if !ok {
		t.Errorf("Expected OllamaError, got %T", err)
	} else {
		if ollamaErr.Title != "Ollama Connection Failed" {
			t.Errorf("Expected title 'Ollama Connection Failed', got %q", ollamaErr.Title)
		}
	}

	// Test that non-local errors are not wrapped
	err = WrapError(
		&testError{msg: "dial tcp: connection refused"},
		"https://api.openai.com/v1",
	)
	if _, ok := err.(*OllamaError); ok {
		t.Error("Should not wrap non-local server errors")
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
