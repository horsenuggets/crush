package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	stateFileName    = "ollama-state.json"
	lockFileName     = ".ollama.lock"
	staleTimeout     = 5 * time.Minute
	serverStartWait  = 15 * time.Second
	healthCheckURL   = "http://localhost:11434/api/version"
	defaultBaseURL   = "http://localhost:11434"
)

// State tracks Ollama server usage across Crush instances.
type State struct {
	Instances  map[string]Instance `json:"instances"`
	ServerPID  int                 `json:"server_pid,omitempty"`
	StartedBy  string              `json:"started_by,omitempty"`
	IsExternal bool                `json:"is_external"`
	Updated    time.Time           `json:"updated"`
}

// Instance represents a Crush instance using Ollama.
type Instance struct {
	ID           string    `json:"id"`
	PID          int       `json:"pid"`
	Model        string    `json:"model,omitempty"`
	LastActivity time.Time `json:"last_activity"`
}

// Manager handles Ollama server lifecycle across Crush instances.
type Manager struct {
	mu           sync.Mutex
	instanceID   string
	stateDir     string
	stateFile    string
	lockFilePath string
	active       bool
	modelInUse   string
	serverCmd    *exec.Cmd
}

// NewManager creates a new Ollama lifecycle manager.
func NewManager(instanceID, stateDir string) *Manager {
	return &Manager{
		instanceID:   instanceID,
		stateDir:     stateDir,
		stateFile:    filepath.Join(stateDir, stateFileName),
		lockFilePath: filepath.Join(stateDir, lockFileName),
	}
}

// EnsureRunning ensures Ollama is installed and running.
// Returns nil if ready, or an OllamaError with guidance.
func (m *Manager) EnsureRunning(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if Ollama is installed
	if !IsInstalled() {
		return NewNotInstalledError()
	}

	// Check if server is already running
	if IsServerRunning() {
		// Check if we started it or it's external
		state, err := m.readState()
		if err != nil {
			slog.Debug("Failed to read Ollama state", "error", err)
		}

		if state.ServerPID == 0 || state.StartedBy == "" {
			// Server running but we don't have record of starting it = external
			slog.Info("Detected external Ollama server, will not auto-stop")
			if err := m.markExternal(); err != nil {
				slog.Warn("Failed to mark Ollama as external", "error", err)
			}
		}
		return nil
	}

	// Server not running - start it
	slog.Info("Starting Ollama server")
	if err := m.startServer(ctx); err != nil {
		return err
	}

	return nil
}

// RegisterUsage marks this instance as using Ollama with the given model.
func (m *Manager) RegisterUsage(model string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.active = true
	m.modelInUse = model

	return m.withLock(func(state *State) error {
		state.Instances[m.instanceID] = Instance{
			ID:           m.instanceID,
			PID:          os.Getpid(),
			Model:        model,
			LastActivity: time.Now(),
		}
		return nil
	})
}

// UnregisterUsage removes this instance from Ollama users.
// If no other instances are using Ollama and we started it, shuts it down.
func (m *Manager) UnregisterUsage() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.active {
		return nil
	}

	m.active = false
	m.modelInUse = ""

	return m.withLock(func(state *State) error {
		delete(state.Instances, m.instanceID)

		// Clean up stale instances
		now := time.Now()
		for id, inst := range state.Instances {
			if now.Sub(inst.LastActivity) > staleTimeout {
				delete(state.Instances, id)
			}
		}

		// Check if we should stop the server
		if len(state.Instances) == 0 && !state.IsExternal && state.StartedBy != "" {
			slog.Info("No instances using Ollama, stopping server")
			return m.stopServer(state)
		}

		return nil
	})
}

// Shutdown cleans up this instance's Ollama usage.
func (m *Manager) Shutdown() error {
	return m.UnregisterUsage()
}

// IsActive returns whether this instance is currently using Ollama.
func (m *Manager) IsActive() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.active
}

