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
// This theme supports animation - accent colors cycle through the rainbow in sync.
// Base/text colors are neutral grayscale for readability and performance.

// hsvColor creates a color from HSV values.
// h: hue (0-360), s: saturation (0-1), v: value (0-1)
func hsvColor(h, s, v float64) colorful.Color {
	return colorful.Hsv(math.Mod(h+360, 360), s, v)
}

// Static neutral colors (grayscale) - these don't animate for performance
var (
	chromaFgBase      = hsvColor(0, 0, 0.9)    // Near white
	chromaFgMuted     = hsvColor(0, 0, 0.6)    // Medium gray
	chromaFgHalfMuted = hsvColor(0, 0, 0.75)   // Light gray
	chromaFgSubtle    = hsvColor(0, 0, 0.45)   // Dark gray
	chromaFgSelected  = hsvColor(0, 0, 0.0)    // Black for contrast on light selection
	chromaBgOverlay   = hsvColor(0, 0, 0.12)   // Dark gray overlay
	chromaBorder      = hsvColor(0, 0, 0.25)   // Gray border
	chromaSelection   = hsvColor(0, 0, 0.35)   // Semi-transparent white for selection
)

// chromaColorFunc generates animated rainbow colors based on hue offset.
// All animated colors use the SAME base hue offset so they move together in sync.
// Primary is kept neutral (gray) for selection to avoid cache issues.
// Returns: Primary, Secondary, Tertiary, Accent, BorderFocus, Success, Error, Warning, Info,
// BgBase, BgBaseLighter, BgSubtle
func chromaColorFunc(baseHue, hueOffset float64) []color.Color {
	h := hueOffset // All colors shift by the same amount

	return []color.Color{
		chromaSelection,            // 0: Primary - neutral gray selection (not animated)
		hsvColor(h, 0.8, 0.9),      // 1: Secondary - animated
		hsvColor(h, 0.7, 0.85),     // 2: Tertiary - animated
		hsvColor(h, 1.0, 1.0),      // 3: Accent - animated rainbow
		hsvColor(h, 1.0, 1.0),      // 4: BorderFocus - animated rainbow
		hsvColor(h, 0.9, 0.9),      // 5: Success - animated (checkmarks)
		hsvColor(h+180, 0.9, 0.85), // 6: Error - opposite hue for contrast
		hsvColor(h, 1.0, 1.0),      // 7: Warning - animated
		hsvColor(h, 0.85, 0.9),     // 8: Info - animated
		hsvColor(h, 0.4, 0.06),     // 9: BgBase - subtle tinted background
		hsvColor(h, 0.35, 0.08),    // 10: BgBaseLighter
		hsvColor(h, 0.3, 0.11),     // 11: BgSubtle
	}
}

// chromaStyleBuilder rebuilds lipgloss styles during animation.
func chromaStyleBuilder(t *styles.Theme, hueOffset float64) {
	h := hueOffset
	accent := hsvColor(h, 1.0, 1.0)
	bgDark := hsvColor(h, 0.4, 0.06)

	// Text selection uses neutral selection color (not animated for caching)
	t.TextSelection = lipgloss.NewStyle().Foreground(chromaFgSelected).Background(chromaSelection)

	// LSP and MCP status - all use animated accent
	t.ItemOfflineIcon = lipgloss.NewStyle().Foreground(chromaFgSubtle).SetString("●")
	t.ItemBusyIcon = t.ItemOfflineIcon.Foreground(accent)
	t.ItemErrorIcon = t.ItemOfflineIcon.Foreground(hsvColor(h+180, 0.9, 0.85))
	t.ItemOnlineIcon = t.ItemOfflineIcon.Foreground(accent)

	// Editor: Yolo Mode - animated accent color
	t.YoloIconFocused = lipgloss.NewStyle().Foreground(bgDark).Background(accent).Bold(true).SetString(" ! ")
	t.YoloIconBlurred = t.YoloIconFocused.Foreground(chromaFgMuted).Background(chromaBgOverlay)
	t.YoloDotsFocused = lipgloss.NewStyle().Foreground(accent).SetString(":::")
	t.YoloDotsBlurred = t.YoloDotsFocused.Foreground(chromaFgSubtle)

	// oAuth Chooser
	t.AuthBorderSelected = lipgloss.NewStyle().BorderForeground(accent)
	t.AuthTextSelected = lipgloss.NewStyle().Foreground(accent)
	t.AuthBorderUnselected = lipgloss.NewStyle().BorderForeground(chromaBorder)
	t.AuthTextUnselected = lipgloss.NewStyle().Foreground(chromaFgMuted)
}

// NewChromaTheme creates an RGB gaming aesthetic theme with rainbow colors.
// Accent colors animate through the rainbow; base colors are neutral for performance.
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
		StyleBuilder:   chromaStyleBuilder,

		// Animated accent colors
		Primary:   initialColors[0],
		Secondary: initialColors[1],
		Tertiary:  initialColors[2],
		Accent:    initialColors[3],

		// Animated backgrounds (subtle tint)
		BgBase:        initialColors[9],
		BgBaseLighter: initialColors[10],
		BgSubtle:      initialColors[11],
		BgOverlay:     chromaBgOverlay,

		// Static neutral foregrounds for readability
		FgBase:      chromaFgBase,
		FgMuted:     chromaFgMuted,
		FgHalfMuted: chromaFgHalfMuted,
		FgSubtle:    chromaFgSubtle,
		FgSelected:  chromaFgSelected,

		// Borders
		Border:      chromaBorder,
		BorderFocus: initialColors[4],

		// Animated status colors
		Success: initialColors[5],
		Error:   initialColors[6],
		Warning: initialColors[7],
		Info:    initialColors[8],

		// Named colors - static neutral values
		White:      chromaFgBase,
		Yellow:     initialColors[7],
		Citron:     initialColors[7],
		Blue:       initialColors[0],
		BlueLight:  initialColors[0],
		BlueDark:   initialColors[1],
		Green:      initialColors[5],
		GreenDark:  initialColors[2],
		GreenLight: initialColors[1],
		Red:        initialColors[6],
		RedDark:    hsvColor(180, 0.7, 0.6),
		RedLight:   hsvColor(180, 0.9, 0.9),
		Cherry:     initialColors[6],
	}

	// Initialize lipgloss styles with initial colors
	chromaStyleBuilder(t, 0)

	return t
}
