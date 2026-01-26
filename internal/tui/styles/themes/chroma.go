package themes

import (
	"image/color"
	"math"

	"charm.land/bubbles/v2/help"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/tui/styles"
	"github.com/lucasb-eyer/go-colorful"
)

// Chroma theme - RGB gaming aesthetic with vibrant rainbow colors.
// Colors are generated using HSV to create a full spectrum effect.
// This theme supports animation - accent colors cycle through the rainbow in sync.
// Base/text colors are neutral grayscale for readability and performance.

// hsvColor creates a color from HSV values with perceptual brightness compensation.
// h: hue (0-360), s: saturation (0-1), v: value (0-1)
// Uses styles.PerceivedBrightnessAdjust for consistent brightness across all chroma usage.
func hsvColor(h, s, v float64) colorful.Color {
	h = math.Mod(h+360, 360)
	// Adjust for perceptually darker hues (blue/purple range)
	// Reducing saturation makes them lighter/more pastel
	satReduce, valBoost := styles.PerceivedBrightnessAdjust(h)
	s = math.Max(0.0, s-satReduce)
	v = math.Min(1.0, v+valBoost)
	return colorful.Hsv(h, s, v)
}

// Static neutral colors (grayscale) - these don't animate for performance
var (
	chromaFgBase      = hsvColor(0, 0, 0.9)    // Near white
	chromaFgMuted     = hsvColor(0, 0, 0.6)    // Medium gray
	chromaFgHalfMuted = hsvColor(0, 0, 0.75)   // Light gray
	chromaFgSubtle    = hsvColor(0, 0, 0.45)   // Dark gray
	chromaFgSelected  = colorful.Color{R: 0, G: 0, B: 0} // Pure black for contrast on light selection
	chromaBgOverlay   = hsvColor(0, 0, 0.12)   // Dark gray overlay
	chromaBorder      = hsvColor(0, 0, 0.25)   // Gray border
	chromaSelection   = hsvColor(0, 0, 1.0)    // Bright white for selection
)

// chromaColorFunc generates animated rainbow colors based on hue offset.
// All animated colors use the SAME base hue offset so they move together in sync.
// Primary is kept neutral (white) for selection to avoid cache issues.
// Returns: Primary, Secondary, Tertiary, Accent, BorderFocus, Success, Error, Warning, Info,
// BgBase, BgBaseLighter, BgSubtle
func chromaColorFunc(baseHue, hueOffset float64) []color.Color {
	h := hueOffset // All colors shift by the same amount

	return []color.Color{
		hsvColor(h, 0.45, 1.0),     // 0: Primary - pastel animated for selection bar
		hsvColor(h+30, 0.45, 1.0),  // 1: Secondary - offset pastel for buttons
		chromaFgBase,               // 2: Tertiary - white (used for ">" prompt, cached)
		hsvColor(h, 1.0, 1.0),      // 3: Accent - animated rainbow
		hsvColor(h, 1.0, 1.0),      // 4: BorderFocus - animated rainbow
		chromaFgBase,               // 5: Success - white (used for checkmarks, cached)
		hsvColor(h+180, 0.9, 0.85), // 6: Error - opposite hue for contrast
		hsvColor(h, 1.0, 1.0),      // 7: Warning - animated
		hsvColor(h, 0.85, 0.9),     // 8: Info - animated
		hsvColor(h, 0.6, 0.10),     // 9: BgBase - more saturated tinted background
		hsvColor(h, 0.55, 0.12),    // 10: BgBaseLighter
		hsvColor(h, 0.5, 0.15),     // 11: BgSubtle
	}
}

