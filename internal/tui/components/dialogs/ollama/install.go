package ollama

import (
	"context"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/ollama"
	"github.com/charmbracelet/crush/internal/tui/components/dialogs"
	"github.com/charmbracelet/crush/internal/tui/styles"
	"github.com/charmbracelet/crush/internal/tui/util"
)

const (
	InstallDialogID dialogs.DialogID = "ollama-install"
)

// InstallCompleteMsg is sent when Ollama installation completes.
type InstallCompleteMsg struct {
	Success bool
	Message string
}

// InstallDialog represents a confirmation dialog for installing Ollama.
type InstallDialog interface {
	dialogs.DialogModel
}

type installDialogCmp struct {
	wWidth     int
	wHeight    int
	selectedNo bool // true if "No" button is selected
	installing bool // true while installation is in progress
	result     *ollama.InstallResult
	keymap     keyMap
}

type keyMap struct {
	LeftRight  key.Binding
	Tab        key.Binding
	EnterSpace key.Binding
	Yes        key.Binding
	No         key.Binding
	Close      key.Binding
}

func defaultKeymap() keyMap {
	return keyMap{
		LeftRight: key.NewBinding(
			key.WithKeys("left", "right", "h", "l"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
		),
		EnterSpace: key.NewBinding(
			key.WithKeys("enter", " "),
		),
		Yes: key.NewBinding(
			key.WithKeys("y", "Y"),
		),
		No: key.NewBinding(
			key.WithKeys("n", "N"),
		),
		Close: key.NewBinding(
			key.WithKeys("esc", "q"),
		),
	}
}

// NewInstallDialog creates a new Ollama install confirmation dialog.
func NewInstallDialog() InstallDialog {
	return &installDialogCmp{
		selectedNo: false, // Default to "Yes" since user likely wants to install
		keymap:     defaultKeymap(),
	}
}

func (d *installDialogCmp) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input for the install dialog.
func (d *installDialogCmp) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.wWidth = msg.Width
		d.wHeight = msg.Height
	case InstallCompleteMsg:
		d.installing = false
		if msg.Success {
			// Close dialog and signal to retry
			return d, tea.Batch(
				util.CmdHandler(dialogs.CloseDialogMsg{}),
				util.CmdHandler(msg),
			)
		}
		// Show error - keep dialog open
		d.result = &ollama.InstallResult{
			Success: false,
			Message: msg.Message,
		}
		return d, nil
	case tea.KeyPressMsg:
		if d.installing {
			return d, nil // Ignore input while installing
		}
		if d.result != nil && !d.result.Success {
			// After failed install, any key closes
			return d, util.CmdHandler(dialogs.CloseDialogMsg{})
		}
		switch {
		case key.Matches(msg, d.keymap.LeftRight, d.keymap.Tab):
			d.selectedNo = !d.selectedNo
			return d, nil
		case key.Matches(msg, d.keymap.EnterSpace):
			if !d.selectedNo {
				return d, d.startInstall()
			}
			return d, util.CmdHandler(dialogs.CloseDialogMsg{})
		case key.Matches(msg, d.keymap.Yes):
			return d, d.startInstall()
		case key.Matches(msg, d.keymap.No, d.keymap.Close):
			return d, util.CmdHandler(dialogs.CloseDialogMsg{})
		}
	}
	return d, nil
}

func (d *installDialogCmp) startInstall() tea.Cmd {
	d.installing = true
	return func() tea.Msg {
		result := ollama.Install(context.Background())
		return InstallCompleteMsg{
			Success: result.Success,
			Message: result.Message,
		}
	}
}

// View renders the install dialog.
func (d *installDialogCmp) View() string {
	t := styles.CurrentTheme()
	baseStyle := t.S().Base

	var content string

	if d.installing {
		content = d.renderInstalling(t, baseStyle)
	} else if d.result != nil && !d.result.Success {
		content = d.renderError(t, baseStyle)
	} else {
		content = d.renderPrompt(t, baseStyle)
	}

	dialogStyle := baseStyle.
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderFocus)

	return dialogStyle.Render(content)
}

func (d *installDialogCmp) renderPrompt(t *styles.Theme, baseStyle lipgloss.Style) string {
	var title, question, method string
	if t.IsAnimated() {
		title = t.DialogTitle.Render("Ollama Not Installed")
		question = t.DialogText.Render("Would you like to install Ollama automatically?")
		method = t.DialogHint.Render(ollama.InstallMethod())
	} else {
		title = t.S().Title.Render("Ollama Not Installed")
		question = "Would you like to install Ollama automatically?"
		method = t.S().Text.Faint(true).Render(ollama.InstallMethod())
	}

	yesStyle := t.S().Text
	noStyle := yesStyle

	if d.selectedNo {
		noStyle = noStyle.Foreground(t.ButtonColor()).Background(t.Secondary)
		yesStyle = yesStyle.Background(t.BgSubtle)
	} else {
		yesStyle = yesStyle.Foreground(t.ButtonColor()).Background(t.Secondary)
		noStyle = noStyle.Background(t.BgSubtle)
	}

	const horizontalPadding = 3
	yesButton := yesStyle.PaddingLeft(horizontalPadding).Underline(true).Render("Y") +
		yesStyle.PaddingRight(horizontalPadding).Render("es")
	noButton := noStyle.PaddingLeft(horizontalPadding).Underline(true).Render("N") +
		noStyle.PaddingRight(horizontalPadding).Render("o")

	buttons := baseStyle.Width(lipgloss.Width(question)).Align(lipgloss.Center).Render(
		lipgloss.JoinHorizontal(lipgloss.Center, yesButton, "  ", noButton),
	)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		question,
		method,
		"",
		buttons,
	)
}

func (d *installDialogCmp) renderInstalling(t *styles.Theme, baseStyle lipgloss.Style) string {
	var title, spinner string
	if t.IsAnimated() {
		title = t.DialogTitle.Render("Installing Ollama")
		spinner = t.DialogHint.Render("Please wait...")
	} else {
		title = t.S().Title.Render("Installing Ollama")
		spinner = t.S().Text.Faint(true).Render("Please wait...")
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		spinner,
	)
}

func (d *installDialogCmp) renderError(t *styles.Theme, baseStyle lipgloss.Style) string {
	var title, message, hint string
	if t.IsAnimated() {
		title = t.DialogTitle.Render("Installation Failed")
		message = t.DialogText.Render(d.result.Message)
		hint = t.DialogHint.Render("Press any key to close")
	} else {
		title = t.S().Error.Render("Installation Failed")
		message = t.S().Text.Render(d.result.Message)
		hint = t.S().Text.Faint(true).Render("Press any key to close")
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		message,
		"",
		hint,
	)
}

func (d *installDialogCmp) Position() (int, int) {
	// Calculate based on the widest content
	question := "Would you like to install Ollama automatically?"
	dialogWidth := lipgloss.Width(question) + 8 // Add padding
	dialogHeight := 8

	row := d.wHeight/2 - dialogHeight/2
	col := d.wWidth/2 - dialogWidth/2

	return row, col
}

func (d *installDialogCmp) ID() dialogs.DialogID {
	return InstallDialogID
}

// AutoCentered implements dialogs.AutoCenteredDialog for automatic centering.
func (d *installDialogCmp) AutoCentered() {}
