package themes

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/tui/styles"
)

// Ocean theme - deep blues and teals inspired by the sea.
var (
	oceanBgBase        = styles.ParseHex("#0b132b") // Deep navy
	oceanBgBaseLighter = styles.ParseHex("#111c38") // Slightly lighter navy
	oceanBgSubtle      = styles.ParseHex("#172646") // Midnight blue
	oceanBgOverlay     = styles.ParseHex("#1e325a") // Ocean blue

	oceanFgBase      = styles.ParseHex("#e0eefa") // Seafoam white
	oceanFgMuted     = styles.ParseHex("#8caac8") // Muted blue-gray
	oceanFgHalfMuted = styles.ParseHex("#a0bed7") // Light blue-gray
	oceanFgSubtle    = styles.ParseHex("#536D89") // Subtle blue
	oceanFgSelected  = styles.ParseHex("#0b132b") // Dark navy for contrast on teal selection

	oceanBorder      = styles.ParseHex("#283c64") // Deep border
	oceanBorderFocus = styles.ParseHex("#40bebc") // Teal

	oceanPrimary   = styles.ParseHex("#40bebc") // Teal
	oceanSecondary = styles.ParseHex("#6495ed") // Cornflower blue
	oceanTertiary  = styles.ParseHex("#4682b4") // Steel blue
	oceanAccent    = styles.ParseHex("#00ced1") // Turquoise

	oceanSuccess = styles.ParseHex("#48d1b0") // Sea green
	oceanError   = styles.ParseHex("#ff6b81") // Coral
	oceanWarning = styles.ParseHex("#ffbe5c") // Sandy gold
	oceanInfo    = styles.ParseHex("#64b5f6") // Light blue

	oceanWhite       = styles.ParseHex("#f0f8ff") // Alice blue
	oceanBlueLight   = styles.ParseHex("#87cefa") // Light sky blue
	oceanBlueDark    = styles.ParseHex("#191970") // Midnight blue
	oceanBlue        = styles.ParseHex("#6495ed") // Cornflower
	oceanYoloFocused = styles.ParseHex("#64dcc8") // Aqua-tinted for YOLO

	oceanYellow = styles.ParseHex("#ffd764") // Sandy
	oceanCitron = styles.ParseHex("#ffbe5c") // Gold

	oceanGreen      = styles.ParseHex("#48d1b0") // Sea green
	oceanGreenDark  = styles.ParseHex("#20b2aa") // Light sea green
	oceanGreenLight = styles.ParseHex("#7fffd4") // Aquamarine

	oceanRed      = styles.ParseHex("#ff6b81") // Coral
	oceanRedDark  = styles.ParseHex("#cd5c5c") // Indian red
	oceanRedLight = styles.ParseHex("#ffa0a0") // Light coral
	oceanCherry   = styles.ParseHex("#ff69b4") // Hot pink
)

// NewOceanTheme creates an ocean-inspired deep blue theme.
func NewOceanTheme() *styles.Theme {
	t := &styles.Theme{
		Name:   "ocean",
		IsDark: true,

		Primary:   oceanPrimary,
		Secondary: oceanSecondary,
		Tertiary:  oceanTertiary,
		Accent:    oceanAccent,

		BgBase:        oceanBgBase,
		BgBaseLighter: oceanBgBaseLighter,
		BgSubtle:      oceanBgSubtle,
		BgOverlay:     oceanBgOverlay,

		FgBase:      oceanFgBase,
		FgMuted:     oceanFgMuted,
		FgHalfMuted: oceanFgHalfMuted,
		FgSubtle:    oceanFgSubtle,
		FgSelected:  oceanFgSelected,

		Border:      oceanBorder,
		BorderFocus: oceanBorderFocus,

		Success: oceanSuccess,
		Error:   oceanError,
		Warning: oceanWarning,
		Info:    oceanInfo,

		White:     oceanWhite,
		BlueLight: oceanBlueLight,
		BlueDark:  oceanBlueDark,
		Blue:      oceanBlue,

		Yellow: oceanYellow,
		Citron: oceanCitron,

		Green:      oceanGreen,
		GreenDark:  oceanGreenDark,
		GreenLight: oceanGreenLight,

		Red:      oceanRed,
		RedDark:  oceanRedDark,
		RedLight: oceanRedLight,
		Cherry:   oceanCherry,
	}

	t.TextSelection = lipgloss.NewStyle().Foreground(oceanFgSelected).Background(oceanPrimary)

	t.ItemOfflineIcon = lipgloss.NewStyle().Foreground(oceanFgSubtle).SetString("●")
	t.ItemBusyIcon = t.ItemOfflineIcon.Foreground(oceanWarning)
	t.ItemErrorIcon = t.ItemOfflineIcon.Foreground(oceanError)
	t.ItemOnlineIcon = t.ItemOfflineIcon.Foreground(oceanSuccess)

	// Editor: Yolo Mode - using aqua-tinted color to fit the ocean aesthetic.
	t.YoloIconFocused = lipgloss.NewStyle().Foreground(oceanBgBase).Background(oceanYoloFocused).Bold(true).SetString(" ! ")
	t.YoloIconBlurred = t.YoloIconFocused.Foreground(oceanBgBase).Background(oceanFgSubtle)
	t.YoloDotsFocused = lipgloss.NewStyle().Foreground(oceanYoloFocused).SetString(":::")
	t.YoloDotsBlurred = t.YoloDotsFocused.Foreground(oceanFgSubtle)

	t.AuthBorderSelected = lipgloss.NewStyle().BorderForeground(oceanSuccess)
	t.AuthTextSelected = lipgloss.NewStyle().Foreground(oceanSuccess)
	t.AuthBorderUnselected = lipgloss.NewStyle().BorderForeground(oceanBorder)
	t.AuthTextUnselected = lipgloss.NewStyle().Foreground(oceanFgMuted)

	return t
}