// chromaStyleBuilder rebuilds lipgloss styles during animation.
func chromaStyleBuilder(t *styles.Theme, hueOffset float64) {
	h := hueOffset
	// Wrap colors in Color8Bit to ensure proper RGBA conversion for lipgloss
	wrap := styles.NewColor8Bit
	accent := wrap(hsvColor(h, 1.0, 1.0))
	bgDark := wrap(hsvColor(h, 0.6, 0.10))

	// Text colors stay neutral gray - they get cached and can't animate reliably
	// Only accents, borders, backgrounds, and specific UI elements animate

	// Text selection uses neutral selection color (not animated for caching)
	t.TextSelection = lipgloss.NewStyle().Foreground(wrap(chromaFgSelected)).Background(wrap(chromaSelection))

	// LSP and MCP status - all use animated accent
	t.ItemOfflineIcon = lipgloss.NewStyle().Foreground(wrap(chromaFgSubtle)).SetString("●")
	t.ItemBusyIcon = t.ItemOfflineIcon.Foreground(accent)
	t.ItemErrorIcon = t.ItemOfflineIcon.Foreground(wrap(hsvColor(h+180, 0.9, 0.85)))
	t.ItemOnlineIcon = t.ItemOfflineIcon.Foreground(accent)

	// Editor: Yolo Mode - animated accent color
	t.YoloIconFocused = lipgloss.NewStyle().Foreground(bgDark).Background(accent).Bold(true).SetString(" ! ")
	t.YoloIconBlurred = t.YoloIconFocused.Foreground(wrap(chromaFgMuted)).Background(wrap(chromaBgOverlay))
	t.YoloDotsFocused = lipgloss.NewStyle().Foreground(accent).SetString(":::")
	t.YoloDotsBlurred = t.YoloDotsFocused.Foreground(wrap(chromaFgSubtle))

	// oAuth Chooser
	t.AuthBorderSelected = lipgloss.NewStyle().BorderForeground(accent)
	t.AuthTextSelected = lipgloss.NewStyle().Foreground(accent)
	t.AuthBorderUnselected = lipgloss.NewStyle().BorderForeground(wrap(chromaBorder))
	t.AuthTextUnselected = lipgloss.NewStyle().Foreground(wrap(chromaFgMuted))

	// Dialog styles - animated rainbow text for dialogs (they render fresh each frame)
	pastel := wrap(hsvColor(h, 0.45, 1.0)) // Match Secondary
	t.DialogTitle = lipgloss.NewStyle().Foreground(accent).Bold(true)
	t.DialogText = lipgloss.NewStyle().Foreground(pastel)
	t.DialogHint = lipgloss.NewStyle().Foreground(wrap(chromaFgMuted)).Faint(true)

	// Help bar styles - animated shortcuts, neutral dots
	base := lipgloss.NewStyle()
	t.HelpStyles = help.Styles{
		ShortKey:       base.Foreground(pastel),
		ShortDesc:      base.Foreground(wrap(chromaFgMuted)),
		ShortSeparator: base.Foreground(wrap(chromaFgSubtle)),
		Ellipsis:       base.Foreground(wrap(chromaFgSubtle)),
		FullKey:        base.Foreground(pastel),
		FullDesc:       base.Foreground(wrap(chromaFgMuted)),
		FullSeparator:  base.Foreground(wrap(chromaFgSubtle)),
	}
}

// NewChromaTheme creates an RGB gaming aesthetic theme with rainbow colors.
// Accent colors animate through the rainbow; base colors are neutral for performance.
func NewChromaTheme() *styles.Theme {
	// Get initial colors at hue offset 0
	initialColors := chromaColorFunc(0, 0)

	// Wrap colors in Color8Bit to ensure proper RGBA conversion for lipgloss
	wrap := styles.NewColor8Bit

	t := &styles.Theme{
		Name:   "chroma",
		IsDark: true,

		// Animation settings
		Animated:       true,
		AnimationSpeed: 30, // 30 degrees per second = full cycle in 12 seconds
		ColorFunc:      chromaColorFunc,
		StyleBuilder:   chromaStyleBuilder,

		// Animated accent colors
		Primary:   wrap(initialColors[0]),
		Secondary: wrap(initialColors[1]),
		Tertiary:  wrap(initialColors[2]),
		Accent:    wrap(initialColors[3]),

		// Animated backgrounds (subtle tint)
		BgBase:        wrap(initialColors[9]),
		BgBaseLighter: wrap(initialColors[10]),
		BgSubtle:      wrap(initialColors[11]),
		BgOverlay:     wrap(chromaBgOverlay),

		// Static neutral foregrounds (grayscale) - text gets cached by components
		FgBase:      wrap(chromaFgBase),      // Near white
		FgMuted:     wrap(chromaFgMuted),     // Medium gray
		FgHalfMuted: wrap(chromaFgHalfMuted), // Light gray
		FgSubtle:    wrap(chromaFgSubtle),    // Dark gray
		FgSelected:  wrap(chromaFgSelected),
		FgButton:    wrap(chromaFgSelected),
		FgCursor:    wrap(chromaFgBase), // Neutral white cursor (not animated)

		// Borders
		Border:      wrap(chromaBorder),
		BorderFocus: wrap(initialColors[4]),

		// Animated status colors
		Success: wrap(initialColors[5]),
		Error:   wrap(initialColors[6]),
		Warning: wrap(initialColors[7]),
		Info:    wrap(initialColors[8]),

		// Named colors - static neutral values
		White:      wrap(chromaFgBase),
		Yellow:     wrap(initialColors[7]),
		Citron:     wrap(initialColors[7]),
		Blue:       wrap(initialColors[0]),
		BlueLight:  wrap(initialColors[0]),
		BlueDark:   wrap(initialColors[1]),
		Green:      wrap(initialColors[5]),
		GreenDark:  wrap(initialColors[2]),
		GreenLight: wrap(initialColors[1]),
		Red:        wrap(initialColors[6]),
		RedDark:    wrap(hsvColor(180, 0.7, 0.6)),
		RedLight:   wrap(hsvColor(180, 0.9, 0.9)),
		Cherry:     wrap(initialColors[6]),
	}

	// Initialize lipgloss styles with initial colors
	chromaStyleBuilder(t, 0)

	// Initialize the atomic color snapshot so GetAnimatedColors never returns nil
	t.InitAnimatedColors()

	return t
}
