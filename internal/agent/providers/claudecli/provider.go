// Package claudecli provides a fantasy.Provider that wraps the local Claude CLI.
// This allows using a Claude Max subscription instead of API access.
package claudecli

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"sync"

	"charm.land/fantasy"
	"charm.land/fantasy/object"
)

const (
	// Name is the provider name.
	Name = "claude-cli"
)

// ErrCLINotFound is returned when the claude executable is not found.
var ErrCLINotFound = errors.New("claude CLI not found in PATH")

// toolNameMap maps Claude CLI tool names to Crush tool names.
// Claude CLI uses PascalCase, Crush uses lowercase.
var toolNameMap = map[string]string{
	"Bash":  "bash",
	"Read":  "view",
	"Edit":  "edit",
	"Write": "write",
	"Glob":  "glob",
	"Grep":  "grep",
	"LS":    "ls",
	"Task":  "agent",
}

type options struct {
	executablePath  string
	workingDir      string
	model           string // default model to use
	skipPermissions bool
}

// Option configures the provider.
type Option = func(*options)

// WithExecutablePath sets the path to the claude executable.
func WithExecutablePath(path string) Option {
	return func(o *options) { o.executablePath = path }
}

// WithWorkingDir sets the working directory for claude commands.
func WithWorkingDir(dir string) Option {
	return func(o *options) { o.workingDir = dir }
}

// WithModel sets the default model to use.
func WithModel(model string) Option {
	return func(o *options) { o.model = model }
}

// WithSkipPermissions skips permission checks in the Claude CLI.
func WithSkipPermissions(skip bool) Option {
	return func(o *options) { o.skipPermissions = skip }
}

type provider struct {
	opts options
}

// New creates a new Claude CLI provider.
func New(opts ...Option) (fantasy.Provider, error) {
	o := options{
		executablePath:  "claude",
		skipPermissions: true, // Default to skipping since Crush has its own permission system
	}
	for _, opt := range opts {
		opt(&o)
	}

	// Verify the executable exists
	path, err := exec.LookPath(o.executablePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCLINotFound, o.executablePath)
	}
	o.executablePath = path

	return &provider{opts: o}, nil
}

func (p *provider) Name() string { return Name }

// LanguageModel implements fantasy.Provider.
func (p *provider) LanguageModel(_ context.Context, modelID string) (fantasy.LanguageModel, error) {
	if modelID == "" {
		modelID = p.opts.model
	}
	if modelID == "" {
		modelID = "opus" // Default to opus
	}
	return &languageModel{
		modelID:  modelID,
		provider: Name,
		opts:     p.opts,
	}, nil
}

type languageModel struct {
	provider string
	modelID  string
	opts     options
}

func (m *languageModel) Provider() string { return m.provider }
func (m *languageModel) Model() string    { return m.modelID }

// GenerateObject implements fantasy.LanguageModel.
func (m *languageModel) GenerateObject(ctx context.Context, call fantasy.ObjectCall) (*fantasy.ObjectResponse, error) {
	return object.GenerateWithTool(ctx, m, call)
}

// StreamObject implements fantasy.LanguageModel.
func (m *languageModel) StreamObject(ctx context.Context, call fantasy.ObjectCall) (fantasy.ObjectStreamResponse, error) {
	return object.StreamWithTool(ctx, m, call)
}

// Generate implements fantasy.LanguageModel.
func (m *languageModel) Generate(ctx context.Context, call fantasy.Call) (*fantasy.Response, error) {
	// For non-streaming, collect all stream parts into a response
	stream, err := m.Stream(ctx, call)
	if err != nil {
		return nil, err
	}

	var response fantasy.Response
	var content []fantasy.Content

	for part := range stream {
		switch part.Type {
		case fantasy.StreamPartTypeTextDelta:
			content = append(content, fantasy.TextContent{Text: part.Delta})
		case fantasy.StreamPartTypeToolCall:
			content = append(content, fantasy.ToolCallContent{
				ToolCallID: part.ID,
				ToolName:   part.ToolCallName,
				Input:      part.ToolCallInput,
			})
		case fantasy.StreamPartTypeError:
			return nil, part.Error
		case fantasy.StreamPartTypeFinish:
			response.FinishReason = part.FinishReason
			response.Usage = part.Usage
		}
	}

	response.Content = content
	slog.Debug("Claude CLI Generate response",
		"content_count", len(content),
		"finish_reason", response.FinishReason,
		"has_tool_calls", response.FinishReason == fantasy.FinishReasonToolCalls)
	return &response, nil
}

