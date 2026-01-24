package themes

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/tui/styles"
)

// Light theme colors - pastel palette for bright environments.
var (
	lightBgBase        = styles.ParseHex("#ffffff") // White
	lightBgBaseLighter = styles.ParseHex("#fafbfc") // Near white
	lightBgSubtle      = styles.ParseHex("#f3f4f6") // Very light gray
	lightBgOverlay     = styles.ParseHex("#e5e7eb") // Light gray

	lightFgBase      = styles.ParseHex("#1f2937") // Dark gray (not pure black)
	lightFgMuted     = styles.ParseHex("#6b7280") // Medium gray
	lightFgHalfMuted = styles.ParseHex("#4b5563") // Darker gray
	lightFgSubtle    = styles.ParseHex("#9ca3af") // Light gray text
	lightFgSelected  = styles.ParseHex("#ffffff") // White for selection

	lightBorder      = styles.ParseHex("#e5e7eb") // Light border
	lightBorderFocus = styles.ParseHex("#6366f1") // Indigo

	lightPrimary   = styles.ParseHex("#6366f1") // Indigo
	lightSecondary = styles.ParseHex("#6b7280") // Gray
	lightTertiary  = styles.ParseHex("#4b5563") // Dark gray
	lightAccent    = styles.ParseHex("#8b5cf6") // Purple

	lightSuccess = styles.ParseHex("#4ade80") // Pastel green
	lightError   = styles.ParseHex("#fb7185") // Pastel red/rose
	lightWarning = styles.ParseHex("#fbbf24") // Pastel amber
	lightInfo    = styles.ParseHex("#60a5fa") // Pastel blue

	lightWhite     = styles.ParseHex("#ffffff")
	lightBlueLight = styles.ParseHex("#bfdbfe") // Very light blue
	lightBlueDark  = styles.ParseHex("#3b82f6") // Medium blue
	lightBlue      = styles.ParseHex("#60a5fa") // Pastel blue

	lightYellow = styles.ParseHex("#fde047") // Pastel yellow
	lightCitron = styles.ParseHex("#fbbf24") // Amber

	lightGreen      = styles.ParseHex("#4ade80") // Pastel green
	lightGreenDark  = styles.ParseHex("#22c55e") // Medium green
	lightGreenLight = styles.ParseHex("#86efac") // Very light green

	lightRed      = styles.ParseHex("#fb7185") // Pastel rose
	lightRedDark  = styles.ParseHex("#f43f5e") // Rose
	lightRedLight = styles.ParseHex("#fda4af") // Very light rose
	lightCherry   = styles.ParseHex("#ec4899") // Pink
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
