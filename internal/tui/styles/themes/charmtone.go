package themes

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/tui/styles"
)

// Charmtone theme colors - official Charm color palette.
var (
	// Backgrounds (dark to light)
	charmtoneBgBase        = styles.ParseHex("#201F26") // Pepper
	charmtoneBgBaseLighter = styles.ParseHex("#2d2c35") // BBQ
	charmtoneBgSubtle      = styles.ParseHex("#3A3943") // Charcoal
	charmtoneBgOverlay     = styles.ParseHex("#4D4C57") // Iron

	// Foregrounds (dark to light)
	charmtoneFgSubtle    = styles.ParseHex("#605F6B") // Oyster
	charmtoneFgMuted     = styles.ParseHex("#858392") // Squid
	charmtoneFgHalfMuted = styles.ParseHex("#BFBCC8") // Smoke
	charmtoneFgBase      = styles.ParseHex("#DFDBDD") // Ash
	charmtoneFgSelected  = styles.ParseHex("#F1EFEF") // Salt

	// Accent colors
	charmtonePrimary   = styles.ParseHex("#6B50FF") // Charple
	charmtoneSecondary = styles.ParseHex("#F379FF") // Lilac
	charmtoneTertiary  = styles.ParseHex("#68FFD6") // Bok
	charmtoneAccent    = styles.ParseHex("#E8FE96") // Zest

	// Status colors
	charmtoneSuccess = styles.ParseHex("#12C78F") // Guac
	charmtoneError   = styles.ParseHex("#EB4268") // Sriracha
	charmtoneWarning = styles.ParseHex("#E8FE96") // Zest
	charmtoneInfo    = styles.ParseHex("#00A4FF") // Malibu

	// Named colors
	charmtoneWhite     = styles.ParseHex("#FFFAF1") // Butter
	charmtoneBlueLight = styles.ParseHex("#4FBEFE") // Sardine
	charmtoneBlueDark  = styles.ParseHex("#007AB8") // Damson
	charmtoneBlue      = styles.ParseHex("#00A4FF") // Malibu
	charmtoneYellow    = styles.ParseHex("#F5EF34") // Mustard
	charmtoneCitron    = styles.ParseHex("#E8FF27") // Citron
	charmtoneGreen     = styles.ParseHex("#00FFB2") // Julep
	charmtoneGreenDark = styles.ParseHex("#12C78F") // Guac
	charmtoneGreenLt   = styles.ParseHex("#68FFD6") // Bok
	charmtoneRed       = styles.ParseHex("#FF577D") // Coral
	charmtoneRedDark   = styles.ParseHex("#EB4268") // Sriracha
	charmtoneRedLight  = styles.ParseHex("#FF7F90") // Salmon
	charmtoneCherry    = styles.ParseHex("#FF388B") // Cherry

	// Border colors
	charmtoneBorder      = styles.ParseHex("#3A3943") // Charcoal
	charmtoneBorderFocus = styles.ParseHex("#6B50FF") // Charple
)

// NewCharmtoneTheme creates a theme using the official Charm color palette.
func NewCharmtoneTheme() *styles.Theme {
	t := &styles.Theme{
		Name:   "charmtone",
		IsDark: true,

		Primary:   charmtonePrimary,
		Secondary: charmtoneSecondary,
		Tertiary:  charmtoneTertiary,
		Accent:    charmtoneAccent,

		// Backgrounds
		BgBase:        charmtoneBgBase,
		BgBaseLighter: charmtoneBgBaseLighter,
		BgSubtle:      charmtoneBgSubtle,
		BgOverlay:     charmtoneBgOverlay,

		// Foregrounds
		FgBase:      charmtoneFgBase,
		FgMuted:     charmtoneFgMuted,
		FgHalfMuted: charmtoneFgHalfMuted,
		FgSubtle:    charmtoneFgSubtle,
		FgSelected: charmtoneFgSelected,

		// Borders
		Border:      charmtoneBorder,
		BorderFocus: charmtoneBorderFocus,

		// Status
		Success: charmtoneSuccess,
		Error:   charmtoneError,
		Warning: charmtoneWarning,
		Info:    charmtoneInfo,

		// Colors
		White: charmtoneWhite,

		BlueLight: charmtoneBlueLight,
		BlueDark:  charmtoneBlueDark,
		Blue:      charmtoneBlue,

		Yellow: charmtoneYellow,
		Citron: charmtoneCitron,

		Green:      charmtoneGreen,
		GreenDark:  charmtoneGreenDark,
		GreenLight: charmtoneGreenLt,

		Red:      charmtoneRed,
		RedDark:  charmtoneRedDark,
		RedLight: charmtoneRedLight,
		Cherry:   charmtoneCherry,
	}

	// Text selection.
	t.TextSelection = lipgloss.NewStyle().Foreground(charmtoneFgSelected).Background(charmtonePrimary)

	// LSP and MCP status.
	t.ItemOfflineIcon = lipgloss.NewStyle().Foreground(charmtoneFgMuted).SetString("●")
	t.ItemBusyIcon = t.ItemOfflineIcon.Foreground(charmtoneCitron)
	t.ItemErrorIcon = t.ItemOfflineIcon.Foreground(charmtoneRed)
	t.ItemOnlineIcon = t.ItemOfflineIcon.Foreground(charmtoneSuccess)

	// Editor: Yolo Mode.
	t.YoloIconFocused = lipgloss.NewStyle().Foreground(charmtoneFgSubtle).Background(charmtoneCitron).Bold(true).SetString(" ! ")
	t.YoloIconBlurred = t.YoloIconFocused.Foreground(charmtoneBgBase).Background(charmtoneFgMuted)
	t.YoloDotsFocused = lipgloss.NewStyle().Foreground(charmtoneAccent).SetString(":::")
	t.YoloDotsBlurred = t.YoloDotsFocused.Foreground(charmtoneFgMuted)

	// oAuth Chooser.
	t.AuthBorderSelected = lipgloss.NewStyle().BorderForeground(charmtoneSuccess)
	t.AuthTextSelected = lipgloss.NewStyle().Foreground(charmtoneGreen)
	t.AuthBorderUnselected = lipgloss.NewStyle().BorderForeground(charmtoneBgOverlay)
	t.AuthTextUnselected = lipgloss.NewStyle().Foreground(charmtoneFgMuted)

	return t
}
