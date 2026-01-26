package themes

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/tui/styles"
)

// Sugarcookie theme - a sweet pastel pink and purple light theme.
var (
	sugarBgBase        = styles.ParseHex("#ffedf8") // Soft pink white
	sugarBgBaseLighter = styles.ParseHex("#fffafd") // Near white with pink tint
	sugarBgSubtle      = styles.ParseHex("#fcf0f8") // Very light pink
	sugarBgOverlay     = styles.ParseHex("#f5e1ee") // Light pink overlay

	sugarFgBase      = styles.ParseHex("#503250") // Dark purple-gray
	sugarFgMuted     = styles.ParseHex("#8c6e8c") // Muted purple
	sugarFgHalfMuted = styles.ParseHex("#785578") // Medium purple
	sugarFgSubtle    = styles.ParseHex("#b496b4") // Subtle purple
	sugarFgSelected  = styles.ParseHex("#ffffff") // White for selection

	sugarBorder      = styles.ParseHex("#ebd2e6") // Light pink border
	sugarBorderFocus = styles.ParseHex("#c882c8") // Orchid focus

	sugarPrimary   = styles.ParseHex("#c882c8") // Orchid
	sugarSecondary = styles.ParseHex("#b48cc8") // Lavender
	sugarTertiary  = styles.ParseHex("#a078b4") // Soft purple
	sugarAccent    = styles.ParseHex("#ff96b4") // Soft pink

	sugarSuccess = styles.ParseHex("#96d2b4") // Mint green
	sugarError   = styles.ParseHex("#ff8ca0") // Soft rose
	sugarWarning = styles.ParseHex("#ffb4c8") // Pastel pink-peach
	sugarInfo    = styles.ParseHex("#aab4e6") // Periwinkle

	sugarWhite     = styles.ParseHex("#ffffff")
	sugarBlueLight = styles.ParseHex("#c8d2f5") // Lavender blue
	sugarBlueDark  = styles.ParseHex("#8282c8") // Medium purple-blue
	sugarBlue      = styles.ParseHex("#aab4e6") // Periwinkle

	sugarYellow = styles.ParseHex("#ffdcc8") // Peachy cream
	sugarCitron = styles.ParseHex("#ffc8b4") // Soft peach

	sugarGreen      = styles.ParseHex("#96d2b4") // Mint
	sugarGreenDark  = styles.ParseHex("#78b496") // Darker mint
	sugarGreenLight = styles.ParseHex("#b4e6c8") // Light mint

	sugarRed      = styles.ParseHex("#ff8ca0") // Soft rose
	sugarRedDark  = styles.ParseHex("#e66482") // Rose
	sugarRedLight = styles.ParseHex("#ffb4c3") // Light rose
	sugarCherry   = styles.ParseHex("#ff78aa") // Pink
)

// NewSugarcookieTheme creates a sweet pastel pink and purple light theme.
func NewSugarcookieTheme() *styles.Theme {
	t := &styles.Theme{
		Name:   "sugarcookie",
		IsDark: false,

		Primary:   sugarPrimary,
		Secondary: sugarSecondary,
		Tertiary:  sugarTertiary,
		Accent:    sugarAccent,

		// Backgrounds
		BgBase:        sugarBgBase,
		BgBaseLighter: sugarBgBaseLighter,
		BgSubtle:      sugarBgSubtle,
		BgOverlay:     sugarBgOverlay,

		// Foregrounds
		FgBase:      sugarFgBase,
		FgMuted:     sugarFgMuted,
		FgHalfMuted: sugarFgHalfMuted,
		FgSubtle:    sugarFgSubtle,
		FgSelected: sugarFgSelected,

		// Borders
		Border:      sugarBorder,
		BorderFocus: sugarBorderFocus,

		// Status
		Success: sugarSuccess,
		Error:   sugarError,
		Warning: sugarWarning,
		Info:    sugarInfo,

		// Colors
		White: sugarWhite,

		BlueLight: sugarBlueLight,
		BlueDark:  sugarBlueDark,
		Blue:      sugarBlue,

		Yellow: sugarYellow,
		Citron: sugarCitron,

		Green:      sugarGreen,
		GreenDark:  sugarGreenDark,
		GreenLight: sugarGreenLight,

		Red:      sugarRed,
		RedDark:  sugarRedDark,
		RedLight: sugarRedLight,
		Cherry:   sugarCherry,
	}

	// Text selection.
	t.TextSelection = lipgloss.NewStyle().Foreground(sugarFgSelected).Background(sugarPrimary)

	// LSP and MCP status.
	t.ItemOfflineIcon = lipgloss.NewStyle().Foreground(sugarFgSubtle).SetString("●")
	t.ItemBusyIcon = t.ItemOfflineIcon.Foreground(sugarCitron)
	t.ItemErrorIcon = t.ItemOfflineIcon.Foreground(sugarError)
	t.ItemOnlineIcon = t.ItemOfflineIcon.Foreground(sugarSuccess)

	// Editor: Yolo Mode - using a soft pink that fits the theme.
	t.YoloIconFocused = lipgloss.NewStyle().Foreground(sugarFgBase).Background(sugarAccent).Bold(true).SetString(" ! ")
	t.YoloIconBlurred = t.YoloIconFocused.Foreground(sugarFgMuted).Background(sugarBgOverlay)
	t.YoloDotsFocused = lipgloss.NewStyle().Foreground(sugarAccent).SetString(":::")
	t.YoloDotsBlurred = t.YoloDotsFocused.Foreground(sugarFgSubtle)

	// oAuth Chooser.
	t.AuthBorderSelected = lipgloss.NewStyle().BorderForeground(sugarSuccess)
	t.AuthTextSelected = lipgloss.NewStyle().Foreground(sugarSuccess)
	t.AuthBorderUnselected = lipgloss.NewStyle().BorderForeground(sugarBorder)
	t.AuthTextUnselected = lipgloss.NewStyle().Foreground(sugarFgMuted)

	return t
}