// startServer starts the Ollama server and records it in state.
func (m *Manager) startServer(ctx context.Context) error {
	cmd := exec.Command("ollama", "serve")
	cmd.Stdout = nil
	cmd.Stderr = nil
	// Detach from process group so it survives if Crush exits
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err := cmd.Start(); err != nil {
		return &OllamaError{
			Err:     ErrStartFailed,
			Title:   "Failed to Start Ollama",
			Message: fmt.Sprintf("Could not start Ollama server: %v", err),
			Action:  "Try starting manually with: ollama serve",
		}
	}

	m.serverCmd = cmd

	// Record that we started the server
	if err := m.withLock(func(state *State) error {
		state.ServerPID = cmd.Process.Pid
		state.StartedBy = m.instanceID
		state.IsExternal = false
		return nil
	}); err != nil {
		slog.Warn("Failed to record Ollama server state", "error", err)
	}

	// Wait for server to become available
	if err := m.waitForServer(ctx); err != nil {
		// Server failed to start, clean up
		_ = cmd.Process.Kill()
		return err
	}

	slog.Info("Ollama server started", "pid", cmd.Process.Pid)
	return nil
}

// waitForServer waits for the Ollama server to respond to health checks.
func (m *Manager) waitForServer(ctx context.Context) error {
	deadline := time.Now().Add(serverStartWait)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if IsServerRunning() {
				return nil
			}
			if time.Now().After(deadline) {
				return &OllamaError{
					Err:     ErrStartFailed,
					Title:   "Ollama Server Timeout",
					Message: "Ollama server did not start within timeout.",
					Action:  "Check Ollama logs or start manually with: ollama serve",
				}
			}
		}
	}
}

// stopServer stops the Ollama server.
func (m *Manager) stopServer(state *State) error {
	if state.ServerPID == 0 {
		return nil
	}

	proc, err := os.FindProcess(state.ServerPID)
	if err != nil {
		slog.Debug("Ollama process not found", "pid", state.ServerPID)
		state.ServerPID = 0
		state.StartedBy = ""
		return nil
	}

	// Send SIGTERM for graceful shutdown
	if err := proc.Signal(syscall.SIGTERM); err != nil {
		slog.Debug("Failed to send SIGTERM to Ollama", "error", err)
	}

	// Wait briefly for graceful shutdown
	time.Sleep(500 * time.Millisecond)

	// Clear state
	state.ServerPID = 0
	state.StartedBy = ""

	slog.Info("Ollama server stopped")
	return nil
}

// markExternal marks the current Ollama server as externally managed.
func (m *Manager) markExternal() error {
	return m.withLock(func(state *State) error {
		state.IsExternal = true
		return nil
	})
}

// withLock executes a function with file locking.
func (m *Manager) withLock(fn func(*State) error) error {
	lock, err := acquireLock(m.lockFilePath)
	if err != nil {
		return fmt.Errorf("failed to acquire Ollama lock: %w", err)
	}
	defer releaseLock(lock)

	state, err := m.readState()
	if err != nil {
		return err
	}

	if err := fn(state); err != nil {
		return err
	}

	state.Updated = time.Now()
	return m.writeState(state)
}

// readState reads the Ollama state from file.
func (m *Manager) readState() (*State, error) {
	data, err := os.ReadFile(m.stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{
				Instances: make(map[string]Instance),
			}, nil
		}
		return nil, fmt.Errorf("failed to read Ollama state: %w", err)
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return &State{
			Instances: make(map[string]Instance),
		}, nil
	}

	if state.Instances == nil {
		state.Instances = make(map[string]Instance)
	}

	return &state, nil
}

// writeState writes the Ollama state to file.
func (m *Manager) writeState(state *State) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Ollama state: %w", err)
	}

	if err := os.WriteFile(m.stateFile, data, 0o600); err != nil {
		return fmt.Errorf("failed to write Ollama state: %w", err)
	}

	return nil
}

// acquireLock creates a lock file for exclusive access.
func acquireLock(lockPath string) (*os.File, error) {
	for i := 0; i < 50; i++ {
		f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err == nil {
			return f, nil
		}
		if !os.IsExist(err) {
			return nil, err
		}

		// Check if lock is stale
		info, statErr := os.Stat(lockPath)
		if statErr == nil && time.Since(info.ModTime()) > 30*time.Second {
			os.Remove(lockPath)
		}

		time.Sleep(100 * time.Millisecond)
	}
	return nil, fmt.Errorf("timeout acquiring Ollama lock")
}

// releaseLock removes the lock file.
func releaseLock(f *os.File) {
	if f != nil {
		name := f.Name()
		f.Close()
		os.Remove(name)
	}
}

// IsInstalled returns whether Ollama is installed on this system.
func IsInstalled() bool {
	_, err := exec.LookPath("ollama")
	return err == nil
}

