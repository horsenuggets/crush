package rename

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/charmbracelet/crush/internal/tui/components/core"
	"github.com/charmbracelet/crush/internal/tui/components/dialogs"
	"github.com/charmbracelet/crush/internal/tui/styles"
	"github.com/charmbracelet/crush/internal/tui/util"
)

const (
	RenameDialogID dialogs.DialogID = "rename"

	defaultWidth int = 50
)

// RenameDialog provides an interface for the rename dialog.
type RenameDialog interface {
	dialogs.DialogModel
}

// SessionRenamedMsg is sent when the session has been renamed.
type SessionRenamedMsg struct {
	SessionID string
	NewTitle  string
}

// RenameDialogKeyMap defines key bindings for the rename dialog.
type RenameDialogKeyMap struct {
	Submit key.Binding
	Close  key.Binding
}

func DefaultRenameDialogKeyMap() RenameDialogKeyMap {
	return RenameDialogKeyMap{
		Submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "save"),
		),
		Close: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc", "cancel"),
		),
	}
}

func (k RenameDialogKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Submit, k.Close}
}

func (k RenameDialogKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Submit, k.Close},
	}
}

type renameDialogCmp struct {
	width     int
	wWidth    int
	wHeight   int
	sessionID string
	input     textinput.Model
	keyMap    RenameDialogKeyMap
	help      help.Model
}

// NewRenameDialog creates a new dialog for renaming a session.
func NewRenameDialog(sessionID, currentTitle string) RenameDialog {
	t := styles.CurrentTheme()

	input := textinput.New()
	input.Placeholder = "Enter new title..."
	input.SetValue(currentTitle)
	input.Focus()
	input.SetWidth(defaultWidth - 6)
	input.SetStyles(t.S().TextInput)

	help := help.New()
	help.Styles = t.S().Help

	return &renameDialogCmp{
		width:     defaultWidth,
		sessionID: sessionID,
		input:     input,
		keyMap:    DefaultRenameDialogKeyMap(),
		help:      help,
	}
}

func (r *renameDialogCmp) Init() tea.Cmd {
	return r.input.Focus()
}

func (r *renameDialogCmp) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.wWidth = msg.Width
		r.wHeight = msg.Height
		return r, nil
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, r.keyMap.Submit):
			newTitle := r.input.Value()
			if newTitle == "" {
				return r, nil // Don't save empty titles
			}
			return r, tea.Sequence(
				util.CmdHandler(dialogs.CloseDialogMsg{}),
				func() tea.Msg {
					return SessionRenamedMsg{
						SessionID: r.sessionID,
						NewTitle:  newTitle,
					}
				},
			)
		case key.Matches(msg, r.keyMap.Close):
			return r, util.CmdHandler(dialogs.CloseDialogMsg{})
		}
	}

	// Forward other messages to the text input
	var cmd tea.Cmd
	r.input, cmd = r.input.Update(msg)
	return r, cmd
}

func (r *renameDialogCmp) View() string {
	t := styles.CurrentTheme()
	if t.IsAnimated() {
		r.help.Styles = t.HelpStyles
	} else {
		r.help.Styles = t.S().Help
	}

	header := t.S().Base.Padding(0, 1, 1, 1).Render(core.Title("Rename Session", r.width-4))
	inputView := t.S().Base.Padding(0, 1).Render(r.input.View())
	helpView := t.S().Base.Width(r.width-2).PaddingLeft(1).AlignHorizontal(lipgloss.Left).Render(r.help.View(r.keyMap))

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		inputView,
		"",
		helpView,
	)
	return r.style().Render(content)
}

func (r *renameDialogCmp) Cursor() *tea.Cursor {
	cursor := r.input.Cursor()
	if cursor != nil {
		row, col := r.Position()
		cursor.Y += row + 2
		cursor.X += col + 2
	}
	return cursor
}

func (r *renameDialogCmp) style() lipgloss.Style {
	t := styles.CurrentTheme()
	return t.S().Base.
		Width(r.width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderFocus)
}

func (r *renameDialogCmp) Position() (int, int) {
	row := r.wHeight/4 - 2
	col := r.wWidth/2 - r.width/2
	return row, col
}

func (r *renameDialogCmp) ID() dialogs.DialogID {
	return RenameDialogID
}