// Stream implements fantasy.LanguageModel.
func (m *languageModel) Stream(ctx context.Context, call fantasy.Call) (fantasy.StreamResponse, error) {
	// Build the command arguments
	args := m.buildArgs(call)

	// Build prompt and log it for debugging
	prompt := m.buildPrompt(call)
	slog.Debug("Starting claude CLI", "args", args, "working_dir", m.opts.workingDir, "prompt_length", len(prompt))
	slog.Debug("Claude CLI prompt", "prompt", prompt[:min(500, len(prompt))])

	cmd := exec.CommandContext(ctx, m.opts.executablePath, args...)
	if m.opts.workingDir != "" {
		cmd.Dir = m.opts.workingDir
	}

	// Set up environment
	cmd.Env = os.Environ()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start claude CLI: %w", err)
	}

	// Send the prompt via stdin (already built above for logging)
	go func() {
		defer stdin.Close()
		_, _ = io.WriteString(stdin, prompt)
	}()

	// Capture stderr for error reporting
	var stderrBuf bytes.Buffer
	var stderrMu sync.Mutex
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				stderrMu.Lock()
				stderrBuf.Write(buf[:n])
				stderrMu.Unlock()
			}
			if err != nil {
				break
			}
		}
	}()

	return func(yield func(fantasy.StreamPart) bool) {
		defer func() {
			_ = cmd.Wait()
		}()

		parser := newStreamParser()
		scanner := bufio.NewScanner(stdout)
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 4*1024*1024)

		var sawFinish bool

		for scanner.Scan() {
			line := scanner.Text()

			// Parse the JSON event
			parts := parser.parseLine(line)
			for _, part := range parts {
				if part.Type == fantasy.StreamPartTypeFinish {
					sawFinish = true
					slog.Debug("Claude CLI emitting finish", "finish_reason", part.FinishReason)
				}
				if part.Type == fantasy.StreamPartTypeToolCall {
					slog.Debug("Claude CLI emitting tool call", "tool_name", part.ToolCallName)
				}
				if !yield(part) {
					return
				}
			}
		}

		if err := scanner.Err(); err != nil &&
			!errors.Is(err, context.Canceled) &&
			!errors.Is(err, context.DeadlineExceeded) {
			yield(fantasy.StreamPart{Type: fantasy.StreamPartTypeError, Error: err})
			return
		}

		// Check for command error
		if err := cmd.Wait(); err != nil {
			stderrMu.Lock()
			errMsg := strings.TrimSpace(stderrBuf.String())
			stderrMu.Unlock()
			if errMsg == "" {
				errMsg = err.Error()
			}
			// Only report if we haven't already finished successfully
			if !sawFinish {
				yield(fantasy.StreamPart{
					Type:  fantasy.StreamPartTypeError,
					Error: fmt.Errorf("claude CLI error: %s", errMsg),
				})
			}
			return
		}

		if !sawFinish {
			finishReason := fantasy.FinishReasonStop
			if parser.hadToolCalls {
				finishReason = fantasy.FinishReasonToolCalls
			}
			finishPart := fantasy.StreamPart{
				Type:         fantasy.StreamPartTypeFinish,
				FinishReason: finishReason,
			}
			if parser.lastUsage != nil {
				finishPart.Usage = *parser.lastUsage
			}
			yield(finishPart)
			slog.Debug("Claude CLI stream ended, emitting finish", "finish_reason", finishReason, "had_tool_calls", parser.hadToolCalls)
		}
	}, nil
}

func (m *languageModel) buildArgs(call fantasy.Call) []string {
	args := []string{
		"--print",                        // Non-interactive mode
		"--output-format", "stream-json", // JSON streaming output
		"--verbose",                      // More detailed output
	}

	// Skip permission checks (Crush has its own permission system)
	if m.opts.skipPermissions {
		args = append(args, "--dangerously-skip-permissions")
	}

	// Add model
	args = append(args, "--model", m.modelID)

	// Extract system prompt from messages
	for _, msg := range call.Prompt {
		if msg.Role == fantasy.MessageRoleSystem {
			// Get text content from the message
			for _, part := range msg.Content {
				if textPart, ok := part.(fantasy.TextPart); ok {
					args = append(args, "--system-prompt", textPart.Text)
					break
				}
			}
			break
		}
	}

	return args
}

