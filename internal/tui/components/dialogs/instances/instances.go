package instances

import (
	"fmt"
	"path/filepath"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	crushctx "github.com/charmbracelet/crush/internal/context"
	"github.com/charmbracelet/crush/internal/tui/components/core"
	"github.com/charmbracelet/crush/internal/tui/components/dialogs"
	"github.com/charmbracelet/crush/internal/tui/styles"
	"github.com/charmbracelet/crush/internal/tui/util"
)

const InstancesDialogID dialogs.DialogID = "instances"

// InstancesDialog interface for the parallel instances dialog.
type InstancesDialog interface {
	dialogs.DialogModel
}

type instancesDialogCmp struct {
	instances   []crushctx.Instance
	currentID   string
	selectedIdx int
	wWidth      int
	wHeight     int
	width       int
	keyMap      KeyMap
	help        help.Model
}

// NewInstancesDialog creates a new dialog showing parallel Crush instances.
func NewInstancesDialog(instances []crushctx.Instance, currentID string) InstancesDialog {
	t := styles.CurrentTheme()
	help := help.New()
	help.Styles = t.S().Help

	return &instancesDialogCmp{
		instances:   instances,
		currentID:   currentID,
		selectedIdx: 0,
		keyMap:      DefaultKeyMap(),
		help:        help,
	}
}

func (d *instancesDialogCmp) Init() tea.Cmd {
	return nil
}

func (d *instancesDialogCmp) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.wWidth = msg.Width
		d.wHeight = msg.Height
		d.width = min(80, d.wWidth-8)
		return d, nil
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, d.keyMap.Close):
			return d, util.CmdHandler(dialogs.CloseDialogMsg{})
		case key.Matches(msg, d.keyMap.Up):
			if d.selectedIdx > 0 {
				d.selectedIdx--
			}
			return d, nil
		case key.Matches(msg, d.keyMap.Down):
			if d.selectedIdx < len(d.instances)-1 {
				d.selectedIdx++
			}
			return d, nil
		}
	}
	return d, nil
}

func (d *instancesDialogCmp) View() string {
	t := styles.CurrentTheme()

	var content string
	if len(d.instances) == 0 {
		content = t.S().Muted.Padding(1).Render("No other Crush instances running")
	} else {
		var rows []string
		for i, inst := range d.instances {
			row := d.renderInstance(inst, i == d.selectedIdx)
			rows = append(rows, row)
		}
		content = lipgloss.JoinVertical(lipgloss.Left, rows...)
	}

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		t.S().Base.Padding(0, 1, 1, 1).Render(core.Title("Parallel Instances", d.width-4)),
		content,
		"",
		t.S().Base.Width(d.width-2).PaddingLeft(1).AlignHorizontal(lipgloss.Left).Render(d.help.View(d.keyMap)),
	)

	return d.style().Render(view)
}

func (d *instancesDialogCmp) renderInstance(inst crushctx.Instance, selected bool) string {
	t := styles.CurrentTheme()

	// Determine if this is the current instance
	isCurrent := inst.ID == d.currentID
	currentMarker := ""
	if isCurrent {
		currentMarker = " (current)"
	}

	// Format time since last activity
	timeSince := time.Since(inst.LastActivity)
	var timeStr string
	if timeSince < time.Minute {
		timeStr = "just now"
	} else if timeSince < time.Hour {
		timeStr = fmt.Sprintf("%dm ago", int(timeSince.Minutes()))
	} else {
		timeStr = fmt.Sprintf("%dh ago", int(timeSince.Hours()))
	}

	// Shorten paths for display
	workDir := filepath.Base(inst.WorkingDir)
	if workDir == "." {
		workDir = inst.WorkingDir
	}

	// Build instance display
	idStyle := t.S().Text.Bold(true)
	if selected {
		idStyle = idStyle.Foreground(t.Primary)
	}

	taskStr := inst.Task
	if taskStr == "" {
		taskStr = "idle"
	}
	if len(taskStr) > 40 {
		taskStr = taskStr[:37] + "..."
	}

	line1 := fmt.Sprintf("%s%s", idStyle.Render(inst.ID), t.S().Muted.Render(currentMarker))
	line2 := t.S().Subtle.Render(fmt.Sprintf("  %s • %s", workDir, timeStr))
	line3 := t.S().Muted.Render(fmt.Sprintf("  %s", taskStr))

	rowStyle := t.S().Base.Padding(0, 1).Width(d.width - 4)
	if selected {
		rowStyle = rowStyle.Background(t.BgSubtle)
	}

	return rowStyle.Render(lipgloss.JoinVertical(lipgloss.Left, line1, line2, line3))
}

func (d *instancesDialogCmp) style() lipgloss.Style {
	t := styles.CurrentTheme()
	return t.S().Base.
		Width(d.width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		Background(t.BgBase).
		Padding(1)
}

func (d *instancesDialogCmp) Cursor() *tea.Cursor {
	return nil
}

func (d *instancesDialogCmp) ID() dialogs.DialogID {
	return InstancesDialogID
}

func (d *instancesDialogCmp) Position() (int, int) {
	// Center the dialog
	x := (d.wWidth - d.width) / 2
	y := (d.wHeight - 20) / 2
	if y < 2 {
		y = 2
	}
	return x, y
}
