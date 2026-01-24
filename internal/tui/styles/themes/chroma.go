package themes

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/tui/styles"
	"github.com/lucasb-eyer/go-colorful"
)

// Chroma theme - RGB gaming aesthetic with vibrant rainbow colors.
// Colors are generated using HSV to create a full spectrum effect.
var (
	// Generate rainbow colors at different hue positions (0-360)
	chromaRed     = hsvColor(0, 1.0, 1.0)    // Red
	chromaOrange  = hsvColor(30, 1.0, 1.0)   // Orange
	chromaYellow  = hsvColor(60, 1.0, 1.0)   // Yellow
	chromaLime    = hsvColor(90, 1.0, 1.0)   // Lime
	chromaGreen   = hsvColor(120, 1.0, 1.0)  // Green
	chromaCyan    = hsvColor(180, 1.0, 1.0)  // Cyan
	chromaBlue    = hsvColor(210, 1.0, 1.0)  // Blue
	chromaIndigo  = hsvColor(240, 1.0, 1.0)  // Indigo
	chromaViolet  = hsvColor(270, 1.0, 1.0)  // Violet
	chromaMagenta = hsvColor(300, 1.0, 1.0)  // Magenta
	chromaPink    = hsvColor(330, 1.0, 1.0)  // Pink

	// Muted rainbow variants (lower saturation/value)
	chromaRedMuted    = hsvColor(0, 0.7, 0.7)
	chromaGreenMuted  = hsvColor(120, 0.7, 0.7)
	chromaBlueMuted   = hsvColor(210, 0.7, 0.7)
	chromaYellowMuted = hsvColor(60, 0.8, 0.9)
	chromaCyanMuted   = hsvColor(180, 0.7, 0.8)

	// Dark backgrounds with subtle color tints
	chromaBgBase        = hsvColor(270, 0.3, 0.08)  // Very dark purple-tinted
	chromaBgBaseLighter = hsvColor(270, 0.25, 0.10) // Slightly lighter
	chromaBgSubtle      = hsvColor(270, 0.2, 0.12)  // Subtle overlay
	chromaBgOverlay     = hsvColor(270, 0.15, 0.18) // Overlay

	// Foreground colors - bright and vibrant
	chromaFgBase      = hsvColor(180, 0.1, 0.95)  // Near-white with cyan tint
	chromaFgMuted     = hsvColor(270, 0.2, 0.6)   // Muted purple
	chromaFgHalfMuted = hsvColor(210, 0.15, 0.75) // Light blue-gray
	chromaFgSubtle    = hsvColor(270, 0.25, 0.45) // Subtle purple
	chromaFgSelected  = styles.ParseHex("#ffffff") // Pure white

	// Borders with color cycling effect
	chromaBorder      = hsvColor(270, 0.4, 0.25) // Purple border
	chromaBorderFocus = chromaCyan               // Bright cyan focus
)

// hsvColor creates a color from HSV values.
// h: hue (0-360), s: saturation (0-1), v: value (0-1)
func hsvColor(h, s, v float64) colorful.Color {
	return colorful.Hsv(h, s, v)
}

// NewChromaTheme creates an RGB gaming aesthetic theme with rainbow colors.
func NewChromaTheme() *styles.Theme {
	t := &styles.Theme{
		Name:   "chroma",
		IsDark: true,

		// Use different rainbow colors for each role
		Primary:   chromaCyan,    // Cyan for primary
		Secondary: chromaViolet,  // Violet for secondary
		Tertiary:  chromaBlue,    // Blue for tertiary
		Accent:    chromaMagenta, // Magenta for accent

		// Backgrounds
		BgBase:        chromaBgBase,
		BgBaseLighter: chromaBgBaseLighter,
		BgSubtle:      chromaBgSubtle,
		BgOverlay:     chromaBgOverlay,

		// Foregrounds
		FgBase:      chromaFgBase,
		FgMuted:     chromaFgMuted,
		FgHalfMuted: chromaFgHalfMuted,
		FgSubtle:    chromaFgSubtle,
		FgSelected:  chromaFgSelected,

		// Borders
		Border:      chromaBorder,
		BorderFocus: chromaBorderFocus,

		// Status colors - each a different rainbow color
		Success: chromaGreen,  // Green for success
		Error:   chromaRed,    // Red for error
		Warning: chromaYellow, // Yellow for warning
		Info:    chromaBlue,   // Blue for info

		// Colors - full rainbow spectrum
		White: chromaFgBase,

		BlueLight: chromaCyan,
		BlueDark:  chromaIndigo,
		Blue:      chromaBlue,

		Yellow: chromaYellow,
		Citron: chromaLime,

		Green:      chromaGreen,
		GreenDark:  chromaGreenMuted,
		GreenLight: chromaLime,

		Red:      chromaRed,
		RedDark:  chromaRedMuted,
		RedLight: chromaPink,
		Cherry:   chromaMagenta,
	}

	// Text selection with rainbow effect
	t.TextSelection = lipgloss.NewStyle().Foreground(chromaBgBase).Background(chromaCyan)

	// LSP and MCP status - each a different rainbow color
	t.ItemOfflineIcon = lipgloss.NewStyle().Foreground(chromaFgSubtle).SetString("●")
	t.ItemBusyIcon = t.ItemOfflineIcon.Foreground(chromaYellow)
	t.ItemErrorIcon = t.ItemOfflineIcon.Foreground(chromaRed)
	t.ItemOnlineIcon = t.ItemOfflineIcon.Foreground(chromaGreen)

	// Editor: Yolo Mode - vibrant magenta/pink for danger
	t.YoloIconFocused = lipgloss.NewStyle().Foreground(chromaBgBase).Background(chromaMagenta).Bold(true).SetString(" ! ")
	t.YoloIconBlurred = t.YoloIconFocused.Foreground(chromaFgBase).Background(chromaBgOverlay)
	t.YoloDotsFocused = lipgloss.NewStyle().Foreground(chromaMagenta).SetString(":::")
	t.YoloDotsBlurred = t.YoloDotsFocused.Foreground(chromaFgSubtle)

	// oAuth Chooser
	t.AuthBorderSelected = lipgloss.NewStyle().BorderForeground(chromaGreen)
	t.AuthTextSelected = lipgloss.NewStyle().Foreground(chromaGreen)
	t.AuthBorderUnselected = lipgloss.NewStyle().BorderForeground(chromaBorder)
	t.AuthTextUnselected = lipgloss.NewStyle().Foreground(chromaFgMuted)

	return t
}
