// Package instances provides UI components for displaying parallel Crush instances.
package instances

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	crushctx "github.com/charmbracelet/crush/internal/context"
	"github.com/charmbracelet/crush/internal/home"
	"github.com/charmbracelet/crush/internal/tui/styles"
)

// RenderOptions configures how the instances block is rendered.
type RenderOptions struct {
	MaxWidth    int
	MaxItems    int
	ShowSection bool
	SectionName string
}

// RenderInstancesBlock renders a block showing other active Crush instances.
func RenderInstancesBlock(instances []crushctx.Instance, opts RenderOptions, isCompact bool) string {
	if len(instances) == 0 {
		return ""
	}

	t := styles.CurrentTheme()
	var lines []string

	if opts.ShowSection && opts.SectionName != "" {
		lines = append(lines, opts.SectionName)
	}

	// Limit items if needed
	displayInstances := instances
	if opts.MaxItems > 0 && len(instances) > opts.MaxItems {
		displayInstances = instances[:opts.MaxItems]
	}

	for _, inst := range displayInstances {
		line := renderInstance(inst, opts.MaxWidth, isCompact)
		lines = append(lines, line)
	}

	// Show count if we had to truncate
	if opts.MaxItems > 0 && len(instances) > opts.MaxItems {
		remaining := len(instances) - opts.MaxItems
		lines = append(lines, t.S().Muted.Render(fmt.Sprintf("  +%d more", remaining)))
	}

	return strings.Join(lines, "\n")
}

// renderInstance renders a single instance line.
func renderInstance(inst crushctx.Instance, maxWidth int, isCompact bool) string {
	t := styles.CurrentTheme()

	// Status indicator (green for active, yellow if slightly stale)
	age := time.Since(inst.LastActivity)
	var statusIcon lipgloss.Style
	if age < 30*time.Second {
		statusIcon = t.ItemOnlineIcon
	} else if age < 2*time.Minute {
		statusIcon = t.ItemBusyIcon
	} else {
		statusIcon = t.ItemOfflineIcon
	}

	// Format the working directory (shorten home paths)
	dir := inst.WorkingDir
	homeDir := home.Dir()
	if strings.HasPrefix(dir, homeDir) {
		dir = "~" + strings.TrimPrefix(dir, homeDir)
	}
	// Get just the base directory name for compact display
	baseName := filepath.Base(dir)

	// Build the line
	idStyle := t.S().Muted
	dirStyle := t.S().Text
	taskStyle := t.S().Subtle

	var line string
	if isCompact {
		// Compact: just icon + id + base dir
		line = fmt.Sprintf("%s %s %s",
			statusIcon.String(),
			idStyle.Render(inst.ID),
			dirStyle.Render(baseName),
		)
	} else {
		// Full: icon + id + dir + task
		idPart := idStyle.Render(inst.ID)
		dirPart := dirStyle.Render(truncateString(dir, maxWidth-15))

		if inst.Task != "" {
			taskPart := taskStyle.Render(truncateString(inst.Task, maxWidth-20))
			line = fmt.Sprintf("%s %s %s\n   %s",
				statusIcon.String(),
				idPart,
				dirPart,
				taskPart,
			)
		} else {
			line = fmt.Sprintf("%s %s %s",
				statusIcon.String(),
				idPart,
				dirPart,
			)
		}
	}

	return line
}

// truncateString truncates a string to maxLen, adding ellipsis if needed.
func truncateString(s string, maxLen int) string {
	if maxLen <= 3 {
		return s
	}
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
