package themes

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/tui/styles"
)

// Monokai theme - classic code editor color scheme.
var (
	monokaiBg        = styles.ParseHex("#272822")
	monokaiBgLight   = styles.ParseHex("#31322c") // Slightly lighter
	monokaiBgSubtle  = styles.ParseHex("#3b3c36") // Subtle bg
	monokaiBgOverlay = styles.ParseHex("#454640") // Overlay

	monokaiFg        = styles.ParseHex("#f8f8f2") // white
	monokaiComment   = styles.ParseHex("#75715e")
	monokaiRed       = styles.ParseHex("#f92672") // pink/red
	monokaiOrange    = styles.ParseHex("#fd971f")
	monokaiYellow    = styles.ParseHex("#e6db74")
	monokaiGreen     = styles.ParseHex("#a6e22e")
	monokaiCyan      = styles.ParseHex("#66d9ef")
	monokaiPurple    = styles.ParseHex("#ae81ff")
	monokaiFgHalf    = styles.ParseHex("#969282")
	monokaiBlueDark  = styles.ParseHex("#3c8ca0")
	monokaiGreenDark = styles.ParseHex("#78b41e")
	monokaiRedDark   = styles.ParseHex("#b41e50")
	monokaiRedLight  = styles.ParseHex("#ff5a96")
)

// NewMonokaiTheme creates the classic Monokai editor theme.
func NewMonokaiTheme() *styles.Theme {
	t := &styles.Theme{
		Name:   "monokai",
		IsDark: true,

		Primary:   monokaiCyan,
		Secondary: monokaiComment,
		Tertiary:  monokaiComment,
		Accent:    monokaiPurple,

		BgBase:        monokaiBg,
		BgBaseLighter: monokaiBgLight,
		BgSubtle:      monokaiBgSubtle,
		BgOverlay:     monokaiBgOverlay,

		FgBase:      monokaiFg,
		FgMuted:     monokaiComment,
		FgHalfMuted: monokaiFgHalf,
		FgSubtle:    monokaiComment,
		FgSelected:  monokaiFg,

		Border:      monokaiBgSubtle,
		BorderFocus: monokaiCyan,

		Success: monokaiGreen,
		Error:   monokaiRed,
		Warning: monokaiOrange,
		Info:    monokaiCyan,

		White:     monokaiFg,
		BlueLight: monokaiCyan,
		BlueDark:  monokaiBlueDark,
		Blue:      monokaiCyan,

		Yellow: monokaiYellow,
		Citron: monokaiOrange,

		Green:      monokaiGreen,
		GreenDark:  monokaiGreenDark,
		GreenLight: monokaiGreen,

		Red:      monokaiRed,
		RedDark:  monokaiRedDark,
		RedLight: monokaiRedLight,
		Cherry:   monokaiRed,
	}

	t.TextSelection = lipgloss.NewStyle().Foreground(monokaiFg).Background(monokaiBgOverlay)

	t.ItemOfflineIcon = lipgloss.NewStyle().Foreground(monokaiComment).SetString("●")
	t.ItemBusyIcon = t.ItemOfflineIcon.Foreground(monokaiOrange)
	t.ItemErrorIcon = t.ItemOfflineIcon.Foreground(monokaiRed)
	t.ItemOnlineIcon = t.ItemOfflineIcon.Foreground(monokaiGreen)

	t.YoloIconFocused = lipgloss.NewStyle().Foreground(monokaiBg).Background(monokaiOrange).Bold(true).SetString(" ! ")
	t.YoloIconBlurred = t.YoloIconFocused.Foreground(monokaiFg).Background(monokaiBgOverlay)
	t.YoloDotsFocused = lipgloss.NewStyle().Foreground(monokaiOrange).SetString(":::")
	t.YoloDotsBlurred = t.YoloDotsFocused.Foreground(monokaiComment)

	t.AuthBorderSelected = lipgloss.NewStyle().BorderForeground(monokaiGreen)
	t.AuthTextSelected = lipgloss.NewStyle().Foreground(monokaiGreen)
	t.AuthBorderUnselected = lipgloss.NewStyle().BorderForeground(monokaiBgSubtle)
	t.AuthTextUnselected = lipgloss.NewStyle().Foreground(monokaiComment)

	return t
}
