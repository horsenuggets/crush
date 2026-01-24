package ollama

import (
	"errors"
	"fmt"
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
	// Test connection refused error with explicit URL
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

	// Test connection refused error with URL extracted from error message
	// This simulates when WrapError is called with provider type instead of URL
	err = WrapError(
		&testError{msg: `Post "http://localhost:11434/v1/chat/completions": dial tcp [::1]:11434: connect: connection refused`},
		"openaicompat", // Provider type, not URL
	)
	if err == nil {
		t.Error("Expected wrapped error, got nil")
	}
	ollamaErr, ok = err.(*OllamaError)
	if !ok {
		t.Errorf("Expected OllamaError when URL is in error message, got %T", err)
	} else {
		if ollamaErr.Title != "Ollama Connection Failed" {
			t.Errorf("Expected title 'Ollama Connection Failed', got %q", ollamaErr.Title)
		}
	}

	// Test LM Studio connection refused
	err = WrapError(
		&testError{msg: `Post "http://localhost:1234/v1/chat/completions": dial tcp localhost:1234: connect: connection refused`},
		"",
	)
	if err == nil {
		t.Error("Expected wrapped error for LM Studio, got nil")
	}
	ollamaErr, ok = err.(*OllamaError)
	if !ok {
		t.Errorf("Expected OllamaError for LM Studio, got %T", err)
	} else {
		if ollamaErr.Title != "LM Studio Connection Failed" {
			t.Errorf("Expected title 'LM Studio Connection Failed', got %q", ollamaErr.Title)
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

func TestCanAutoInstall(t *testing.T) {
	canInstall := CanAutoInstall()
	t.Logf("Can auto-install: %v", canInstall)
	t.Logf("Install method: %s", InstallMethod())
}

// TestExactProviderErrorScenario tests the exact scenario where a fantasy.ProviderError
// wraps a connection refused error - this is what actually happens in production.
func TestExactProviderErrorScenario(t *testing.T) {
	// This is the EXACT error message format from the screenshot
	errMsg := `Post "http://localhost:11434/v1/chat/completions": dial tcp [::1]:11434: connect: connection refused`

	// Simulate fantasy.ProviderError (which just returns Message from Error())
	providerErr := &mockProviderError{message: errMsg}

	t.Logf("Provider error message: %s", providerErr.Error())

	// Call WrapError with empty baseURL (as done in agent.go)
	wrapped := WrapError(providerErr, "")

	t.Logf("Wrapped error type: %T", wrapped)
	t.Logf("Wrapped error: %s", wrapped.Error())

	// Check if it's an OllamaError
	var ollamaErr *OllamaError
	if !errors.As(wrapped, &ollamaErr) {
		t.Fatalf("FAILED: WrapError did not return OllamaError for connection refused to localhost:11434")
	}

	t.Logf("SUCCESS: Found OllamaError")
	t.Logf("Title: %s", ollamaErr.Title)
	t.Logf("Message: %s", ollamaErr.Message)
	t.Logf("Action: %s", ollamaErr.Action)

	if ollamaErr.Title != "Ollama Connection Failed" {
		t.Errorf("Expected title 'Ollama Connection Failed', got %q", ollamaErr.Title)
	}
}

// mockProviderError simulates fantasy.ProviderError which returns Message from Error()
type mockProviderError struct {
	message string
	title   string
}

func (e *mockProviderError) Error() string {
	if e.title == "" {
		return e.message
	}
	return e.title + ": " + e.message
}

// TestWrapErrorWithWrappedErrors tests that WrapError correctly handles errors
// that are wrapped by other error types (like fantasy.ProviderError).
// This simulates the real integration path where HTTP errors are wrapped.
func TestWrapErrorWithWrappedErrors(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		baseURL        string
		expectOllama   bool
		expectTitle    string
		expectContains string
	}{
		{
			name:           "wrapped connection refused to Ollama",
			err:            fmt.Errorf("provider error: %w", &testError{msg: `Post "http://localhost:11434/v1/chat/completions": dial tcp [::1]:11434: connect: connection refused`}),
			baseURL:        "",
			expectOllama:   true,
			expectTitle:    "Ollama Connection Failed",
			expectContains: "ollama serve",
		},
		{
			name:           "double wrapped connection refused",
			err:            fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", &testError{msg: `Post "http://localhost:11434/v1/chat/completions": dial tcp [::1]:11434: connect: connection refused`})),
			baseURL:        "",
			expectOllama:   true,
			expectTitle:    "Ollama Connection Failed",
			expectContains: "ollama serve",
		},
		{
			name:           "wrapped LM Studio error",
			err:            fmt.Errorf("request failed: %w", &testError{msg: `Post "http://localhost:1234/v1/chat/completions": dial tcp localhost:1234: connect: connection refused`}),
			baseURL:        "",
			expectOllama:   true,
			expectTitle:    "LM Studio Connection Failed",
			expectContains: "LM Studio",
		},
		{
			name:           "wrapped llama.cpp error",
			err:            fmt.Errorf("request failed: %w", &testError{msg: `Post "http://localhost:8080/v1/chat/completions": dial tcp localhost:8080: connect: connection refused`}),
			baseURL:        "",
			expectOllama:   true,
			expectTitle:    "llama.cpp/LocalAI Connection Failed",
			expectContains: "llama.cpp",
		},
		{
			name:         "wrapped remote API error should not be Ollama",
			err:          fmt.Errorf("request failed: %w", &testError{msg: `Post "https://api.openai.com/v1/chat/completions": dial tcp: connection refused`}),
			baseURL:      "",
			expectOllama: false,
		},
		{
			name:           "model not found error",
			err:            fmt.Errorf("provider error: %w", &testError{msg: `model "llama3" not found`}),
			baseURL:        "http://localhost:11434/v1",
			expectOllama:   true,
			expectTitle:    "Model Not Found",
			expectContains: "ollama pull",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapError(tt.err, tt.baseURL)

			var ollamaErr *OllamaError
			isOllamaErr := errors.As(result, &ollamaErr)

			if tt.expectOllama {
				if !isOllamaErr {
					t.Errorf("Expected OllamaError, got %T: %v", result, result)
					return
				}
				if ollamaErr.Title != tt.expectTitle {
					t.Errorf("Expected title %q, got %q", tt.expectTitle, ollamaErr.Title)
				}
				errMsg := ollamaErr.Error()
				if tt.expectContains != "" && !contains(errMsg, tt.expectContains) {
					t.Errorf("Expected error message to contain %q, got %q", tt.expectContains, errMsg)
				}
			} else {
				if isOllamaErr {
					t.Errorf("Did not expect OllamaError, got %v", ollamaErr)
				}
			}
		})
	}
}

// TestOllamaErrorUnwrap tests that OllamaError properly implements error unwrapping
// so errors.Is and errors.As work correctly through the error chain.
func TestOllamaErrorUnwrap(t *testing.T) {
	// Test that wrapped OllamaError can be detected with errors.Is
	notInstalledErr := NewNotInstalledError()
	wrappedErr := fmt.Errorf("failed to update models: %w", notInstalledErr)
	doubleWrapped := fmt.Errorf("coordinator error: %w", wrappedErr)

	if !errors.Is(doubleWrapped, ErrNotInstalled) {
		t.Error("errors.Is should find ErrNotInstalled through wrapped OllamaError")
	}

	// Test errors.As with wrapped OllamaError
	var ollamaErr *OllamaError
	if !errors.As(doubleWrapped, &ollamaErr) {
		t.Error("errors.As should find OllamaError through wrapping")
	}
	if ollamaErr.Title != "Ollama Not Installed" {
		t.Errorf("Expected title 'Ollama Not Installed', got %q", ollamaErr.Title)
	}
}

// TestOllamaErrorMessage tests that OllamaError.Error() includes both message and action.
func TestOllamaErrorMessage(t *testing.T) {
	err := NewNotInstalledError()
	msg := err.Error()

	if !contains(msg, "not installed") {
		t.Errorf("Error message should contain 'not installed', got: %s", msg)
	}
	if !contains(msg, "brew install") || !contains(msg, "ollama.ai") {
		t.Errorf("Error message should contain install instructions, got: %s", msg)
	}
}

// TestAllLocalServerPorts tests that all known local server ports are detected.
func TestAllLocalServerPorts(t *testing.T) {
	ports := map[string]string{
		"11434": "Ollama",
		"1234":  "LM Studio",
		"8080":  "llama.cpp/LocalAI",
		"5000":  "text-generation-webui",
	}

	for port, expectedName := range ports {
		t.Run(expectedName, func(t *testing.T) {
			errMsg := fmt.Sprintf(`Post "http://localhost:%s/v1/chat/completions": dial tcp localhost:%s: connect: connection refused`, port, port)
			result := WrapError(&testError{msg: errMsg}, "")

			var ollamaErr *OllamaError
			if !errors.As(result, &ollamaErr) {
				t.Errorf("Expected OllamaError for %s (port %s), got %T", expectedName, port, result)
				return
			}

			expectedTitle := expectedName + " Connection Failed"
			if ollamaErr.Title != expectedTitle {
				t.Errorf("Expected title %q, got %q", expectedTitle, ollamaErr.Title)
			}
		})
	}
}

// TestExtractModelName tests model name extraction from various error message formats.
func TestExtractModelName(t *testing.T) {
	tests := []struct {
		name     string
		errStr   string
		expected string
	}{
		{
			name:     "single quotes",
			errStr:   "not found: model 'llama3.2' not found",
			expected: "llama3.2",
		},
		{
			name:     "double quotes",
			errStr:   `model "llama3.2" not found in Ollama`,
			expected: "llama3.2",
		},
		{
			name:     "no quotes",
			errStr:   "model not found",
			expected: "",
		},
		{
			name:     "model with colon",
			errStr:   "model 'deepseek-coder:6.7b' not found",
			expected: "deepseek-coder:6.7b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractModelName(tt.errStr)
			if result != tt.expected {
				t.Errorf("extractModelName(%q) = %q, want %q", tt.errStr, result, tt.expected)
			}
		})
	}
}

