package themes

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/tui/styles"
)

// Dark theme colors - minimal dark theme with neutral colors.
var (
	// Backgrounds (dark to light)
	darkBgBase        = styles.ParseHex("#201F26") // Pepper
	darkBgBaseLighter = styles.ParseHex("#2d2c35") // BBQ
	darkBgSubtle      = styles.ParseHex("#3A3943") // Charcoal
	darkBgOverlay     = styles.ParseHex("#4D4C57") // Iron

	// Foregrounds (dark to light)
	darkFgSubtle    = styles.ParseHex("#605F6B") // Oyster
	darkFgMuted     = styles.ParseHex("#858392") // Squid
	darkFgHalfMuted = styles.ParseHex("#BFBCC8") // Smoke
	darkFgBase      = styles.ParseHex("#DFDBDD") // Ash
	darkFgSelected  = styles.ParseHex("#F1EFEF") // Salt

	// Accent colors
	darkPrimary   = styles.ParseHex("#00A4FF") // Malibu
	darkSecondary = styles.ParseHex("#00A4FF") // Malibu
	darkTertiary  = styles.ParseHex("#DFDBDD") // Ash
	darkAccent    = styles.ParseHex("#00A4FF") // Malibu

	// Status colors
	darkSuccess = styles.ParseHex("#12C78F") // Guac
	darkError   = styles.ParseHex("#FF577D") // Coral
	darkWarning = styles.ParseHex("#E8FF27") // Citron
	darkInfo    = styles.ParseHex("#00A4FF") // Malibu

	// Named colors
	darkWhite     = styles.ParseHex("#F1EFEF") // Salt
	darkBlueLight = styles.ParseHex("#4FBEFE") // Sardine
	darkBlueDark  = styles.ParseHex("#007AB8") // Damson
	darkBlue      = styles.ParseHex("#00A4FF") // Malibu
	darkYellow    = styles.ParseHex("#F5EF34") // Mustard
	darkCitron    = styles.ParseHex("#E8FF27") // Citron
	darkGreen     = styles.ParseHex("#00FFB2") // Julep
	darkGreenDark = styles.ParseHex("#12C78F") // Guac
	darkGreenLt   = styles.ParseHex("#68FFD6") // Bok
	darkRed       = styles.ParseHex("#FF577D") // Coral
	darkRedDark   = styles.ParseHex("#EB4268") // Sriracha
	darkRedLight  = styles.ParseHex("#FF7F90") // Salmon
	darkCherry    = styles.ParseHex("#FF388B") // Cherry

	// Border colors
	darkBorder      = styles.ParseHex("#3A3943") // Charcoal
	darkBorderFocus = styles.ParseHex("#00A4FF") // Malibu
)

// NewDarkTheme creates a minimal dark theme with neutral colors.
func NewDarkTheme() *styles.Theme {
	t := &styles.Theme{
		Name:   "dark",
		IsDark: true,

		Primary:   darkPrimary,
		Secondary: darkSecondary,
		Tertiary:  darkTertiary,
		Accent:    darkAccent,

		// Backgrounds
		BgBase:        darkBgBase,
		BgBaseLighter: darkBgBaseLighter,
		BgSubtle:      darkBgSubtle,
		BgOverlay:     darkBgOverlay,

		// Foregrounds
		FgBase:      darkFgBase,
		FgMuted:     darkFgMuted,
		FgHalfMuted: darkFgHalfMuted,
		FgSubtle:    darkFgSubtle,
		FgSelected: darkBgBase,

		// Borders
		Border:      darkBorder,
		BorderFocus: darkBorderFocus,

		// Status
		Success: darkSuccess,
		Error:   darkError,
		Warning: darkWarning,
		Info:    darkInfo,

		// Colors
		White: darkWhite,

		BlueLight: darkBlueLight,
		BlueDark:  darkBlueDark,
		Blue:      darkBlue,

		Yellow: darkYellow,
		Citron: darkCitron,

		Green:      darkGreen,
		GreenDark:  darkGreenDark,
		GreenLight: darkGreenLt,

		Red:      darkRed,
		RedDark:  darkRedDark,
		RedLight: darkRedLight,
		Cherry:   darkCherry,
	}

	// Text selection.
	t.TextSelection = lipgloss.NewStyle().Foreground(darkFgSelected).Background(darkPrimary)

	// LSP and MCP status.
	t.ItemOfflineIcon = lipgloss.NewStyle().Foreground(darkFgSubtle).SetString("●")
	t.ItemBusyIcon = t.ItemOfflineIcon.Foreground(darkWarning)
	t.ItemErrorIcon = t.ItemOfflineIcon.Foreground(darkError)
	t.ItemOnlineIcon = t.ItemOfflineIcon.Foreground(darkSuccess)

	// Editor: Yolo Mode.
	t.YoloIconFocused = lipgloss.NewStyle().Foreground(darkBgBase).Background(darkWarning).Bold(true).SetString(" ! ")
	t.YoloIconBlurred = t.YoloIconFocused.Foreground(darkBgBase).Background(darkFgSubtle)
	t.YoloDotsFocused = lipgloss.NewStyle().Foreground(darkWarning).SetString(":::")
	t.YoloDotsBlurred = t.YoloDotsFocused.Foreground(darkFgSubtle)

	// oAuth Chooser.
	t.AuthBorderSelected = lipgloss.NewStyle().BorderForeground(darkAccent)
	t.AuthTextSelected = lipgloss.NewStyle().Foreground(darkAccent)
	t.AuthBorderUnselected = lipgloss.NewStyle().BorderForeground(darkPrimary)
	t.AuthTextUnselected = lipgloss.NewStyle().Foreground(darkFgSubtle)

	return t
}
