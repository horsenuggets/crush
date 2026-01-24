package themes

import (
	"image/color"
	"math"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/tui/styles"
	"github.com/lucasb-eyer/go-colorful"
)

// Chroma theme - RGB gaming aesthetic with vibrant rainbow colors.
// Colors are generated using HSV to create a full spectrum effect.
// This theme supports animation - colors cycle through the rainbow over time.
var (
	// Static rainbow colors for style initialization (these are used for lipgloss styles)
	chromaRed     = hsvColor(0, 1.0, 1.0)   // Red
	chromaYellow  = hsvColor(60, 1.0, 1.0)  // Yellow
	chromaGreen   = hsvColor(120, 1.0, 1.0) // Green
	chromaCyan    = hsvColor(180, 1.0, 1.0) // Cyan
	chromaMagenta = hsvColor(300, 1.0, 1.0) // Magenta

	// Dark backgrounds with subtle color tints
	chromaBgBase    = hsvColor(270, 0.3, 0.08)  // Very dark purple-tinted
	chromaBgOverlay = hsvColor(270, 0.15, 0.18) // Overlay (static)

	// Foreground colors - bright and vibrant (static for readability)
	chromaFgBase      = hsvColor(180, 0.1, 0.95)   // Near-white with cyan tint
	chromaFgMuted     = hsvColor(270, 0.2, 0.6)    // Muted purple
	chromaFgHalfMuted = hsvColor(210, 0.15, 0.75)  // Light blue-gray
	chromaFgSubtle    = hsvColor(270, 0.25, 0.45)  // Subtle purple
	chromaFgSelected  = styles.ParseHex("#ffffff") // Pure white

	// Border (static)
	chromaBorder = hsvColor(270, 0.4, 0.25) // Purple border
)

// hsvColor creates a color from HSV values.
// h: hue (0-360), s: saturation (0-1), v: value (0-1)
func hsvColor(h, s, v float64) colorful.Color {
	return colorful.Hsv(h, s, v)
}

// chromaColorFunc generates animated rainbow colors based on hue offset.
// Returns colors in this order: Primary, Secondary, Tertiary, Accent, BorderFocus,
// Success, Error, Warning, Info, BgBase, BgBaseLighter, BgSubtle
func chromaColorFunc(baseHue, hueOffset float64) []color.Color {
	// Helper to create color at offset hue
	atHue := func(h float64) color.Color {
		return hsvColor(math.Mod(h+hueOffset, 360), 1.0, 1.0)
	}
	atHueMuted := func(h float64) color.Color {
		return hsvColor(math.Mod(h+hueOffset, 360), 0.7, 0.7)
	}
	atHueBg := func(h, s, v float64) color.Color {
		return hsvColor(math.Mod(h+hueOffset, 360), s, v)
	}

	return []color.Color{
		atHue(180),            // Primary (cyan base)
		atHue(270),            // Secondary (violet base)
		atHue(210),            // Tertiary (blue base)
		atHue(300),            // Accent (magenta base)
		atHue(180),            // BorderFocus (cyan base)
		atHue(120),            // Success (green base)
		atHueMuted(0),         // Error (red base, slightly muted so it's still visible)
		atHue(60),             // Warning (yellow base)
		atHue(210),            // Info (blue base)
		atHueBg(270, 0.3, 0.08),  // BgBase
		atHueBg(270, 0.25, 0.10), // BgBaseLighter
		atHueBg(270, 0.2, 0.12),  // BgSubtle
	}
}

// NewChromaTheme creates an RGB gaming aesthetic theme with rainbow colors.
// This theme supports animation - colors cycle through the rainbow over time.
func NewChromaTheme() *styles.Theme {
	// Get initial colors at hue offset 0
	initialColors := chromaColorFunc(0, 0)

	t := &styles.Theme{
		Name:   "chroma",
		IsDark: true,

		// Animation settings
		Animated:       true,
		AnimationSpeed: 30, // 30 degrees per second = full cycle in 12 seconds
		ColorFunc:      chromaColorFunc,

		// Use different rainbow colors for each role (initial values)
		Primary:   initialColors[0],
		Secondary: initialColors[1],
		Tertiary:  initialColors[2],
		Accent:    initialColors[3],

		// Backgrounds (initial values, will animate)
		BgBase:        initialColors[9],
		BgBaseLighter: initialColors[10],
		BgSubtle:      initialColors[11],
		BgOverlay:     chromaBgOverlay,

		// Foregrounds (static for readability)
		FgBase:      chromaFgBase,
		FgMuted:     chromaFgMuted,
		FgHalfMuted: chromaFgHalfMuted,
		FgSubtle:    chromaFgSubtle,
		FgSelected:  chromaFgSelected,

		// Borders
		Border:      chromaBorder,
		BorderFocus: initialColors[4],

		// Status colors - each a different rainbow color (initial values)
		Success: initialColors[5],
		Error:   initialColors[6],
		Warning: initialColors[7],
		Info:    initialColors[8],

		// Colors - full rainbow spectrum (static fallbacks)
		White: chromaFgBase,

		BlueLight: hsvColor(180, 1.0, 1.0), // Cyan
		BlueDark:  hsvColor(240, 1.0, 1.0), // Indigo
		Blue:      hsvColor(210, 1.0, 1.0), // Blue

		Yellow: hsvColor(60, 1.0, 1.0),  // Yellow
		Citron: hsvColor(90, 1.0, 1.0),  // Lime

		Green:      hsvColor(120, 1.0, 1.0), // Green
		GreenDark:  hsvColor(120, 0.7, 0.7), // Muted green
		GreenLight: hsvColor(90, 1.0, 1.0),  // Lime

		Red:      hsvColor(0, 1.0, 1.0),   // Red
		RedDark:  hsvColor(0, 0.7, 0.7),   // Muted red
		RedLight: hsvColor(330, 1.0, 1.0), // Pink
		Cherry:   hsvColor(300, 1.0, 1.0), // Magenta
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
