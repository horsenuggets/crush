// Package voice provides voice input using OpenAI Whisper API.
package voice

import (
	"context"
	"fmt"
	"os"
)

// VoiceInput manages voice recording and transcription.
type VoiceInput struct {
	recorder *Recorder
	whisper  *WhisperClient
}

// New creates a new VoiceInput instance.
// If openaiKey is empty, it will try to read from OPENAI_API_KEY environment variable.
func New(openaiKey string) (*VoiceInput, error) {
	whisper, err := NewWhisperClient(openaiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create whisper client: %w", err)
	}

	return &VoiceInput{
		recorder: NewRecorder(),
		whisper:  whisper,
	}, nil
}

// StartRecording begins capturing audio from the microphone.
func (v *VoiceInput) StartRecording(ctx context.Context) error {
	return v.recorder.Start(ctx)
}

// StopRecording stops recording and transcribes the audio.
// Returns the transcribed text.
// dynamicContext is optional session-specific terms to improve transcription accuracy.
func (v *VoiceInput) StopRecording(ctx context.Context, dynamicContext string) (string, error) {
	audioPath, err := v.recorder.Stop()
	if err != nil {
		return "", fmt.Errorf("failed to stop recording: %w", err)
	}
	defer os.Remove(audioPath)

	text, err := v.whisper.TranscribeFile(ctx, audioPath, dynamicContext)
	if err != nil {
		return "", fmt.Errorf("failed to transcribe: %w", err)
	}

	return text, nil
}

// CancelRecording stops recording without transcribing.
func (v *VoiceInput) CancelRecording() {
	v.recorder.Cancel()
}

// IsRecording returns true if currently recording.
func (v *VoiceInput) IsRecording() bool {
	return v.recorder.IsRecording()
}

// TranscribeFile transcribes an existing audio file.
func (v *VoiceInput) TranscribeFile(ctx context.Context, audioPath string, dynamicContext string) (string, error) {
	return v.whisper.TranscribeFile(ctx, audioPath, dynamicContext)
}

// Cleanup releases any resources.
func (v *VoiceInput) Cleanup() {
	v.recorder.Cleanup()
}