// TestExtractURLFromError tests URL extraction from various error message formats.
func TestExtractURLFromError(t *testing.T) {
	tests := []struct {
		name     string
		errStr   string
		expected string
	}{
		{
			name:     "quoted HTTP URL",
			errStr:   `Post "http://localhost:11434/v1/chat/completions": dial tcp`,
			expected: "http://localhost:11434/v1/chat/completions",
		},
		{
			name:     "localhost with port in dial error",
			errStr:   "dial tcp localhost:11434: connect: connection refused",
			expected: "http://localhost:11434",
		},
		{
			name:     "127.0.0.1 with port",
			errStr:   "dial tcp 127.0.0.1:1234: connect: connection refused",
			expected: "http://localhost:1234",
		},
		{
			name:     "no URL in error",
			errStr:   "some other error without URL",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractURLFromError(tt.errStr)
			if result != tt.expected {
				t.Errorf("extractURLFromError(%q) = %q, want %q", tt.errStr, result, tt.expected)
			}
		})
	}
}

// TestIsLocalServer tests local server detection.
func TestIsLocalServer(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"http://localhost:11434", true},
		{"http://127.0.0.1:11434", true},
		{"http://0.0.0.0:11434", true},
		{"https://api.openai.com", false},
		{"https://api.anthropic.com", false},
		{"http://192.168.1.100:11434", false}, // Local network but not localhost
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			result := isLocalServer(tt.url)
			if result != tt.expected {
				t.Errorf("isLocalServer(%q) = %v, want %v", tt.url, result, tt.expected)
			}
		})
	}
}

