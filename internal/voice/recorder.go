package voice

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

// RecorderState represents the current state of the recorder.
type RecorderState int

const (
	RecorderStateIdle RecorderState = iota
	RecorderStateRecording
	RecorderStateStopping
)

// Recorder captures audio from the microphone.
type Recorder struct {
	mu       sync.Mutex
	state    RecorderState
	cmd      *exec.Cmd
	cancel   context.CancelFunc
	tempFile string
}

// NewRecorder creates a new audio recorder.
func NewRecorder() *Recorder {
	return &Recorder{
		state: RecorderStateIdle,
	}
}

// State returns the current recorder state.
func (r *Recorder) State() RecorderState {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.state
}

// IsRecording returns true if currently recording.
func (r *Recorder) IsRecording() bool {
	return r.State() == RecorderStateRecording
}

// Start begins recording audio.
func (r *Recorder) Start(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state != RecorderStateIdle {
		return fmt.Errorf("recorder is not idle")
	}

	// Create temp file for recording
	tempDir := os.TempDir()
	r.tempFile = filepath.Join(tempDir, fmt.Sprintf("crush_voice_%d.wav", time.Now().UnixNano()))

	// Create cancellable context
	recordCtx, cancel := context.WithCancel(ctx)
	r.cancel = cancel

	// Try sox first, then ffmpeg
	if err := r.startWithSox(recordCtx); err != nil {
		if err := r.startWithFFmpeg(recordCtx); err != nil {
			return fmt.Errorf("no audio recording tool available (tried sox and ffmpeg): %w", err)
		}
	}

	r.state = RecorderStateRecording
	return nil
}

// startWithSox attempts to start recording using sox.
func (r *Recorder) startWithSox(ctx context.Context) error {
	// Check if sox is available
	if _, err := exec.LookPath("sox"); err != nil {
		return err
	}

	// sox -d -r 16000 -c 1 -b 16 output.wav
	r.cmd = exec.CommandContext(ctx, "sox", "-d", "-r", "16000", "-c", "1", "-b", "16", r.tempFile)
	r.cmd.Stdout = nil
	r.cmd.Stderr = nil

	if err := r.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start sox: %w", err)
	}

	return nil
}

// startWithFFmpeg attempts to start recording using ffmpeg.
func (r *Recorder) startWithFFmpeg(ctx context.Context) error {
	// Check if ffmpeg is available
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return err
	}

	// Determine input device based on OS
	var inputArgs []string
	switch {
	case fileExists("/dev/snd"):
		// Linux ALSA
		inputArgs = []string{"-f", "alsa", "-i", "default"}
	case fileExists("/dev/audio"):
		// BSD/Solaris
		inputArgs = []string{"-f", "oss", "-i", "/dev/audio"}
	default:
		// macOS - use "default" to pick system default input, not hardcoded index
		// Index :0 may pick the wrong device (e.g., iPhone instead of MacBook mic)
		inputArgs = []string{"-f", "avfoundation", "-i", ":default"}
	}

	args := append(inputArgs,
		"-ar", "16000",
		"-ac", "1",
		"-acodec", "pcm_s16le",
		"-f", "wav",
		"-y",
		r.tempFile,
	)

	r.cmd = exec.CommandContext(ctx, "ffmpeg", args...)
	r.cmd.Stdout = nil
	r.cmd.Stderr = nil

	if err := r.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	return nil
}

// Stop stops recording and returns the path to the audio file.
func (r *Recorder) Stop() (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state != RecorderStateRecording {
		return "", fmt.Errorf("recorder is not recording")
	}

	r.state = RecorderStateStopping

	// Send SIGINT to ffmpeg to gracefully stop and finalize the file
	// Context cancellation kills the process abruptly which corrupts WAV headers
	if r.cmd != nil && r.cmd.Process != nil {
		_ = r.cmd.Process.Signal(syscall.SIGINT)
		// Give ffmpeg a moment to finalize the file
		time.Sleep(100 * time.Millisecond)
		_ = r.cmd.Wait()
	}

	// Cancel context after process is done
	if r.cancel != nil {
		r.cancel()
	}

	r.state = RecorderStateIdle

	// Check if file was created
	if _, err := os.Stat(r.tempFile); err != nil {
		return "", fmt.Errorf("recording file not created: %w", err)
	}

	return r.tempFile, nil
}

// Cancel stops recording without saving.
func (r *Recorder) Cancel() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state != RecorderStateRecording {
		return
	}

	r.state = RecorderStateStopping

	if r.cancel != nil {
		r.cancel()
	}

	if r.cmd != nil {
		_ = r.cmd.Wait()
	}

	// Clean up temp file
	if r.tempFile != "" {
		_ = os.Remove(r.tempFile)
	}

	r.state = RecorderStateIdle
}

// Cleanup removes any temporary files.
func (r *Recorder) Cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.tempFile != "" {
		_ = os.Remove(r.tempFile)
		r.tempFile = ""
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
