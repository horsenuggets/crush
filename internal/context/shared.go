// Package context provides shared context management for parallel Crush instances.
package context

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	sharedContextFile = "shared-context.json"
	lockFile          = ".shared-context.lock"
	staleTimeout      = 5 * time.Minute
)

// Instance represents a running Crush instance.
type Instance struct {
	ID           string    `json:"id"`
	SpawnDir     string    `json:"spawn_dir"`      // Directory where Crush was started
	WorkingDir   string    `json:"working_dir"`    // Current directory the model is operating in
	Task         string    `json:"task,omitempty"` // Current task description
	LastActivity time.Time `json:"last_activity"`
	PID          int       `json:"pid"`
}

// SharedContext holds the shared state between Crush instances.
type SharedContext struct {
	Instances map[string]Instance `json:"instances"`
	Updated   time.Time           `json:"updated"`
}

// Manager handles shared context operations.
type Manager struct {
	mu           sync.Mutex
	instanceID   string
	spawnDir     string
	contextDir   string
	contextFile  string
	lockFilePath string
}

// NewManager creates a new context manager.
func NewManager(spawnDir string) (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	contextDir := filepath.Join(homeDir, ".crush")
	if err := os.MkdirAll(contextDir, 0o700); err != nil {
		return nil, fmt.Errorf("failed to create context directory: %w", err)
	}

	return &Manager{
		instanceID:   uuid.New().String()[:8],
		spawnDir:     spawnDir,
		contextDir:   contextDir,
		contextFile:  filepath.Join(contextDir, sharedContextFile),
		lockFilePath: filepath.Join(contextDir, lockFile),
	}, nil
}

// InstanceID returns the unique ID for this instance.
func (m *Manager) InstanceID() string {
	return m.instanceID
}

// Register registers this instance in the shared context.
func (m *Manager) Register() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.withLock(func(ctx *SharedContext) error {
		ctx.Instances[m.instanceID] = Instance{
			ID:           m.instanceID,
			SpawnDir:     m.spawnDir,
			WorkingDir:   m.spawnDir,
			LastActivity: time.Now(),
			PID:          os.Getpid(),
		}
		return nil
	})
}

// Unregister removes this instance from the shared context.
func (m *Manager) Unregister() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.withLock(func(ctx *SharedContext) error {
		delete(ctx.Instances, m.instanceID)
		return nil
	})
}

// UpdateWorkingDir updates the working directory for this instance.
func (m *Manager) UpdateWorkingDir(workingDir string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.withLock(func(ctx *SharedContext) error {
		if instance, ok := ctx.Instances[m.instanceID]; ok {
			instance.WorkingDir = workingDir
			instance.LastActivity = time.Now()
			ctx.Instances[m.instanceID] = instance
		}
		return nil
	})
}

// UpdateTask updates the current task for this instance.
func (m *Manager) UpdateTask(task string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.withLock(func(ctx *SharedContext) error {
		if instance, ok := ctx.Instances[m.instanceID]; ok {
			instance.Task = task
			instance.LastActivity = time.Now()
			ctx.Instances[m.instanceID] = instance
		}
		return nil
	})
}

// Heartbeat updates the last activity time for this instance.
func (m *Manager) Heartbeat() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.withLock(func(ctx *SharedContext) error {
		if instance, ok := ctx.Instances[m.instanceID]; ok {
			instance.LastActivity = time.Now()
			ctx.Instances[m.instanceID] = instance
		}
		return nil
	})
}

