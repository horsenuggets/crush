package files

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/charmbracelet/crush/internal/history"
	"github.com/charmbracelet/crush/internal/tui/components/core"
	"github.com/charmbracelet/crush/internal/tui/styles"
)

// FileHistory represents a file history with initial and latest versions.
type FileHistory struct {
	InitialVersion history.File
	LatestVersion  history.File
}

// SessionFile represents a file with its history information.
type SessionFile struct {
	History   FileHistory
	FilePath  string
	Additions int
	Deletions int
}

// RenderOptions contains options for rendering file lists.
type RenderOptions struct {
	MaxWidth    int
	MaxItems    int
	ShowSection bool
	SectionName string
}

// RenderFileList renders a list of file status items with the given options.
func RenderFileList(fileSlice []SessionFile, opts RenderOptions) []string {
	t := styles.CurrentTheme()
	fileList := []string{}

	if opts.ShowSection {
		sectionName := opts.SectionName
		if sectionName == "" {
			sectionName = "Modified Files"
		}
		section := t.S().Subtle.Render(sectionName)
		fileList = append(fileList, section, "")
	}

	if len(fileSlice) == 0 {
		fileList = append(fileList, t.S().Base.Foreground(t.Border).Render("None"))
		return fileList
	}

	// Sort files alphabetically by filename
	sort.Slice(fileSlice, func(i, j int) bool {
		return strings.Compare(filepath.Base(fileSlice[i].FilePath), filepath.Base(fileSlice[j].FilePath)) < 0
	})

	// Build a map of filenames to detect duplicates
	filenameCounts := make(map[string]int)
	for _, file := range fileSlice {
		if file.Additions == 0 && file.Deletions == 0 {
			continue
		}
		filename := filepath.Base(file.FilePath)
		filenameCounts[filename]++
	}

	// Determine how many items to show
	maxItems := len(fileSlice)
	if opts.MaxItems > 0 {
		maxItems = min(opts.MaxItems, len(fileSlice))
	}

	filesShown := 0
	for _, file := range fileSlice {
		if file.Additions == 0 && file.Deletions == 0 {
			continue // skip files with no changes
		}
		if filesShown >= maxItems {
			break
		}

		var statusParts []string
		if file.Additions > 0 {
			statusParts = append(statusParts, t.S().Base.Foreground(t.Success).Render(fmt.Sprintf("+%d", file.Additions)))
		}
		if file.Deletions > 0 {
			statusParts = append(statusParts, t.S().Base.Foreground(t.Error).Render(fmt.Sprintf("-%d", file.Deletions)))
		}

		extraContent := strings.Join(statusParts, " ")

		// Default to just the filename
		filename := filepath.Base(file.FilePath)
		displayPath := filename

		// If there are duplicate filenames, add parent directory for disambiguation
		if filenameCounts[filename] > 1 {
			parentDir := filepath.Base(filepath.Dir(file.FilePath))
			displayPath = parentDir + "/" + filename
		}

		displayPath = ansi.Truncate(displayPath, opts.MaxWidth-lipgloss.Width(extraContent)-2, "…")

		fileList = append(fileList,
			core.Status(
				core.StatusOpts{
					Title:        displayPath,
					ExtraContent: extraContent,
				},
				opts.MaxWidth,
			),
		)
		filesShown++
	}

	return fileList
}

// RenderFileBlock renders a complete file block with optional truncation indicator.
func RenderFileBlock(fileSlice []SessionFile, opts RenderOptions, showTruncationIndicator bool) string {
	t := styles.CurrentTheme()
	fileList := RenderFileList(fileSlice, opts)

	// Add truncation indicator if needed
	if showTruncationIndicator && opts.MaxItems > 0 {
		totalFilesWithChanges := 0
		for _, file := range fileSlice {
			if file.Additions > 0 || file.Deletions > 0 {
				totalFilesWithChanges++
			}
		}
		if totalFilesWithChanges > opts.MaxItems {
			remaining := totalFilesWithChanges - opts.MaxItems
			if remaining == 1 {
				fileList = append(fileList, t.S().Base.Foreground(t.FgMuted).Render("…"))
			} else {
				fileList = append(fileList,
					t.S().Base.Foreground(t.FgSubtle).Render(fmt.Sprintf("…and %d more", remaining)),
				)
			}
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, fileList...)
	if opts.MaxWidth > 0 {
		return lipgloss.NewStyle().Width(opts.MaxWidth).Render(content)
	}
	return content
}