func (m *languageModel) buildPrompt(call fantasy.Call) string {
	var sb strings.Builder
	hasToolResults := false

	for _, msg := range call.Prompt {
		// Skip system messages (handled via --system-prompt flag)
		if msg.Role == fantasy.MessageRoleSystem {
			continue
		}

		for _, part := range msg.Content {
			switch p := part.(type) {
			case fantasy.TextPart:
				if msg.Role == fantasy.MessageRoleUser {
					sb.WriteString(p.Text)
					sb.WriteString("\n")
				} else if msg.Role == fantasy.MessageRoleAssistant {
					// Include assistant responses for context
					sb.WriteString("[Assistant]: ")
					sb.WriteString(p.Text)
					sb.WriteString("\n")
				}
			case fantasy.ToolCallPart:
				// Show what tool was called
				sb.WriteString(fmt.Sprintf("[Tool Call: %s]\n", p.ToolName))
			case fantasy.ToolResultPart:
				hasToolResults = true
				// Include tool results so Claude knows what happened
				sb.WriteString("[Tool Result]:\n")
				if p.Output != nil {
					if textResult, ok := p.Output.(fantasy.ToolResultOutputContentText); ok {
						sb.WriteString(textResult.Text)
						sb.WriteString("\n")
					}
				}
			}
		}
	}

	// If we have tool results, add instruction to respond
	if hasToolResults {
		sb.WriteString("\nBased on the tool results above, please provide your response to the user's original question.")
	}

	return strings.TrimSpace(sb.String())
}

// Claude CLI JSON event types
type claudeEvent struct {
	Type string `json:"type"`

	// For message_start, assistant events
	Message *claudeMessage `json:"message,omitempty"`

	// For content_block_start
	Index        int           `json:"index,omitempty"`
	ContentBlock *contentBlock `json:"content_block,omitempty"`

	// For content_block_delta
	Delta *contentDelta `json:"delta,omitempty"`

	// For message_delta
	Usage *usage `json:"usage,omitempty"`
}

type claudeMessage struct {
	ID         string         `json:"id,omitempty"`
	Role       string         `json:"role,omitempty"`
	Model      string         `json:"model,omitempty"`
	Content    []contentBlock `json:"content,omitempty"`
	StopReason string         `json:"stop_reason,omitempty"`
	Usage      *usage         `json:"usage,omitempty"`
}

type contentBlock struct {
	Type  string          `json:"type"`
	Text  string          `json:"text,omitempty"`
	ID    string          `json:"id,omitempty"`
	Name  string          `json:"name,omitempty"`
	Input json.RawMessage `json:"input,omitempty"`
}

type contentDelta struct {
	Type        string `json:"type"`
	Text        string `json:"text,omitempty"`
	PartialJSON string `json:"partial_json,omitempty"`
	Thinking    string `json:"thinking,omitempty"`
}

type usage struct {
	InputTokens  int64 `json:"input_tokens"`
	OutputTokens int64 `json:"output_tokens"`
}

// streamParser maintains state while parsing Claude CLI output
type streamParser struct {
	currentToolID    string
	currentToolName  string
	currentToolInput strings.Builder
	hadToolCalls     bool
	lastUsage        *fantasy.Usage
}

func newStreamParser() *streamParser {
	return &streamParser{}
}

func (p *streamParser) parseLine(line string) []fantasy.StreamPart {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	var event claudeEvent
	if err := json.Unmarshal([]byte(line), &event); err != nil {
		slog.Debug("Failed to parse claude CLI output", "line", line, "error", err)
		return nil
	}

	slog.Debug("Claude CLI event received", "type", event.Type)
	return p.processEvent(&event)
}

