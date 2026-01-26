package ollama

import (
	"context"
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/ollama"
	"github.com/charmbracelet/crush/internal/tui/components/dialogs"
	"github.com/charmbracelet/crush/internal/tui/styles"
	"github.com/charmbracelet/crush/internal/tui/util"
)

const (
	PullDialogID dialogs.DialogID = "ollama-pull"
)

// PullCompleteMsg is sent when model pull completes.
type PullCompleteMsg struct {
	Success bool
	Message string
	Model   string
}

// PullProgressMsg is sent during model download with progress updates.
type PullProgressMsg struct {
	Status  string
	Percent float64
	Done    bool
	Error   error
}

// PullDialog represents a dialog for pulling Ollama models.
type PullDialog interface {
	dialogs.DialogModel
}

type pullDialogCmp struct {
	wWidth       int
	wHeight      int
	model        string
	pulling      bool
	progress     float64
	status       string
	result       *ollama.PullResult
	selectedNo   bool
	keymap       pullKeyMap
	spinner      spinner.Model
	progressChan <-chan ollama.PullProgress
}

type pullKeyMap struct {
	LeftRight  key.Binding
	Tab        key.Binding
	EnterSpace key.Binding
	Yes        key.Binding
	No         key.Binding
	Close      key.Binding
}

