package theme

import (
	"slices"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/charmbracelet/crush/internal/tui/components/core"
	"github.com/charmbracelet/crush/internal/tui/components/dialogs"
	"github.com/charmbracelet/crush/internal/tui/exp/list"
	"github.com/charmbracelet/crush/internal/tui/styles"
	"github.com/charmbracelet/crush/internal/tui/util"
)

const (
	ThemeDialogID dialogs.DialogID = "theme"

	defaultWidth int = 40
)

type listModel = list.FilterableList[list.CompletionItem[ThemeOption]]

// ThemeOption represents a theme choice.
type ThemeOption struct {
	Name    string
	Display string
}

// ThemeDialog represents the theme selection dialog.
type ThemeDialog interface {
	dialogs.DialogModel
}

type themeDialogCmp struct {
	width   int
	wWidth  int
	wHeight int

	themeList listModel
	keyMap    ThemeDialogKeyMap
	help      help.Model
}

// ThemeSelectedMsg is sent when a theme is selected.
type ThemeSelectedMsg struct {
	Theme string
}

// ThemeDialogKeyMap defines key bindings for the theme dialog.
type ThemeDialogKeyMap struct {
	Next     key.Binding
	Previous key.Binding
	Select   key.Binding
	Close    key.Binding
}

// DefaultThemeDialogKeyMap returns the default key bindings.
func DefaultThemeDialogKeyMap() ThemeDialogKeyMap {
	return ThemeDialogKeyMap{
		Next: key.NewBinding(
			key.WithKeys("down", "j", "ctrl+n"),
			key.WithHelp("↓/j", "next"),
		),
		Previous: key.NewBinding(
			key.WithKeys("up", "k", "ctrl+p"),
			key.WithHelp("↑/k", "previous"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Close: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc", "close"),
		),
	}
}

// ShortHelp returns key bindings for the short help view.
func (k ThemeDialogKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Select, k.Close}
}

// FullHelp returns key bindings for the full help view.
func (k ThemeDialogKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Previous},
		{k.Select, k.Close},
	}
}

// NewThemeDialog creates a new theme selection dialog.
func NewThemeDialog() ThemeDialog {
	keyMap := DefaultThemeDialogKeyMap()
	listKeyMap := list.DefaultKeyMap()
	listKeyMap.Down.SetEnabled(false)
	listKeyMap.Up.SetEnabled(false)
	listKeyMap.DownOneItem = keyMap.Next
	listKeyMap.UpOneItem = keyMap.Previous

	t := styles.CurrentTheme()
	inputStyle := t.S().Base.PaddingLeft(1).PaddingBottom(1)
	themeList := list.NewFilterableList(
		[]list.CompletionItem[ThemeOption]{},
		list.WithFilterInputStyle(inputStyle),
		list.WithFilterListOptions(
			list.WithKeyMap(listKeyMap),
			list.WithWrapNavigation(),
			list.WithResizeByList(),
		),
	)
	h := help.New()
	h.Styles = t.S().Help

	return &themeDialogCmp{
		themeList: themeList,
		width:     defaultWidth,
		keyMap:    keyMap,
		help:      h,
	}
}

func (d *themeDialogCmp) Init() tea.Cmd {
	return d.populateThemes()
}

func (d *themeDialogCmp) populateThemes() tea.Cmd {
	mgr := styles.DefaultManager()
	currentTheme := mgr.Current().Name
	themeNames := mgr.List()
	slices.Sort(themeNames)

	themeItems := []list.CompletionItem[ThemeOption]{}
	selectedID := ""
	for _, name := range themeNames {
		display := name
		opts := []list.CompletionItemOption{
			list.WithCompletionID(name),
		}
		if name == currentTheme {
			opts = append(opts, list.WithCompletionShortcut("current"))
			selectedID = name
		}
		themeItems = append(themeItems, list.NewCompletionItem(
			display,
			ThemeOption{Name: name, Display: display},
			opts...,
		))
	}

	cmd := d.themeList.SetItems(themeItems)
	if selectedID != "" {
		return tea.Sequence(cmd, d.themeList.SetSelected(selectedID))
	}
	return cmd
}

func (d *themeDialogCmp) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.wWidth = msg.Width
		d.wHeight = msg.Height
		return d, d.themeList.SetSize(d.listWidth(), d.listHeight())
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, d.keyMap.Select):
			selectedItem := d.themeList.SelectedItem()
			if selectedItem == nil {
				return d, nil
			}
			theme := (*selectedItem).Value()
			return d, tea.Sequence(
				util.CmdHandler(dialogs.CloseDialogMsg{}),
				func() tea.Msg {
					return ThemeSelectedMsg{Theme: theme.Name}
				},
			)
		case key.Matches(msg, d.keyMap.Close):
			return d, util.CmdHandler(dialogs.CloseDialogMsg{})
		default:
			u, cmd := d.themeList.Update(msg)
			d.themeList = u.(listModel)
			return d, cmd
		}
	}
	return d, nil
}

func (d *themeDialogCmp) View() string {
	t := styles.CurrentTheme()
	listView := d.themeList

	header := t.S().Base.Padding(0, 1, 1, 1).Render(core.Title("Select Theme", d.width-4))
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		listView.View(),
		"",
		t.S().Base.Width(d.width-2).PaddingLeft(1).AlignHorizontal(lipgloss.Left).Render(d.help.View(d.keyMap)),
	)
	return d.style().Render(content)
}

func (d *themeDialogCmp) Cursor() *tea.Cursor {
	if cursor, ok := d.themeList.(util.Cursor); ok {
		cursor := cursor.Cursor()
		if cursor != nil {
			cursor = d.moveCursor(cursor)
		}
		return cursor
	}
	return nil
}

func (d *themeDialogCmp) listWidth() int {
	return d.width - 2
}

func (d *themeDialogCmp) listHeight() int {
	listHeight := len(d.themeList.Items()) + 2 + 4
	return min(listHeight, d.wHeight/2)
}

func (d *themeDialogCmp) moveCursor(cursor *tea.Cursor) *tea.Cursor {
	row, col := d.Position()
	offset := row + 3
	cursor.Y += offset
	cursor.X = cursor.X + col + 2
	return cursor
}

func (d *themeDialogCmp) style() lipgloss.Style {
	t := styles.CurrentTheme()
	return t.S().Base.
		Width(d.width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderFocus)
}

func (d *themeDialogCmp) Position() (int, int) {
	row := d.wHeight/4 - 2
	col := d.wWidth / 2
	col -= d.width / 2
	return row, col
}

func (d *themeDialogCmp) ID() dialogs.DialogID {
	return ThemeDialogID
}