// TestNewErrorFunctions tests all error constructor functions.
func TestNewErrorFunctions(t *testing.T) {
	t.Run("NewNotInstalledError", func(t *testing.T) {
		err := NewNotInstalledError()
		if !errors.Is(err, ErrNotInstalled) {
			t.Error("Should wrap ErrNotInstalled")
		}
		if err.Title != "Ollama Not Installed" {
			t.Errorf("Wrong title: %s", err.Title)
		}
	})

	t.Run("NewNotRunningError without auto-start failure", func(t *testing.T) {
		err := NewNotRunningError(false)
		if !errors.Is(err, ErrNotRunning) {
			t.Error("Should wrap ErrNotRunning")
		}
		if contains(err.Action, "Failed to auto-start") {
			t.Error("Should not mention auto-start failure")
		}
	})

	t.Run("NewNotRunningError with auto-start failure", func(t *testing.T) {
		err := NewNotRunningError(true)
		if !contains(err.Action, "Failed to auto-start") {
			t.Error("Should mention auto-start failure")
		}
	})

	t.Run("NewConnectionError for Ollama", func(t *testing.T) {
		err := NewConnectionError("http://localhost:11434/v1")
		if !errors.Is(err, ErrConnectionRefused) {
			t.Error("Should wrap ErrConnectionRefused")
		}
		if err.Title != "Ollama Connection Failed" {
			t.Errorf("Wrong title: %s", err.Title)
		}
		if !contains(err.Action, "ollama serve") {
			t.Error("Should suggest ollama serve")
		}
	})

	t.Run("NewConnectionError for LM Studio", func(t *testing.T) {
		err := NewConnectionError("http://localhost:1234/v1")
		if err.Title != "LM Studio Connection Failed" {
			t.Errorf("Wrong title: %s", err.Title)
		}
		if !contains(err.Action, "LM Studio") {
			t.Error("Should mention LM Studio")
		}
	})

	t.Run("NewModelNotFoundError", func(t *testing.T) {
		err := NewModelNotFoundError("llama3.2")
		if !errors.Is(err, ErrModelNotFound) {
			t.Error("Should wrap ErrModelNotFound")
		}
		if !contains(err.Message, "llama3.2") {
			t.Error("Should include model name in message")
		}
		if !contains(err.Action, "ollama pull llama3.2") {
			t.Error("Should suggest pulling the specific model")
		}
		if err.Model != "llama3.2" {
			t.Errorf("Should set Model field, got: %s", err.Model)
		}
	})

	t.Run("WrapError extracts model name", func(t *testing.T) {
		rawErr := &testError{msg: `model "llama3.2" not found`}
		wrapped := WrapError(rawErr, "")

		var ollamaErr *OllamaError
		if !errors.As(wrapped, &ollamaErr) {
			t.Fatal("Should wrap as OllamaError")
		}
		if ollamaErr.Model != "llama3.2" {
			t.Errorf("Should extract model name, got: %s", ollamaErr.Model)
		}
		if !contains(ollamaErr.Action, "ollama pull llama3.2") {
			t.Error("Should suggest pulling the specific model")
		}
	})

	t.Run("WrapError handles model not found with single quotes", func(t *testing.T) {
		rawErr := &testError{msg: `model 'gemma2' not found`}
		wrapped := WrapError(rawErr, "")

		var ollamaErr *OllamaError
		if !errors.As(wrapped, &ollamaErr) {
			t.Fatal("Should wrap as OllamaError")
		}
		if ollamaErr.Model != "gemma2" {
			t.Errorf("Should extract model name from single quotes, got: %s", ollamaErr.Model)
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && searchString(s, substr)))
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
