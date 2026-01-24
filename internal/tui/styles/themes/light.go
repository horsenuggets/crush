package themes

import (
	"image/color"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/tui/styles"
)

// Light theme colors - pastel palette for bright environments.
var (
	lightBgBase        = color.RGBA{R: 255, G: 255, B: 255, A: 255} // White
	lightBgBaseLighter = color.RGBA{R: 250, G: 251, B: 252, A: 255} // Near white
	lightBgSubtle      = color.RGBA{R: 243, G: 244, B: 246, A: 255} // Very light gray
	lightBgOverlay     = color.RGBA{R: 229, G: 231, B: 235, A: 255} // Light gray

	lightFgBase      = color.RGBA{R: 31, G: 41, B: 55, A: 255}    // Dark gray (not pure black)
	lightFgMuted     = color.RGBA{R: 107, G: 114, B: 128, A: 255} // Medium gray
	lightFgHalfMuted = color.RGBA{R: 75, G: 85, B: 99, A: 255}    // Darker gray
	lightFgSubtle    = color.RGBA{R: 156, G: 163, B: 175, A: 255} // Light gray text
	lightFgSelected  = color.RGBA{R: 255, G: 255, B: 255, A: 255} // White for selection

	lightBorder      = color.RGBA{R: 229, G: 231, B: 235, A: 255} // Light border
	lightBorderFocus = color.RGBA{R: 99, G: 102, B: 241, A: 255}  // Indigo

	lightPrimary   = color.RGBA{R: 99, G: 102, B: 241, A: 255}  // Indigo
	lightSecondary = color.RGBA{R: 107, G: 114, B: 128, A: 255} // Gray
	lightTertiary  = color.RGBA{R: 75, G: 85, B: 99, A: 255}    // Dark gray
	lightAccent    = color.RGBA{R: 139, G: 92, B: 246, A: 255}  // Purple

	lightSuccess = color.RGBA{R: 74, G: 222, B: 128, A: 255}  // Pastel green
	lightError   = color.RGBA{R: 251, G: 113, B: 133, A: 255} // Pastel red/rose
	lightWarning = color.RGBA{R: 251, G: 191, B: 36, A: 255}  // Pastel amber
	lightInfo    = color.RGBA{R: 96, G: 165, B: 250, A: 255}  // Pastel blue

	lightWhite     = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	lightBlueLight = color.RGBA{R: 191, G: 219, B: 254, A: 255} // Very light blue
	lightBlueDark  = color.RGBA{R: 59, G: 130, B: 246, A: 255}  // Medium blue
	lightBlue      = color.RGBA{R: 96, G: 165, B: 250, A: 255}  // Pastel blue

	lightYellow = color.RGBA{R: 253, G: 224, B: 71, A: 255}  // Pastel yellow
	lightCitron = color.RGBA{R: 251, G: 191, B: 36, A: 255}  // Amber

	lightGreen      = color.RGBA{R: 74, G: 222, B: 128, A: 255}  // Pastel green
	lightGreenDark  = color.RGBA{R: 34, G: 197, B: 94, A: 255}   // Medium green
	lightGreenLight = color.RGBA{R: 134, G: 239, B: 172, A: 255} // Very light green

	lightRed      = color.RGBA{R: 251, G: 113, B: 133, A: 255} // Pastel rose
	lightRedDark  = color.RGBA{R: 244, G: 63, B: 94, A: 255}   // Rose
	lightRedLight = color.RGBA{R: 253, G: 164, B: 175, A: 255} // Very light rose
	lightCherry   = color.RGBA{R: 236, G: 72, B: 153, A: 255}  // Pink
)

// NewLightTheme creates a light theme suitable for bright environments.
func NewLightTheme() *styles.Theme {
	t := &styles.Theme{
		Name:   "light",
		IsDark: false,

		Primary:   lightPrimary,
		Secondary: lightSecondary,
		Tertiary:  lightTertiary,
		Accent:    lightAccent,

		// Backgrounds
		BgBase:        lightBgBase,
		BgBaseLighter: lightBgBaseLighter,
		BgSubtle:      lightBgSubtle,
		BgOverlay:     lightBgOverlay,

		// Foregrounds
		FgBase:      lightFgBase,
		FgMuted:     lightFgMuted,
		FgHalfMuted: lightFgHalfMuted,
		FgSubtle:    lightFgSubtle,
		FgSelected:  lightFgSelected,

		// Borders
		Border:      lightBorder,
		BorderFocus: lightBorderFocus,

		// Status
		Success: lightSuccess,
		Error:   lightError,
		Warning: lightWarning,
		Info:    lightInfo,

		// Colors
		White: lightWhite,

		BlueLight: lightBlueLight,
		BlueDark:  lightBlueDark,
		Blue:      lightBlue,

		Yellow: lightYellow,
		Citron: lightCitron,

		Green:      lightGreen,
		GreenDark:  lightGreenDark,
		GreenLight: lightGreenLight,

		Red:      lightRed,
		RedDark:  lightRedDark,
		RedLight: lightRedLight,
		Cherry:   lightCherry,
	}

	// Text selection.
	t.TextSelection = lipgloss.NewStyle().Foreground(lightFgSelected).Background(lightPrimary)

	// LSP and MCP status.
	t.ItemOfflineIcon = lipgloss.NewStyle().Foreground(lightFgSubtle).SetString("●")
	t.ItemBusyIcon = t.ItemOfflineIcon.Foreground(lightWarning)
	t.ItemErrorIcon = t.ItemOfflineIcon.Foreground(lightError)
	t.ItemOnlineIcon = t.ItemOfflineIcon.Foreground(lightSuccess)

	// Editor: Yolo Mode.
	t.YoloIconFocused = lipgloss.NewStyle().Foreground(lightFgBase).Background(lightWarning).Bold(true).SetString(" ! ")
	t.YoloIconBlurred = t.YoloIconFocused.Foreground(lightFgMuted).Background(lightBgOverlay)
	t.YoloDotsFocused = lipgloss.NewStyle().Foreground(lightWarning).SetString(":::")
	t.YoloDotsBlurred = t.YoloDotsFocused.Foreground(lightFgSubtle)

	// oAuth Chooser.
	t.AuthBorderSelected = lipgloss.NewStyle().BorderForeground(lightSuccess)
	t.AuthTextSelected = lipgloss.NewStyle().Foreground(lightSuccess)
	t.AuthBorderUnselected = lipgloss.NewStyle().BorderForeground(lightBorder)
	t.AuthTextUnselected = lipgloss.NewStyle().Foreground(lightFgMuted)

	return t
}