func (p *streamParser) processEvent(event *claudeEvent) []fantasy.StreamPart {
	var parts []fantasy.StreamPart

	switch event.Type {
	case "message_start", "assistant":
		if event.Message != nil {
			slog.Debug("Claude CLI assistant event", "content_count", len(event.Message.Content), "stop_reason", event.Message.StopReason)
			// Process any content in the message
			for _, block := range event.Message.Content {
				slog.Debug("Claude CLI content block", "type", block.Type, "text_len", len(block.Text))
				parts = append(parts, p.processContentBlock(&block)...)
			}
			// Store usage if present (will be used when stream ends)
			if event.Message.Usage != nil {
				p.lastUsage = &fantasy.Usage{
					InputTokens:  event.Message.Usage.InputTokens,
					OutputTokens: event.Message.Usage.OutputTokens,
				}
			}
		}

	case "result":
		// Claude CLI final result - just log it, don't emit finish here
		// Finish will be emitted when the stream ends
		slog.Debug("Claude CLI result event received")

	case "content_block_start":
		if event.ContentBlock != nil {
			parts = append(parts, p.processContentBlock(event.ContentBlock)...)
		}

	case "content_block_delta":
		if event.Delta != nil {
			switch event.Delta.Type {
			case "text_delta":
				if event.Delta.Text != "" {
					parts = append(parts, fantasy.StreamPart{
						Type:  fantasy.StreamPartTypeTextDelta,
						Delta: event.Delta.Text,
					})
				}
			case "input_json_delta":
				// Accumulate tool input
				p.currentToolInput.WriteString(event.Delta.PartialJSON)
			case "thinking_delta":
				if event.Delta.Thinking != "" {
					parts = append(parts, fantasy.StreamPart{
						Type:  fantasy.StreamPartTypeReasoningDelta,
						Delta: event.Delta.Thinking,
					})
				}
			}
		}

	case "content_block_stop":
		// Finalize tool use if we were accumulating one
		if p.currentToolName != "" {
			input := p.currentToolInput.String()
			parts = append(parts, fantasy.StreamPart{
				Type:          fantasy.StreamPartTypeToolCall,
				ID:            p.currentToolID,
				ToolCallName:  mapToolName(p.currentToolName),
				ToolCallInput: input,
			})
			p.hadToolCalls = true
			p.currentToolID = ""
			p.currentToolName = ""
			p.currentToolInput.Reset()
		}

	case "message_delta":
		// Could contain usage info
		if event.Usage != nil {
			parts = append(parts, fantasy.StreamPart{
				Type: fantasy.StreamPartTypeFinish,
				Usage: fantasy.Usage{
					InputTokens:  event.Usage.InputTokens,
					OutputTokens: event.Usage.OutputTokens,
				},
			})
		}

	case "message_stop":
		finishReason := fantasy.FinishReasonStop
		if p.hadToolCalls {
			finishReason = fantasy.FinishReasonToolCalls
		}
		parts = append(parts, fantasy.StreamPart{
			Type:         fantasy.StreamPartTypeFinish,
			FinishReason: finishReason,
		})
	}

	return parts
}

func (p *streamParser) processContentBlock(block *contentBlock) []fantasy.StreamPart {
	var parts []fantasy.StreamPart

	switch block.Type {
	case "text":
		if block.Text != "" {
			parts = append(parts, fantasy.StreamPart{
				Type:  fantasy.StreamPartTypeTextDelta,
				Delta: block.Text,
			})
		}

	case "tool_use":
		// Store tool info, we'll emit when we get all input
		p.currentToolID = block.ID
		p.currentToolName = block.Name
		p.currentToolInput.Reset()
		if len(block.Input) > 0 {
			p.currentToolInput.Write(block.Input)
		}
		// If input is already complete (no deltas expected), emit now
		if len(block.Input) > 0 {
			parts = append(parts, fantasy.StreamPart{
				Type:          fantasy.StreamPartTypeToolCall,
				ID:            block.ID,
				ToolCallName:  mapToolName(block.Name),
				ToolCallInput: string(block.Input),
			})
			p.hadToolCalls = true
			p.currentToolID = ""
			p.currentToolName = ""
			p.currentToolInput.Reset()
		}

	case "thinking":
		// Extended thinking - handled via deltas
	}

	return parts
}

// mapToolName converts Claude CLI tool names to Crush tool names.
func mapToolName(name string) string {
	if mapped, ok := toolNameMap[name]; ok {
		return mapped
	}
	return name
}
