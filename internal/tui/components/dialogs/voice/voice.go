// Package voice provides a dialog for collecting OpenAI API key for voice input.
package voice

import (
	"fmt"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/config"
	"github.com/charmbracelet/crush/internal/tui/components/core"
	"github.com/charmbracelet/crush/internal/tui/components/dialogs"
	"github.com/charmbracelet/crush/internal/tui/components/dialogs/models"
	"github.com/charmbracelet/crush/internal/tui/styles"
	"github.com/charmbracelet/crush/internal/tui/util"
	"github.com/charmbracelet/crush/internal/voice"
)

const (
	VoiceDialogID dialogs.DialogID = "voice-api-key"
	defaultWidth                   = 56
)

// VoiceAPIKeySetMsg is sent when the API key is successfully set and verified.
type VoiceAPIKeySetMsg struct{}

// VoiceDialog interface for the voice API key dialog.
type VoiceDialog interface {
	dialogs.DialogModel
}

type voiceDialogCmp struct {
	width   int
	wWidth  int
	wHeight int

	apiKeyInput  *models.APIKeyInput
	keyMap       KeyMap
	help         help.Model
	isAPIKeyValid bool
	apiKeyValue   string
}

// NewVoiceDialog creates a new voice API key dialog.
func NewVoiceDialog() VoiceDialog {
	t := styles.CurrentTheme()
	apiKeyInput := models.NewAPIKeyInput()
	apiKeyInput.SetProviderName("OpenAI")
	apiKeyInput.SetShowTitle(false)

	help := help.New()
	help.Styles = t.S().Help

	return &voiceDialogCmp{
		apiKeyInput: apiKeyInput,
		width:       defaultWidth,
		keyMap:      DefaultKeyMap(),
		help:        help,
	}
}

func (v *voiceDialogCmp) Init() tea.Cmd {
	return v.apiKeyInput.Init()
}

func (v *voiceDialogCmp) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.wWidth = msg.Width
		v.wHeight = msg.Height
		v.apiKeyInput.SetWidth(v.width - 2)
		v.help.SetWidth(v.width - 2)
		return v, nil
	case models.APIKeyStateChangeMsg:
		u, cmd := v.apiKeyInput.Update(msg)
		v.apiKeyInput = u.(*models.APIKeyInput)
		return v, cmd
	case spinner.TickMsg:
		u, cmd := v.apiKeyInput.Update(msg)
		v.apiKeyInput = u.(*models.APIKeyInput)
		return v, cmd
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, v.keyMap.Submit):
			if v.isAPIKeyValid {
				// Save and close
				return v, v.saveAPIKeyAndClose()
			}
			// Verify the API key
			v.apiKeyValue = v.apiKeyInput.Value()
			if v.apiKeyValue == "" {
				return v, util.ReportWarn("Please enter an API key")
			}
			return v, tea.Sequence(
				util.CmdHandler(models.APIKeyStateChangeMsg{
					State: models.APIKeyInputStateVerifying,
				}),
				func() tea.Msg {
					start := time.Now()
					// Test the API key by trying to create a voice input
					_, err := voice.New(v.apiKeyValue)
					// Wait at least 500ms so user sees spinner
					elapsed := time.Since(start)
					if elapsed < 500*time.Millisecond {
						time.Sleep(500*time.Millisecond - elapsed)
					}
					if err == nil {
						v.isAPIKeyValid = true
						return models.APIKeyStateChangeMsg{
							State: models.APIKeyInputStateVerified,
						}
					}
					return models.APIKeyStateChangeMsg{
						State: models.APIKeyInputStateError,
					}
				},
			)
		case key.Matches(msg, v.keyMap.Close):
			return v, util.CmdHandler(dialogs.CloseDialogMsg{})
		default:
			u, cmd := v.apiKeyInput.Update(msg)
			v.apiKeyInput = u.(*models.APIKeyInput)
			return v, cmd
		}
	case tea.PasteMsg:
		u, cmd := v.apiKeyInput.Update(msg)
		v.apiKeyInput = u.(*models.APIKeyInput)
		return v, cmd
	default:
		u, cmd := v.apiKeyInput.Update(msg)
		v.apiKeyInput = u.(*models.APIKeyInput)
		return v, cmd
	}
}

func (v *voiceDialogCmp) View() string {
	t := styles.CurrentTheme()
	if t.IsAnimated() {
		v.help.Styles = t.HelpStyles
	} else {
		v.help.Styles = t.S().Help
	}

	title := core.Title("Voice Input Setup", v.width-4)
	description := t.S().Muted.Render("Enter your OpenAI API key to enable voice transcription via Whisper.")

	apiKeyView := v.apiKeyInput.View()
	apiKeyView = t.S().Base.Width(v.width - 3).Height(lipgloss.Height(apiKeyView)).PaddingLeft(1).Render(apiKeyView)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		t.S().Base.Padding(0, 1, 0, 1).Render(title),
		t.S().Base.Padding(0, 1, 1, 1).Render(description),
		apiKeyView,
		"",
		t.S().Base.Width(v.width-2).PaddingLeft(1).AlignHorizontal(lipgloss.Left).Render(v.help.View(v.keyMap)),
	)

	return v.style().Render(content)
}

func (v *voiceDialogCmp) Cursor() *tea.Cursor {
	cursor := v.apiKeyInput.Cursor()
	if cursor != nil {
		row, col := v.Position()
		// Adjust for dialog position and padding
		cursor.Y += row + 4 // Border + title + description + spacing
		cursor.X += col + 2
	}
	return cursor
}

func (v *voiceDialogCmp) style() lipgloss.Style {
	t := styles.CurrentTheme()
	return t.S().Base.
		Width(v.width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderFocus)
}

func (v *voiceDialogCmp) Position() (int, int) {
	row := v.wHeight/4 - 2
	col := v.wWidth/2 - v.width/2
	return row, col
}

func (v *voiceDialogCmp) ID() dialogs.DialogID {
	return VoiceDialogID
}

func (v *voiceDialogCmp) saveAPIKeyAndClose() tea.Cmd {
	cfg := config.Get()

	// Ensure TUI and Voice options exist
	if cfg.Options == nil {
		cfg.Options = &config.Options{}
	}
	if cfg.Options.TUI == nil {
		cfg.Options.TUI = &config.TUIOptions{}
	}
	if cfg.Options.TUI.Voice == nil {
		cfg.Options.TUI.Voice = &config.VoiceOptions{}
	}

	// Save the API key
	cfg.Options.TUI.Voice.APIKey = v.apiKeyValue
	if err := cfg.SetConfigField("options.tui.voice.api_key", v.apiKeyValue); err != nil {
		return util.ReportError(fmt.Errorf("failed to save API key: %w", err))
	}

	return tea.Sequence(
		util.CmdHandler(dialogs.CloseDialogMsg{}),
		util.CmdHandler(VoiceAPIKeySetMsg{}),
	)
}
