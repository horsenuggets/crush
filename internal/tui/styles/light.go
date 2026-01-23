package styles

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// Light theme colors.
var (
	lightBgBase        = color.RGBA{R: 255, G: 255, B: 255, A: 255} // White
	lightBgBaseLighter = color.RGBA{R: 250, G: 250, B: 250, A: 255} // Near white
	lightBgSubtle      = color.RGBA{R: 240, G: 240, B: 240, A: 255} // Light gray
	lightBgOverlay     = color.RGBA{R: 230, G: 230, B: 230, A: 255} // Slightly darker

	lightFgBase      = color.RGBA{R: 32, G: 32, B: 32, A: 255}    // Near black
	lightFgMuted     = color.RGBA{R: 100, G: 100, B: 100, A: 255} // Dark gray
	lightFgHalfMuted = color.RGBA{R: 80, G: 80, B: 80, A: 255}    // Medium gray
	lightFgSubtle    = color.RGBA{R: 140, G: 140, B: 140, A: 255} // Light gray text
	lightFgSelected  = color.RGBA{R: 255, G: 255, B: 255, A: 255} // White for selection

	lightBorder      = color.RGBA{R: 210, G: 210, B: 210, A: 255} // Border gray
	lightBorderFocus = color.RGBA{R: 59, G: 130, B: 246, A: 255}  // Blue

	lightPrimary   = color.RGBA{R: 59, G: 130, B: 246, A: 255}  // Blue
	lightSecondary = color.RGBA{R: 100, G: 100, B: 100, A: 255} // Gray
	lightTertiary  = color.RGBA{R: 80, G: 80, B: 80, A: 255}    // Dark gray
	lightAccent    = color.RGBA{R: 59, G: 130, B: 246, A: 255}  // Blue

	lightSuccess = color.RGBA{R: 34, G: 150, B: 84, A: 255}   // Green
	lightError   = color.RGBA{R: 220, G: 53, B: 69, A: 255}   // Red
	lightWarning = color.RGBA{R: 200, G: 150, B: 30, A: 255}  // Yellow/orange
	lightInfo    = color.RGBA{R: 59, G: 130, B: 246, A: 255}  // Blue

	lightWhite     = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	lightBlueLight = color.RGBA{R: 147, G: 197, B: 253, A: 255}
	lightBlueDark  = color.RGBA{R: 30, G: 64, B: 175, A: 255}
	lightBlue      = color.RGBA{R: 59, G: 130, B: 246, A: 255}

	lightYellow = color.RGBA{R: 202, G: 138, B: 4, A: 255}
	lightCitron = color.RGBA{R: 234, G: 179, B: 8, A: 255}

	lightGreen      = color.RGBA{R: 34, G: 150, B: 84, A: 255}
	lightGreenDark  = color.RGBA{R: 22, G: 101, B: 52, A: 255}
	lightGreenLight = color.RGBA{R: 74, G: 222, B: 128, A: 255}

	lightRed      = color.RGBA{R: 220, G: 53, B: 69, A: 255}
	lightRedDark  = color.RGBA{R: 153, G: 27, B: 27, A: 255}
	lightRedLight = color.RGBA{R: 248, G: 113, B: 113, A: 255}
	lightCherry   = color.RGBA{R: 190, G: 24, B: 93, A: 255}
)

// NewLightTheme creates a light theme suitable for bright environments.
func NewLightTheme() *Theme {
	t := &Theme{
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
