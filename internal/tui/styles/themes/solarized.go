package themes

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/tui/styles"
)

// Solarized Dark theme - classic Ethan Schoonover color scheme.
var (
	// Base colors
	solBase03 = styles.ParseHex("#002b36") // background
	solBase02 = styles.ParseHex("#073642") // bg highlights
	solBase01 = styles.ParseHex("#586e75") // comments
	solBase00 = styles.ParseHex("#657b83") // secondary
	solBase0  = styles.ParseHex("#839496") // body text
	solBase1  = styles.ParseHex("#93a1a1") // emphasis
	solBase2  = styles.ParseHex("#eee8d5") // bg (light)
	solBase3  = styles.ParseHex("#fdf6e3") // bg (light)

	// Accent colors
	solYellow  = styles.ParseHex("#b58900")
	solOrange  = styles.ParseHex("#cb4b16")
	solRed     = styles.ParseHex("#dc322f")
	solMagenta = styles.ParseHex("#d33682")
	solViolet  = styles.ParseHex("#6c71c4")
	solBlue    = styles.ParseHex("#268bd2")
	solCyan    = styles.ParseHex("#2aa198")
	solGreen   = styles.ParseHex("#859900")
)

// NewSolarizedTheme creates the classic Solarized Dark theme.
func NewSolarizedTheme() *styles.Theme {
	t := &styles.Theme{
		Name:   "solarized",
		IsDark: true,

		Primary:   solBlue,
		Secondary: solBase0,
		Tertiary:  solBase01,
		Accent:    solCyan,

		BgBase:        solBase03,
		BgBaseLighter: solBase02,
		BgSubtle:      solBase02,
		BgOverlay:     solBase01,

		FgBase:      solBase0,
		FgMuted:     solBase01,
		FgHalfMuted: solBase00,
		FgSubtle:    solBase01,
		FgSelected:  solBase3,
		FgButton:    solBase3,

		Border:      solBase02,
		BorderFocus: solBlue,

		Success: solGreen,
		Error:   solRed,
		Warning: solYellow,
		Info:    solBlue,

		White:     solBase1,
		BlueLight: solCyan,
		BlueDark:  solViolet,
		Blue:      solBlue,

		Yellow: solYellow,
		Citron: solOrange,

		Green:      solGreen,
		GreenDark:  solCyan,
		GreenLight: solGreen,

		Red:      solRed,
		RedDark:  solOrange,
		RedLight: solMagenta,
		Cherry:   solMagenta,
	}

	t.TextSelection = lipgloss.NewStyle().Foreground(solBase03).Background(solBlue)

	t.ItemOfflineIcon = lipgloss.NewStyle().Foreground(solBase01).SetString("●")
	t.ItemBusyIcon = t.ItemOfflineIcon.Foreground(solYellow)
	t.ItemErrorIcon = t.ItemOfflineIcon.Foreground(solRed)
	t.ItemOnlineIcon = t.ItemOfflineIcon.Foreground(solGreen)

	t.YoloIconFocused = lipgloss.NewStyle().Foreground(solBase03).Background(solYellow).Bold(true).SetString(" ! ")
	t.YoloIconBlurred = t.YoloIconFocused.Foreground(solBase0).Background(solBase02)
	t.YoloDotsFocused = lipgloss.NewStyle().Foreground(solYellow).SetString(":::")
	t.YoloDotsBlurred = t.YoloDotsFocused.Foreground(solBase01)

	t.AuthBorderSelected = lipgloss.NewStyle().BorderForeground(solGreen)
	t.AuthTextSelected = lipgloss.NewStyle().Foreground(solGreen)
	t.AuthBorderUnselected = lipgloss.NewStyle().BorderForeground(solBase02)
	t.AuthTextUnselected = lipgloss.NewStyle().Foreground(solBase01)

	return t
}