// IsServerRunning checks if Ollama server is responding.
func IsServerRunning() bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(healthCheckURL)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// GetInstalledModels returns a list of models available in Ollama.
func GetInstalledModels() ([]string, error) {
	cmd := exec.Command("ollama", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var models []string
	lines := string(output)
	for _, line := range splitLines(lines) {
		if line == "" || line == "NAME" {
			continue
		}
		// First column is model name
		parts := splitFields(line)
		if len(parts) > 0 {
			models = append(models, parts[0])
		}
	}
	return models, nil
}

// splitLines splits a string into lines.
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// splitFields splits a line into whitespace-separated fields.
func splitFields(s string) []string {
	var fields []string
	start := -1
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' || s[i] == '\t' {
			if start >= 0 {
				fields = append(fields, s[start:i])
				start = -1
			}
		} else {
			if start < 0 {
				start = i
			}
		}
	}
	if start >= 0 {
		fields = append(fields, s[start:])
	}
	return fields
}

// PullResult contains the result of a model pull operation.
type PullResult struct {
	Success bool
	Message string
}

// PullProgress contains progress information during model download.
type PullProgress struct {
	Status    string  // Current status message
	Digest    string  // Current layer digest being downloaded
	Total     int64   // Total bytes for current layer
	Completed int64   // Bytes completed for current layer
	Percent   float64 // Overall progress percentage (0-100)
	Done      bool    // True when pull is complete
	Error     error   // Error if pull failed
}

// PullModelWithProgress downloads a model and sends progress updates to the channel.
func PullModelWithProgress(ctx context.Context, model string, progress chan<- PullProgress) {
	defer close(progress)
	slog.Info("Pulling Ollama model with progress", "model", model)

	reqBody := fmt.Sprintf(`{"name":%q}`, model)
	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:11434/api/pull", strings.NewReader(reqBody))
	if err != nil {
		progress <- PullProgress{Error: err, Done: true}
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		progress <- PullProgress{Error: fmt.Errorf("failed to connect to Ollama: %w", err), Done: true}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		progress <- PullProgress{Error: fmt.Errorf("pull failed with status %d", resp.StatusCode), Done: true}
		return
	}

	decoder := json.NewDecoder(resp.Body)
	for {
		var line struct {
			Status    string `json:"status"`
			Digest    string `json:"digest"`
			Total     int64  `json:"total"`
			Completed int64  `json:"completed"`
			Error     string `json:"error"`
		}

		if err := decoder.Decode(&line); err != nil {
			if err.Error() == "EOF" {
				progress <- PullProgress{Status: "Complete", Percent: 100, Done: true}
				return
			}
			progress <- PullProgress{Error: err, Done: true}
			return
		}

		if line.Error != "" {
			progress <- PullProgress{Error: fmt.Errorf("%s", line.Error), Done: true}
			return
		}

		p := PullProgress{
			Status:    line.Status,
			Digest:    line.Digest,
			Total:     line.Total,
			Completed: line.Completed,
		}

		if line.Total > 0 {
			p.Percent = float64(line.Completed) / float64(line.Total) * 100
		}

		if line.Status == "success" {
			p.Done = true
			p.Percent = 100
		}

		progress <- p

		if p.Done {
			return
		}
	}
}

// PullModel downloads a model from Ollama's registry (blocking, no progress).
func PullModel(ctx context.Context, model string) PullResult {
	slog.Info("Pulling Ollama model", "model", model)

	progress := make(chan PullProgress)
	go PullModelWithProgress(ctx, model, progress)

	var lastProgress PullProgress
	for p := range progress {
		lastProgress = p
	}

	if lastProgress.Error != nil {
		slog.Error("Failed to pull model", "model", model, "error", lastProgress.Error)
		return PullResult{
			Success: false,
			Message: fmt.Sprintf("Failed to pull model: %s", lastProgress.Error),
		}
	}

	slog.Info("Successfully pulled model", "model", model)
	return PullResult{
		Success: true,
		Message: fmt.Sprintf("Successfully pulled model %s", model),
	}
}

// IsModelAvailable checks if a model is available locally.
func IsModelAvailable(model string) bool {
	models, err := GetInstalledModels()
	if err != nil {
		return false
	}
	for _, m := range models {
		// Check exact match and match without tag
		if m == model {
			return true
		}
		// Also check if model matches the base name (e.g., "llama3.2" matches "llama3.2:latest")
		if idx := indexOf(m, ':'); idx >= 0 {
			if m[:idx] == model {
				return true
			}
		}
	}
	return false
}

// indexOf returns the index of the first occurrence of c in s, or -1 if not found.
func indexOf(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}
