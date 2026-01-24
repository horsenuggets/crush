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

// perceivedBrightnessAdjust returns saturation reduction and value boost for
// perceptually darker hues. Blue (240°) and purple (270-300°) appear darker
// than yellow/green/cyan at the same HSV values, so we compensate by reducing
// saturation (makes colors lighter/more pastel) and boosting value.
func perceivedBrightnessAdjust(h float64) (saturationReduce, valueBoost float64) {
	h = math.Mod(h+360, 360)

	// Center the adjustment around 250° (deep blue)
	// This is where colors appear darkest perceptually
	center := 250.0
	// Half-width of the adjustment region in degrees
	width := 70.0

	// Calculate distance from center (handling wrap-around)
	dist := math.Abs(h - center)
	if dist > 180 {
		dist = 360 - dist
	}

	// Outside adjustment region - no changes needed
	if dist > width {
		return 0, 0
	}

	// Smooth cosine curve for gradual transition
	factor := (1 + math.Cos(math.Pi*dist/width)) / 2

	// Lean more heavily on saturation reduction for lighter appearance
	// Also add stronger value boost for the deep blue range
	saturationReduce = 0.20 * factor // reduce saturation by up to 20%
	valueBoost = 0.10 * factor       // boost value by up to 10%

	return saturationReduce, valueBoost
}

// hsvColor creates a color from HSV values with perceptual brightness compensation.
// h: hue (0-360), s: saturation (0-1), v: value (0-1)
func hsvColor(h, s, v float64) colorful.Color {
	h = math.Mod(h+360, 360)
	// Adjust for perceptually darker hues (blue/purple range)
	// Reducing saturation makes them lighter/more pastel
	satReduce, valBoost := perceivedBrightnessAdjust(h)
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
	accent := hsvColor(h, 1.0, 1.0)
	bgDark := hsvColor(h, 0.6, 0.10)

	// Text colors stay neutral gray - they get cached and can't animate reliably
	// Only accents, borders, backgrounds, and specific UI elements animate

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

	// Dialog styles - animated rainbow text for dialogs (they render fresh each frame)
	pastel := hsvColor(h, 0.45, 1.0) // Match Secondary
	t.DialogTitle = lipgloss.NewStyle().Foreground(accent).Bold(true)
	t.DialogText = lipgloss.NewStyle().Foreground(pastel)
	t.DialogHint = lipgloss.NewStyle().Foreground(chromaFgMuted).Faint(true)

	// Help bar styles - animated shortcuts, neutral dots
	base := lipgloss.NewStyle()
	t.HelpStyles = help.Styles{
		ShortKey:       base.Foreground(pastel),
		ShortDesc:      base.Foreground(chromaFgMuted),
		ShortSeparator: base.Foreground(chromaFgSubtle),
		Ellipsis:       base.Foreground(chromaFgSubtle),
		FullKey:        base.Foreground(pastel),
		FullDesc:       base.Foreground(chromaFgMuted),
		FullSeparator:  base.Foreground(chromaFgSubtle),
	}
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

		// Static neutral foregrounds (grayscale) - text gets cached by components
		FgBase:      chromaFgBase,      // Near white
		FgMuted:     chromaFgMuted,     // Medium gray
		FgHalfMuted: chromaFgHalfMuted, // Light gray
		FgSubtle:    chromaFgSubtle,    // Dark gray
		FgSelected:  chromaFgSelected,
		FgCursor:    chromaFgBase,      // Neutral white cursor (not animated)

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
