package themes

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/tui/styles"
)

// Monokai theme - classic code editor color scheme.
// UI elements use brown/olive tones; bright colors are for code syntax only.
var (
	monokaiBg        = styles.ParseHex("#272822")
	monokaiBgLight   = styles.ParseHex("#31322c") // Slightly lighter
	monokaiBgSubtle  = styles.ParseHex("#3b3c36") // Subtle bg
	monokaiBgOverlay = styles.ParseHex("#454640") // Overlay

	// UI accent browns (for selection bars, buttons, borders)
	monokaiUISelect = styles.ParseHex("#75756a") // Selection bar - brighter olive brown
	monokaiUIButton = styles.ParseHex("#BEBEAC") // Button highlight
	monokaiUIBorder = styles.ParseHex("#6a6a5c") // Focused border - visible brown

	monokaiFg        = styles.ParseHex("#f8f8f2") // white
	monokaiComment   = styles.ParseHex("#75715e")
	monokaiFgHalf    = styles.ParseHex("#969282")

	// Syntax colors (for code highlighting only)
	monokaiRed       = styles.ParseHex("#f92672") // pink/red
	monokaiOrange    = styles.ParseHex("#fd971f")
	monokaiYellow    = styles.ParseHex("#e6db74")
	monokaiGreen     = styles.ParseHex("#a6e22e")
	monokaiCyan      = styles.ParseHex("#66d9ef")
	monokaiPurple    = styles.ParseHex("#ae81ff")
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

		// UI colors use brown tones, not syntax colors
		Primary:   monokaiUISelect, // Selection bar - lighter brown
		Secondary: monokaiUIButton, // Buttons - medium brown
		Tertiary:  monokaiComment,
		Accent:    monokaiUIBorder,

		BgBase:        monokaiBg,
		BgBaseLighter: monokaiBgLight,
		BgSubtle:      monokaiBgSubtle,
		BgOverlay:     monokaiBgOverlay,

		FgBase:      monokaiFg,
		FgMuted:     monokaiComment,
		FgHalfMuted: monokaiFgHalf,
		FgSubtle:    monokaiComment,
		FgSelected:  monokaiBg, // Dark text on brown selection bar
		FgButton:    monokaiBg, // Dark text on light beige buttons

		Border:      monokaiBgSubtle,
		BorderFocus: monokaiUIBorder, // Brown border, not cyan

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