func defaultPullKeymap() pullKeyMap {
	return pullKeyMap{
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

// NewPullDialog creates a new dialog for pulling an Ollama model.
func NewPullDialog(model string) PullDialog {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return &pullDialogCmp{
		model:      model,
		selectedNo: false,
		keymap:     defaultPullKeymap(),
		spinner:    s,
	}
}

func (d *pullDialogCmp) Init() tea.Cmd {
	return d.spinner.Tick
}

// Update handles keyboard input for the pull dialog.
func (d *pullDialogCmp) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Update spinner when pulling
	if d.pulling {
		d.spinner, cmd = d.spinner.Update(msg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.wWidth = msg.Width
		d.wHeight = msg.Height
		return d, cmd
	case PullProgressMsg:
		d.progress = msg.Percent
		d.status = msg.Status
		if msg.Done {
			d.pulling = false
			d.progressChan = nil
			if msg.Error != nil {
				d.result = &ollama.PullResult{
					Success: false,
					Message: msg.Error.Error(),
				}
				return d, nil
			}
			// Success - close dialog and signal to retry
			return d, tea.Batch(
				util.CmdHandler(dialogs.CloseDialogMsg{}),
				util.CmdHandler(PullCompleteMsg{
					Success: true,
					Message: "Download complete",
					Model:   d.model,
				}),
			)
		}
		// Continue listening for more progress
		return d, tea.Batch(cmd, d.waitForProgress())
	case PullCompleteMsg:
		d.pulling = false
		if msg.Success {
			// Close dialog and signal to retry
			return d, tea.Batch(
				util.CmdHandler(dialogs.CloseDialogMsg{}),
				util.CmdHandler(msg),
			)
		}
		// Show error - keep dialog open
		d.result = &ollama.PullResult{
			Success: false,
			Message: msg.Message,
		}
		return d, nil
	case tea.KeyPressMsg:
		if d.pulling {
			return d, cmd // Continue spinner while pulling
		}
		if d.result != nil && !d.result.Success {
			// After failed pull, any key closes
			return d, tea.Batch(
				util.CmdHandler(dialogs.CloseDialogMsg{}),
				util.CmdHandler(PullCompleteMsg{Success: false, Model: d.model}),
			)
		}
		switch {
		case key.Matches(msg, d.keymap.LeftRight, d.keymap.Tab):
			d.selectedNo = !d.selectedNo
			return d, nil
		case key.Matches(msg, d.keymap.EnterSpace):
			if !d.selectedNo {
				return d, d.startPull()
			}
			return d, tea.Batch(
				util.CmdHandler(dialogs.CloseDialogMsg{}),
				util.CmdHandler(PullCompleteMsg{Success: false, Model: d.model}),
			)
		case key.Matches(msg, d.keymap.Yes):
			return d, d.startPull()
		case key.Matches(msg, d.keymap.No, d.keymap.Close):
			return d, tea.Batch(
				util.CmdHandler(dialogs.CloseDialogMsg{}),
				util.CmdHandler(PullCompleteMsg{Success: false, Model: d.model}),
			)
		}
	}
	return d, cmd
}

func (d *pullDialogCmp) startPull() tea.Cmd {
	d.pulling = true
	d.progress = 0
	d.status = "Starting..."

	progress := make(chan ollama.PullProgress)
	d.progressChan = progress

	go ollama.PullModelWithProgress(context.Background(), d.model, progress)

	return tea.Batch(d.spinner.Tick, d.waitForProgress())
}

func (d *pullDialogCmp) waitForProgress() tea.Cmd {
	return func() tea.Msg {
		if d.progressChan == nil {
			return nil
		}
		p, ok := <-d.progressChan
		if !ok {
			return PullProgressMsg{Done: true}
		}
		return PullProgressMsg{
			Status:  p.Status,
			Percent: p.Percent,
			Done:    p.Done,
			Error:   p.Error,
		}
	}
}

// View renders the pull dialog.
func (d *pullDialogCmp) View() string {
	t := styles.CurrentTheme()
	baseStyle := t.S().Base

	var content string

	if d.pulling {
		content = d.renderPulling(t, baseStyle)
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

func (d *pullDialogCmp) renderPrompt(t *styles.Theme, baseStyle lipgloss.Style) string {
	var title, question, hint string
	if t.IsAnimated() {
		title = t.DialogTitle.Render("Model Not Found")
		question = t.DialogText.Render("Would you like to download " + d.model + "?")
		hint = t.DialogHint.Render("This may take a few minutes depending on model size")
	} else {
		title = t.S().Title.Render("Model Not Found")
		question = "Would you like to download " + d.model + "?"
		hint = t.S().Text.Faint(true).Render("This may take a few minutes depending on model size")
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

	buttons := baseStyle.Width(lipgloss.Width(hint)).Align(lipgloss.Center).Render(
		lipgloss.JoinHorizontal(lipgloss.Center, yesButton, "  ", noButton),
	)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		question,
		hint,
		"",
		buttons,
	)
}

func (d *pullDialogCmp) renderPulling(t *styles.Theme, baseStyle lipgloss.Style) string {
	var title string
	if t.IsAnimated() {
		title = t.DialogTitle.Render("Downloading " + d.model)
	} else {
		title = t.S().Title.Render("Downloading " + d.model)
	}

	// Progress bar
	const barWidth = 40
	filled := int(d.progress / 100 * barWidth)
	if filled > barWidth {
		filled = barWidth
	}
	empty := barWidth - filled

	barFilled := t.S().Base.Foreground(t.Primary).Render(strings.Repeat("█", filled))
	barEmpty := t.S().Base.Foreground(t.FgHalfMuted).Render(strings.Repeat("░", empty))
	progressBar := barFilled + barEmpty

	// Percentage (padded to 3 digits for consistent width)
	var percent string
	if t.IsAnimated() {
		percent = t.DialogText.Render(fmt.Sprintf(" %3.0f%%", d.progress))
	} else {
		percent = t.S().Text.Render(fmt.Sprintf(" %3.0f%%", d.progress))
	}

	// Status
	status := d.status
	if status == "" {
		status = "Downloading..."
	}
	var statusView string
	if t.IsAnimated() {
		statusView = d.spinner.View() + " " + t.DialogHint.Render(status)
	} else {
		statusView = d.spinner.View() + " " + t.S().Text.Faint(true).Render(status)
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		progressBar+percent,
		statusView,
	)
}

func (d *pullDialogCmp) renderError(t *styles.Theme, baseStyle lipgloss.Style) string {
	var title, message, hint string
	if t.IsAnimated() {
		title = t.DialogTitle.Render("Download Failed")
		message = t.DialogText.Render(d.result.Message)
		hint = t.DialogHint.Render("Press any key to close")
	} else {
		title = t.S().Error.Render("Download Failed")
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

func (d *pullDialogCmp) Position() (int, int) {
	var contentWidth int
	var dialogHeight int

	if d.pulling {
		// Downloading state
		title := "Downloading " + d.model
		contentWidth = lipgloss.Width(title)
		dialogHeight = 5
	} else if d.result != nil && !d.result.Success {
		// Error state
		contentWidth = lipgloss.Width(d.result.Message)
		if w := lipgloss.Width("Download Failed"); w > contentWidth {
			contentWidth = w
		}
		dialogHeight = 7
	} else {
		// Prompt state - use hint as it's the widest
		hint := "This may take a few minutes depending on model size"
		contentWidth = lipgloss.Width(hint)
		dialogHeight = 9
	}

	// Add padding (2 on each side) + border (1 on each side)
	dialogWidth := contentWidth + 6

	row := d.wHeight/2 - dialogHeight/2
	col := d.wWidth/2 - dialogWidth/2

	return row, col
}

func (d *pullDialogCmp) ID() dialogs.DialogID {
	return PullDialogID
}

// AutoCentered implements dialogs.AutoCenteredDialog for automatic centering.
func (d *pullDialogCmp) AutoCentered() {}