// ListInstances returns all active instances.
func (m *Manager) ListInstances() ([]Instance, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var instances []Instance
	err := m.withLock(func(ctx *SharedContext) error {
		// Clean up stale instances and dead processes
		now := time.Now()
		for id, inst := range ctx.Instances {
			// Remove if stale (no activity for staleTimeout)
			if now.Sub(inst.LastActivity) > staleTimeout {
				delete(ctx.Instances, id)
				continue
			}
			// Remove if process is no longer running
			if inst.PID > 0 && !isProcessRunning(inst.PID) {
				delete(ctx.Instances, id)
				continue
			}
			instances = append(instances, inst)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Sort by ID for consistent ordering
	sortInstancesByID(instances)

	return instances, nil
}

// OtherInstances returns all instances except this one.
func (m *Manager) OtherInstances() ([]Instance, error) {
	instances, err := m.ListInstances()
	if err != nil {
		return nil, err
	}

	var others []Instance
	for _, inst := range instances {
		if inst.ID != m.instanceID {
			others = append(others, inst)
		}
	}
	return others, nil
}

// GetInstance returns information about a specific instance.
func (m *Manager) GetInstance(id string) (*Instance, error) {
	instances, err := m.ListInstances()
	if err != nil {
		return nil, err
	}

	for _, inst := range instances {
		if inst.ID == id {
			return &inst, nil
		}
	}
	return nil, fmt.Errorf("instance %s not found", id)
}

// withLock executes a function with file locking.
func (m *Manager) withLock(fn func(*SharedContext) error) error {
	// Acquire lock
	lock, err := acquireLock(m.lockFilePath)
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	defer releaseLock(lock)

	// Read existing context
	ctx, err := m.readContext()
	if err != nil {
		return err
	}

	// Execute function
	if err := fn(ctx); err != nil {
		return err
	}

	// Write updated context
	ctx.Updated = time.Now()
	return m.writeContext(ctx)
}

// readContext reads the shared context from file.
func (m *Manager) readContext() (*SharedContext, error) {
	data, err := os.ReadFile(m.contextFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &SharedContext{
				Instances: make(map[string]Instance),
			}, nil
		}
		return nil, fmt.Errorf("failed to read context file: %w", err)
	}

	var ctx SharedContext
	if err := json.Unmarshal(data, &ctx); err != nil {
		// If JSON is corrupted, start fresh
		return &SharedContext{
			Instances: make(map[string]Instance),
		}, nil
	}

	if ctx.Instances == nil {
		ctx.Instances = make(map[string]Instance)
	}

	return &ctx, nil
}

// writeContext writes the shared context to file.
func (m *Manager) writeContext(ctx *SharedContext) error {
	data, err := json.MarshalIndent(ctx, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal context: %w", err)
	}

	if err := os.WriteFile(m.contextFile, data, 0o600); err != nil {
		return fmt.Errorf("failed to write context file: %w", err)
	}

	return nil
}

// acquireLock creates a lock file for exclusive access.
func acquireLock(lockPath string) (*os.File, error) {
	for i := 0; i < 50; i++ { // Try for up to 5 seconds
		f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err == nil {
			return f, nil
		}
		if !os.IsExist(err) {
			return nil, err
		}

		// Check if lock is stale (older than 30 seconds)
		info, statErr := os.Stat(lockPath)
		if statErr == nil && time.Since(info.ModTime()) > 30*time.Second {
			os.Remove(lockPath)
		}

		time.Sleep(100 * time.Millisecond)
	}
	return nil, fmt.Errorf("timeout acquiring lock")
}

// releaseLock removes the lock file.
func releaseLock(f *os.File) {
	if f != nil {
		name := f.Name()
		f.Close()
		os.Remove(name)
	}
}

// isProcessRunning checks if a process with the given PID is still running.
func isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// On Unix, FindProcess always succeeds. Send signal 0 to check if process exists.
	err = process.Signal(os.Signal(nil))
	return err == nil
}

// sortInstancesByID sorts instances by their ID for consistent ordering.
func sortInstancesByID(instances []Instance) {
	for i := 0; i < len(instances)-1; i++ {
		for j := i + 1; j < len(instances); j++ {
			if instances[i].ID > instances[j].ID {
				instances[i], instances[j] = instances[j], instances[i]
			}
		}
	}
}
