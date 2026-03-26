package styles

import (
	"fmt"
	"image/color"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"charm.land/bubbles/v2/filepicker"
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2/ansi"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/shared/colors"
	"github.com/charmbracelet/crush/internal/tui/exp/diffview"
	"github.com/charmbracelet/x/exp/charmtone"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/rivo/uniseg"
)

const (
	defaultListIndent      = 2
	defaultListLevelIndent = 4
	defaultMargin          = 2
)

type Theme struct {
	Name   string
	IsDark bool

	// Animation support for themes like chroma
	Animated       bool                                            // Whether this theme has animated colors
	AnimationSpeed float64                                         // Hue degrees per second (e.g., 30 = full cycle in 12 seconds)
	ColorFunc      func(baseHue, hueOffset float64) []color.Color  // Function to generate colors based on hue offset
	StyleBuilder   func(t *Theme, hueOffset float64)               // Function to rebuild lipgloss styles during animation

	Primary   color.Color
	Secondary color.Color
	Tertiary  color.Color
	Accent    color.Color

	BgBase        color.Color
	BgBaseLighter color.Color
	BgSubtle      color.Color
	BgOverlay     color.Color

	FgBase      color.Color
	FgMuted     color.Color
	FgHalfMuted color.Color
	FgSubtle    color.Color
	FgSelected color.Color // Text color on selection bars and buttons (Primary background)
	FgCursor   color.Color // Optional cursor color, defaults to Secondary if nil

	Border      color.Color
	BorderFocus color.Color

	Success color.Color
	Error   color.Color
	Warning color.Color
	Info    color.Color

	// Colors
	// White
	White color.Color

	// Blues
	BlueLight color.Color
	BlueDark  color.Color
	Blue      color.Color

	// Yellows
	Yellow color.Color
	Citron color.Color

	// Greens
	Green      color.Color
	GreenDark  color.Color
	GreenLight color.Color

	// Reds
	Red      color.Color
	RedDark  color.Color
	RedLight color.Color
	Cherry   color.Color

	// Text selection.
	TextSelection lipgloss.Style

	// LSP and MCP status indicators.
	ItemOfflineIcon lipgloss.Style
	ItemBusyIcon    lipgloss.Style
	ItemErrorIcon   lipgloss.Style
	ItemOnlineIcon  lipgloss.Style

	// Editor: Yolo Mode.
	YoloIconFocused lipgloss.Style
	YoloIconBlurred lipgloss.Style
	YoloDotsFocused lipgloss.Style
	YoloDotsBlurred lipgloss.Style

	// oAuth Chooser.
	AuthBorderSelected   lipgloss.Style
	AuthTextSelected     lipgloss.Style
	AuthBorderUnselected lipgloss.Style
	AuthTextUnselected   lipgloss.Style

	// Dialog styles (can be animated since dialogs render fresh each frame)
	DialogTitle lipgloss.Style
	DialogText  lipgloss.Style
	DialogHint  lipgloss.Style

	// Help bar styles (can be animated)
	HelpStyles help.Styles

	styles     *Styles
	stylesOnce sync.Once
	stylesMu   sync.RWMutex // Protects styles and stylesOnce for thread-safe reloading

	// Animation state
	animStartTime    time.Time
	animHueOffset    float64
	animLastUpdate   time.Time                       // Track last update time for rate limiting
	animColors       atomic.Pointer[AnimatedColors]  // Atomic snapshot for lock-free reads
}

// AnimatedColors holds a snapshot of animated color values.
// This struct is swapped atomically to ensure components see consistent colors.
type AnimatedColors struct {
	Primary       color.Color
	Secondary     color.Color
	Tertiary      color.Color
	Accent        color.Color
	BorderFocus   color.Color
	Success       color.Color
	Error         color.Color
	Warning       color.Color
	Info          color.Color
	BgBase        color.Color
	BgBaseLighter color.Color
	BgSubtle      color.Color
	HueOffset     float64
}

// Color8Bit wraps a color.Color to ensure RGBA() returns 8-bit scaled values.
// This fixes compatibility with lipgloss which expects 8-bit (0-255) values
// but colorful.Color returns 16-bit (0-65535) values.
type Color8Bit struct {
	c color.Color
}

// NewColor8Bit wraps a color to ensure proper 8-bit RGBA conversion.
func NewColor8Bit(c color.Color) Color8Bit {
	return Color8Bit{c: c}
}

// RGBA returns the color's RGBA values scaled to 8-bit range (0-255)
// represented as 16-bit values (0-65535) as required by the color.Color interface.
// This ensures lipgloss interprets the color correctly.
func (c Color8Bit) RGBA() (r, g, b, a uint32) {
	r, g, b, a = c.c.RGBA()
	// Convert from 16-bit to 8-bit, then back to 16-bit format
	// This ensures lipgloss sees values it can properly convert
	r = (r >> 8) * 0x101
	g = (g >> 8) * 0x101
	b = (b >> 8) * 0x101
	a = (a >> 8) * 0x101
	return r, g, b, a
}

type Styles struct {
	Base         lipgloss.Style
	SelectedBase lipgloss.Style

	Title        lipgloss.Style
	Subtitle     lipgloss.Style
	Text         lipgloss.Style
	TextSelected lipgloss.Style
	Muted        lipgloss.Style
	Subtle       lipgloss.Style

	Success lipgloss.Style
	Error   lipgloss.Style
	Warning lipgloss.Style
	Info    lipgloss.Style

	// Markdown & Chroma
	Markdown ansi.StyleConfig

	// Inputs
	TextInput textinput.Styles
	TextArea  textarea.Styles

	// Help
	Help help.Styles

	// Diff
	Diff diffview.Style

	// FilePicker
	FilePicker filepicker.Styles
}

func (t *Theme) S() *Styles {
	// Fast path: check if styles already exist with read lock
	t.stylesMu.RLock()
	if t.styles != nil {
		s := t.styles
		t.stylesMu.RUnlock()
		return s
	}
	t.stylesMu.RUnlock()

	// Slow path: need to build styles with write lock
	t.stylesMu.Lock()
	defer t.stylesMu.Unlock()

	// Double-check after acquiring write lock
	if t.styles != nil {
		return t.styles
	}

	t.styles = t.buildStyles()
	return t.styles
}

// IsAnimated returns true if this theme has animated colors.
func (t *Theme) IsAnimated() bool {
	return t.Animated && t.ColorFunc != nil
}

// StartAnimation initializes the animation start time.
func (t *Theme) StartAnimation() {
	t.stylesMu.Lock()
	defer t.stylesMu.Unlock()
	now := time.Now()
	t.animStartTime = now
	t.animLastUpdate = now
	t.animHueOffset = 0

	// Initialize the atomic color snapshot
	if t.ColorFunc != nil {
		colors := t.ColorFunc(0, 0)
		if len(colors) >= 12 {
			snapshot := &AnimatedColors{
				Primary:       NewColor8Bit(colors[0]),
				Secondary:     NewColor8Bit(colors[1]),
				Tertiary:      NewColor8Bit(colors[2]),
				Accent:        NewColor8Bit(colors[3]),
				BorderFocus:   NewColor8Bit(colors[4]),
				Success:       NewColor8Bit(colors[5]),
				Error:         NewColor8Bit(colors[6]),
				Warning:       NewColor8Bit(colors[7]),
				Info:          NewColor8Bit(colors[8]),
				BgBase:        NewColor8Bit(colors[9]),
				BgBaseLighter: NewColor8Bit(colors[10]),
				BgSubtle:      NewColor8Bit(colors[11]),
				HueOffset:     0,
			}
			t.animColors.Store(snapshot)
		}
	}
}

// AdvanceAnimation updates the theme colors based on elapsed time.
// Returns true if colors were updated and styles need rebuilding.
func (t *Theme) AdvanceAnimation() bool {
	if !t.IsAnimated() {
		return false
	}

	now := time.Now()

	// Calculate time since last update for rate limiting
	timeSinceLastUpdate := now.Sub(t.animLastUpdate).Seconds()
	if timeSinceLastUpdate < 0 {
		timeSinceLastUpdate = 0
	}

	// Calculate new hue offset based on elapsed time from start
	elapsed := now.Sub(t.animStartTime).Seconds()
	newHueOffset := math.Mod(elapsed*t.AnimationSpeed, 360)

	// Sanity check: ensure hue is valid (not NaN, Inf, or negative)
	if math.IsNaN(newHueOffset) || math.IsInf(newHueOffset, 0) || newHueOffset < 0 {
		newHueOffset = t.animHueOffset // Keep current value if calculation fails
	}

	// Calculate hue difference accounting for circular wrap-around (0-360)
	oldOffset := t.animHueOffset
	hueDiff := newHueOffset - oldOffset

	// Normalize to -180 to +180 range for proper wrap handling
	if hueDiff > 180 {
		hueDiff -= 360
	} else if hueDiff < -180 {
		hueDiff += 360
	}

	absHueDiff := math.Abs(hueDiff)

	// Only update if hue changed significantly (at least 1 degree)
	if absHueDiff < 1 {
		return false
	}

	// Rate limit: max hue change based on time elapsed (with 2x safety margin)
	// At 30 deg/sec, in 100ms we should move at most 3 degrees
	// Allow 2x for timing jitter, so max 6 degrees per 100ms
	maxAllowedChange := timeSinceLastUpdate * t.AnimationSpeed * 2.0
	if maxAllowedChange < 5 {
		maxAllowedChange = 5 // Minimum threshold to allow animation to start
	}

	// If the jump is too large, something is wrong - use incremental update instead
	if absHueDiff > maxAllowedChange && t.animLastUpdate.After(t.animStartTime) {
		// Clamp to maximum allowed change in the same direction
		if hueDiff > 0 {
			newHueOffset = math.Mod(oldOffset+maxAllowedChange, 360)
		} else {
			newHueOffset = math.Mod(oldOffset-maxAllowedChange+360, 360)
		}
	}

	t.animHueOffset = newHueOffset
	t.animLastUpdate = now

	// Call the color function to generate new colors
	colors := t.ColorFunc(0, t.animHueOffset)

	// Create a new atomic snapshot of animated colors
	// This is swapped atomically so readers see a consistent set of colors
	// Wrap colors in Color8Bit to ensure proper RGBA conversion for lipgloss
	if len(colors) >= 12 {
		snapshot := &AnimatedColors{
			Primary:       NewColor8Bit(colors[0]),
			Secondary:     NewColor8Bit(colors[1]),
			Tertiary:      NewColor8Bit(colors[2]),
			Accent:        NewColor8Bit(colors[3]),
			BorderFocus:   NewColor8Bit(colors[4]),
			Success:       NewColor8Bit(colors[5]),
			Error:         NewColor8Bit(colors[6]),
			Warning:       NewColor8Bit(colors[7]),
			Info:          NewColor8Bit(colors[8]),
			BgBase:        NewColor8Bit(colors[9]),
			BgBaseLighter: NewColor8Bit(colors[10]),
			BgSubtle:      NewColor8Bit(colors[11]),
			HueOffset:     t.animHueOffset,
		}
		t.animColors.Store(snapshot)

		// Also update the theme fields for backward compatibility
		t.Primary = NewColor8Bit(colors[0])
		t.Secondary = NewColor8Bit(colors[1])
		t.Tertiary = NewColor8Bit(colors[2])
		t.Accent = NewColor8Bit(colors[3])
		t.BorderFocus = NewColor8Bit(colors[4])
		t.Success = NewColor8Bit(colors[5])
		t.Error = NewColor8Bit(colors[6])
		t.Warning = NewColor8Bit(colors[7])
		t.Info = NewColor8Bit(colors[8])
		t.BgBase = NewColor8Bit(colors[9])
		t.BgBaseLighter = NewColor8Bit(colors[10])
		t.BgSubtle = NewColor8Bit(colors[11])
	}

	// Rebuild lipgloss styles if the theme has a style builder
	if t.StyleBuilder != nil {
		t.StyleBuilder(t, t.animHueOffset)
	}

	// Rebuild cached styles
	t.stylesMu.Lock()
	t.styles = t.buildStyles()
	t.stylesMu.Unlock()

	return true
}

// GetAnimatedColors returns the current animated color snapshot.
// This is safe to call concurrently and returns a consistent set of colors.
func (t *Theme) GetAnimatedColors() *AnimatedColors {
	if !t.IsAnimated() {
		return nil
	}
	return t.animColors.Load()
}

// InitAnimatedColors initializes the atomic color snapshot from the theme's current colors.
// Call this after creating an animated theme to ensure GetAnimatedColors never returns nil.
func (t *Theme) InitAnimatedColors() {
	if !t.Animated {
		return
	}
	snapshot := &AnimatedColors{
		Primary:       t.Primary,
		Secondary:     t.Secondary,
		Tertiary:      t.Tertiary,
		Accent:        t.Accent,
		BorderFocus:   t.BorderFocus,
		Success:       t.Success,
		Error:         t.Error,
		Warning:       t.Warning,
		Info:          t.Info,
		BgBase:        t.BgBase,
		BgBaseLighter: t.BgBaseLighter,
		BgSubtle:      t.BgSubtle,
		HueOffset:     0,
	}
	t.animColors.Store(snapshot)
}

// GetHueOffset returns the current animation hue offset.
// Uses the atomic snapshot for lock-free access.
func (t *Theme) GetHueOffset() float64 {
	if snapshot := t.animColors.Load(); snapshot != nil {
		return snapshot.HueOffset
	}
	return t.animHueOffset
}

// cursorColor returns the color to use for the text cursor.
// Uses FgCursor if set, otherwise falls back to Secondary.
func (t *Theme) cursorColor() color.Color {
	if t.FgCursor != nil {
		return t.FgCursor
	}
	return t.Secondary
}

func (t *Theme) buildStyles() *Styles {
	base := lipgloss.NewStyle().
		Foreground(t.FgBase)
	return &Styles{
		Base: base,

		SelectedBase: base.Background(t.Primary),

		Title: base.
			Foreground(t.Accent).
			Bold(true),

		Subtitle: base.
			Foreground(t.Secondary).
			Bold(true),

		Text:         base,
		TextSelected: base.Background(t.Primary).Foreground(t.FgSelected),

		Muted: base.Foreground(t.FgMuted),

		Subtle: base.Foreground(t.FgSubtle),

		Success: base.Foreground(t.Success),

		Error: base.Foreground(t.Error),

		Warning: base.Foreground(t.Warning),

		Info: base.Foreground(t.Info),

		TextInput: textinput.Styles{
			Focused: textinput.StyleState{
				Text:        base,
				Placeholder: base.Foreground(t.FgSubtle),
				Prompt:      base.Foreground(t.Tertiary),
				Suggestion:  base.Foreground(t.FgSubtle),
			},
			Blurred: textinput.StyleState{
				Text:        base.Foreground(t.FgMuted),
				Placeholder: base.Foreground(t.FgSubtle),
				Prompt:      base.Foreground(t.FgMuted),
				Suggestion:  base.Foreground(t.FgSubtle),
			},
			Cursor: textinput.CursorStyle{
				Color: t.cursorColor(),
				Shape: DefaultManager().CursorStyle(),
				Blink: true,
			},
		},
		TextArea: textarea.Styles{
			Focused: textarea.StyleState{
				Base:             base,
				Text:             base,
				LineNumber:       base.Foreground(t.FgSubtle),
				CursorLine:       base,
				CursorLineNumber: base.Foreground(t.FgSubtle),
				Placeholder:      base.Foreground(t.FgSubtle),
				Prompt:           base.Foreground(t.Tertiary),
			},
			Blurred: textarea.StyleState{
				Base:             base,
				Text:             base.Foreground(t.FgMuted),
				LineNumber:       base.Foreground(t.FgMuted),
				CursorLine:       base,
				CursorLineNumber: base.Foreground(t.FgMuted),
				Placeholder:      base.Foreground(t.FgSubtle),
				Prompt:           base.Foreground(t.FgMuted),
			},
			Cursor: textarea.CursorStyle{
				Color: t.cursorColor(),
				Shape: DefaultManager().CursorStyle(),
				Blink: true,
			},
		},

		Markdown: ansi.StyleConfig{
			Document: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					// BlockPrefix: "\n",
					// BlockSuffix: "\n",
					Color: stringPtr(charmtone.Smoke.Hex()),
				},
				// Margin: uintPtr(defaultMargin),
			},
			BlockQuote: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{},
				Indent:         uintPtr(1),
				IndentToken:    stringPtr("│ "),
			},
			List: ansi.StyleList{
				LevelIndent: defaultListIndent,
			},
			Heading: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					BlockSuffix: "\n",
					Color:       stringPtr(charmtone.Malibu.Hex()),
					Bold:        boolPtr(true),
				},
			},
			H1: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Prefix:          " ",
					Suffix:          " ",
					Color:           stringPtr(charmtone.Zest.Hex()),
					BackgroundColor: stringPtr(charmtone.Charple.Hex()),
					Bold:            boolPtr(true),
				},
			},
			H2: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Prefix: "## ",
				},
			},
			H3: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Prefix: "### ",
				},
			},
			H4: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Prefix: "#### ",
				},
			},
			H5: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Prefix: "##### ",
				},
			},
			H6: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Prefix: "###### ",
					Color:  stringPtr(charmtone.Guac.Hex()),
					Bold:   boolPtr(false),
				},
			},
			Strikethrough: ansi.StylePrimitive{
				CrossedOut: boolPtr(true),
			},
			Emph: ansi.StylePrimitive{
				Italic: boolPtr(true),
			},
			Strong: ansi.StylePrimitive{
				Bold: boolPtr(true),
			},
			HorizontalRule: ansi.StylePrimitive{
				Color:  stringPtr(charmtone.Charcoal.Hex()),
				Format: "\n--------\n",
			},
			Item: ansi.StylePrimitive{
				BlockPrefix: "• ",
			},
			Enumeration: ansi.StylePrimitive{
				BlockPrefix: ". ",
			},
			Task: ansi.StyleTask{
				StylePrimitive: ansi.StylePrimitive{},
				Ticked:         "[✓] ",
				Unticked:       "[ ] ",
			},
			Link: ansi.StylePrimitive{
				Color:     stringPtr(charmtone.Zinc.Hex()),
				Underline: boolPtr(true),
			},
			LinkText: ansi.StylePrimitive{
				Color: stringPtr(charmtone.Guac.Hex()),
				Bold:  boolPtr(true),
			},
			Image: ansi.StylePrimitive{
				Color:     stringPtr(charmtone.Cheeky.Hex()),
				Underline: boolPtr(true),
			},
			ImageText: ansi.StylePrimitive{
				Color:  stringPtr(charmtone.Squid.Hex()),
				Format: "Image: {{.text}} →",
			},
			Code: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Prefix:          " ",
					Suffix:          " ",
					Color:           stringPtr(charmtone.Coral.Hex()),
					BackgroundColor: stringPtr(charmtone.Charcoal.Hex()),
				},
			},
			CodeBlock: ansi.StyleCodeBlock{
				StyleBlock: ansi.StyleBlock{
					StylePrimitive: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Charcoal.Hex()),
					},
					Margin: uintPtr(defaultMargin),
				},
				Chroma: &ansi.Chroma{
					Text: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Smoke.Hex()),
					},
					Error: ansi.StylePrimitive{
						Color:           stringPtr(charmtone.Butter.Hex()),
						BackgroundColor: stringPtr(charmtone.Sriracha.Hex()),
					},
					Comment: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Oyster.Hex()),
					},
					CommentPreproc: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Bengal.Hex()),
					},
					Keyword: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Malibu.Hex()),
					},
					KeywordReserved: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Pony.Hex()),
					},
					KeywordNamespace: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Pony.Hex()),
					},
					KeywordType: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Guppy.Hex()),
					},
					Operator: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Salmon.Hex()),
					},
					Punctuation: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Zest.Hex()),
					},
					Name: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Smoke.Hex()),
					},
					NameBuiltin: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Cheeky.Hex()),
					},
					NameTag: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Mauve.Hex()),
					},
					NameAttribute: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Hazy.Hex()),
					},
					NameClass: ansi.StylePrimitive{
						Color:     stringPtr(charmtone.Salt.Hex()),
						Underline: boolPtr(true),
						Bold:      boolPtr(true),
					},
					NameDecorator: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Citron.Hex()),
					},
					NameFunction: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Guac.Hex()),
					},
					LiteralNumber: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Julep.Hex()),
					},
					LiteralString: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Cumin.Hex()),
					},
					LiteralStringEscape: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Bok.Hex()),
					},
					GenericDeleted: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Coral.Hex()),
					},
					GenericEmph: ansi.StylePrimitive{
						Italic: boolPtr(true),
					},
					GenericInserted: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Guac.Hex()),
					},
					GenericStrong: ansi.StylePrimitive{
						Bold: boolPtr(true),
					},
					GenericSubheading: ansi.StylePrimitive{
						Color: stringPtr(charmtone.Squid.Hex()),
					},
					Background: ansi.StylePrimitive{
						BackgroundColor: stringPtr(charmtone.Charcoal.Hex()),
					},
				},
			},
			Table: ansi.StyleTable{
				StyleBlock: ansi.StyleBlock{
					StylePrimitive: ansi.StylePrimitive{},
				},
			},
			DefinitionDescription: ansi.StylePrimitive{
				BlockPrefix: "\n ",
			},
		},

		Help: help.Styles{
			ShortKey:       base.Foreground(t.FgMuted),
			ShortDesc:      base.Foreground(t.FgSubtle),
			ShortSeparator: base.Foreground(t.Border),
			Ellipsis:       base.Foreground(t.Border),
			FullKey:        base.Foreground(t.FgMuted),
			FullDesc:       base.Foreground(t.FgSubtle),
			FullSeparator:  base.Foreground(t.Border),
		},

		Diff: diffview.Style{
			DividerLine: diffview.LineStyle{
				LineNumber: lipgloss.NewStyle().
					Foreground(t.FgHalfMuted).
					Background(t.BgBaseLighter),
				Code: lipgloss.NewStyle().
					Foreground(t.FgHalfMuted).
					Background(t.BgBaseLighter),
			},
			MissingLine: diffview.LineStyle{
				LineNumber: lipgloss.NewStyle().
					Background(t.BgBaseLighter),
				Code: lipgloss.NewStyle().
					Background(t.BgBaseLighter),
			},
			EqualLine: diffview.LineStyle{
				LineNumber: lipgloss.NewStyle().
					Foreground(t.FgMuted).
					Background(t.BgBase),
				Code: lipgloss.NewStyle().
					Foreground(t.FgMuted).
					Background(t.BgBase),
			},
			InsertLine: diffview.LineStyle{
				LineNumber: lipgloss.NewStyle().
					Foreground(lipgloss.Color("#629657")).
					Background(lipgloss.Color("#2b322a")),
				Symbol: lipgloss.NewStyle().
					Foreground(lipgloss.Color("#629657")).
					Background(lipgloss.Color("#323931")),
				Code: lipgloss.NewStyle().
					Background(lipgloss.Color("#323931")),
			},
			DeleteLine: diffview.LineStyle{
				LineNumber: lipgloss.NewStyle().
					Foreground(lipgloss.Color("#a45c59")).
					Background(lipgloss.Color("#312929")),
				Symbol: lipgloss.NewStyle().
					Foreground(lipgloss.Color("#a45c59")).
					Background(lipgloss.Color("#383030")),
				Code: lipgloss.NewStyle().
					Background(lipgloss.Color("#383030")),
			},
		},
		FilePicker: filepicker.Styles{
			DisabledCursor:   base.Foreground(t.FgMuted),
			Cursor:           base.Foreground(t.FgBase),
			Symlink:          base.Foreground(t.FgSubtle),
			Directory:        base.Foreground(t.Primary),
			File:             base.Foreground(t.FgBase),
			DisabledFile:     base.Foreground(t.FgMuted),
			DisabledSelected: base.Background(t.BgOverlay).Foreground(t.FgMuted),
			Permission:       base.Foreground(t.FgMuted),
			Selected:         base.Background(t.Primary).Foreground(t.FgBase),
			FileSize:         base.Foreground(t.FgMuted),
			EmptyDirectory:   base.Foreground(t.FgMuted).PaddingLeft(2).SetString("Empty directory"),
		},
	}
}

type Manager struct {
	themes      map[string]*Theme
	current     *Theme
	cursorStyle tea.CursorShape
	scrollStep  int
}

var (
	defaultManager     *Manager
	defaultManagerOnce sync.Once
)

func initDefaultManager() *Manager {
	defaultManagerOnce.Do(func() {
		defaultManager = newManager()
	})
	return defaultManager
}

func SetDefaultManager(m *Manager) {
	defaultManager = m
}

func DefaultManager() *Manager {
	return initDefaultManager()
}

func CurrentTheme() *Theme {
	return initDefaultManager().Current()
}

// pendingThemes holds themes registered before the manager is initialized.
var pendingThemes []*Theme

// RegisterTheme registers a theme. Call this from theme init() functions.
func RegisterTheme(theme *Theme) {
	if defaultManager != nil {
		defaultManager.Register(theme)
	} else {
		pendingThemes = append(pendingThemes, theme)
	}
}

func newManager() *Manager {
	m := &Manager{
		themes: make(map[string]*Theme),
	}

	// Register themes that were added before manager initialization
	for _, theme := range pendingThemes {
		m.Register(theme)
	}
	pendingThemes = nil

	// Set default theme (fallback to first registered if "dark" not found)
	if theme, ok := m.themes["dark"]; ok {
		m.current = theme
	} else if len(m.themes) > 0 {
		for _, theme := range m.themes {
			m.current = theme
			break
		}
	} else {
		// Fallback theme for tests that don't import themes package
		m.current = newFallbackTheme()
	}

	return m
}

// newFallbackTheme creates a minimal theme for testing.
func newFallbackTheme() *Theme {
	return &Theme{
		Name:        "fallback",
		IsDark:      true,
		Primary:     charmtone.Iron,
		Secondary:   charmtone.Smoke,
		Tertiary:    charmtone.Ash,
		Accent:      charmtone.Malibu,
		BgBase:      charmtone.Pepper,
		BgSubtle:    charmtone.Charcoal,
		BgOverlay:   charmtone.Iron,
		FgBase:      charmtone.Ash,
		FgMuted:     charmtone.Oyster,
		FgHalfMuted: charmtone.Smoke,
		FgSubtle:    charmtone.Squid,
		FgSelected:  charmtone.Salt,
		Border:      charmtone.Charcoal,
		BorderFocus: charmtone.Malibu,
		Success:     charmtone.Guac,
		Error:       charmtone.Coral,
		Warning:     charmtone.Citron,
		Info:        charmtone.Malibu,
	}
}

func (m *Manager) Register(theme *Theme) {
	m.themes[theme.Name] = theme
}

func (m *Manager) Current() *Theme {
	return m.current
}

func (m *Manager) SetTheme(name string) error {
	if theme, ok := m.themes[name]; ok {
		m.current = theme
		return nil
	}
	return fmt.Errorf("theme %s not found", name)
}

// SetThemeWithNotify sets the theme and returns a ThemeChangedMsg for the TUI to handle.
func (m *Manager) SetThemeWithNotify(name string) (tea.Msg, error) {
	if err := m.SetTheme(name); err != nil {
		return nil, err
	}
	return ThemeChangedMsg{ThemeName: name}, nil
}

func (m *Manager) List() []string {
	names := make([]string, 0, len(m.themes))
	for name := range m.themes {
		names = append(names, name)
	}
	return names
}

// SetCursorStyle sets the cursor style from a string (bar, block, underline).
func (m *Manager) SetCursorStyle(style string) {
	switch strings.ToLower(style) {
	case "block":
		m.cursorStyle = tea.CursorBlock
	case "underline":
		m.cursorStyle = tea.CursorUnderline
	default:
		m.cursorStyle = tea.CursorBar
	}
}

// CursorStyle returns the configured cursor style.
func (m *Manager) CursorStyle() tea.CursorShape {
	if m.cursorStyle == 0 {
		return tea.CursorBar // default
	}
	return m.cursorStyle
}

// SetScrollStep sets the scroll step size.
func (m *Manager) SetScrollStep(step int) {
	if step < 1 {
		step = 1
	} else if step > 10 {
		step = 10
	}
	m.scrollStep = step
}

// ScrollStep returns the configured scroll step size.
func (m *Manager) ScrollStep() int {
	if m.scrollStep == 0 {
		return 2 // default
	}
	return m.scrollStep
}

// ReloadThemes resets all theme styles so they will be rebuilt on next access.
// This allows hot-reloading of cursor style and other theme options.
func (m *Manager) ReloadThemes() {
	for _, theme := range m.themes {
		theme.stylesMu.Lock()
		theme.styles = nil
		theme.stylesMu.Unlock()
	}
}

// ParseHex converts hex string to color
func ParseHex(hex string) color.Color {
	var r, g, b uint8
	fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

// Alpha returns a color with transparency
func Alpha(c color.Color, alpha uint8) color.Color {
	r, g, b, _ := c.RGBA()
	return color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: alpha,
	}
}

// Darken makes a color darker by percentage (0-100)
func Darken(c color.Color, percent float64) color.Color {
	r, g, b, a := c.RGBA()
	factor := 1.0 - percent/100.0
	return color.RGBA{
		R: uint8(float64(r>>8) * factor),
		G: uint8(float64(g>>8) * factor),
		B: uint8(float64(b>>8) * factor),
		A: uint8(a >> 8),
	}
}

// Lighten makes a color lighter by percentage (0-100)
func Lighten(c color.Color, percent float64) color.Color {
	r, g, b, a := c.RGBA()
	factor := percent / 100.0
	return color.RGBA{
		R: uint8(min(255, float64(r>>8)+255*factor)),
		G: uint8(min(255, float64(g>>8)+255*factor)),
		B: uint8(min(255, float64(b>>8)+255*factor)),
		A: uint8(a >> 8),
	}
}

func ForegroundGrad(input string, bold bool, color1, color2 color.Color) []string {
	return colors.ForegroundGrad(CurrentTheme().S().Base, input, bold, color1, color2)
}

// ApplyForegroundGrad renders a given string with a horizontal gradient
// foreground.
func ApplyForegroundGrad(input string, color1, color2 color.Color) string {
	return colors.ApplyForegroundGrad(CurrentTheme().S().Base, input, color1, color2)
}

// PerceivedBrightnessAdjust returns saturation reduction for perceptually
// darker or more intense hues. Blue/violet appear darker and red appears
// more intense than other colors at the same HSV values, so we compensate.
// This is the single source of truth for chroma brightness compensation.
func PerceivedBrightnessAdjust(h float64) (saturationReduce, valueBoost float64) {
	h = math.Mod(h+360, 360)

	// Blue/Violet adjustment (25% max saturation reduction)
	// Flat-top curve covering blue (240°) and violet (280°)
	blueVioletAdj := func() float64 {
		blueCenter := 240.0
		violetCenter := 280.0
		midpoint := (blueCenter + violetCenter) / 2 // 260°

		dist := math.Abs(h - midpoint)
		if dist > 180 {
			dist = 360 - dist
		}

		flatWidth := (violetCenter - blueCenter) / 2 // 20°
		taperWidth := 50.0

		if dist <= flatWidth {
			return 1.0
		} else if dist <= flatWidth+taperWidth {
			taperDist := dist - flatWidth
			return (1 + math.Cos(math.Pi*taperDist/taperWidth)) / 2
		}
		return 0
	}()

	// Red adjustment (15% max saturation reduction)
	// Centered at 0°/360° with smooth taper
	redAdj := func() float64 {
		dist := h
		if dist > 180 {
			dist = 360 - dist
		}

		taperWidth := 40.0
		if dist <= taperWidth {
			return (1 + math.Cos(math.Pi*dist/taperWidth)) / 2
		}
		return 0
	}()

	// Combine adjustments (they don't overlap, so take max)
	blueVioletSatReduce := 0.25 * blueVioletAdj
	redSatReduce := 0.15 * redAdj
	saturationReduce = math.Max(blueVioletSatReduce, redSatReduce)

	// No value boost needed since we're using saturation reduction
	valueBoost = 0

	return saturationReduce, valueBoost
}

// ApplyAnimatedGrad renders a string with a sweeping rainbow gradient that
// moves based on the theme's current animation offset. Creates a wave effect.
func ApplyAnimatedGrad(input string) string {
	if input == "" {
		return ""
	}
	t := CurrentTheme()
	if !t.IsAnimated() {
		return t.S().Text.Render(input)
	}

	// Get current hue offset (0-360)
	hueOffset := t.GetHueOffset()

	// Create a sweeping rainbow by generating colors based on position + offset
	var clusters []string
	gr := uniseg.NewGraphemes(input)
	for gr.Next() {
		clusters = append(clusters, string(gr.Runes()))
	}

	if len(clusters) == 0 {
		return ""
	}

	// Generate rainbow colors that sweep based on hueOffset
	// Each character gets a hue shifted by position, plus the animation offset
	var o strings.Builder
	for i, cluster := range clusters {
		// Spread 60 degrees of hue across the text, offset by animation
		charHue := hueOffset + float64(i)*60.0/float64(len(clusters))
		charHue = math.Mod(charHue, 360)

		// Apply perceptual brightness compensation for blue/purple hues
		// Reducing saturation makes colors lighter (more pastel = brighter)
		satReduce, valBoost := PerceivedBrightnessAdjust(charHue)
		saturation := math.Max(0.3, 0.7-satReduce) // base 0.7, reduce for dark hues
		value := math.Min(1.0, 1.0+valBoost)

		c := colorful.Hsv(charHue, saturation, value)
		style := t.S().Base.Foreground(c)
		o.WriteString(style.Render(cluster))
	}
	return o.String()
}

// ApplyBoldForegroundGrad renders a given string with a horizontal gradient
// foreground.
func ApplyBoldForegroundGrad(input string, color1, color2 color.Color) string {
	return colors.ApplyBoldForegroundGrad(CurrentTheme().S().Base, input, color1, color2)
}

// Animation support for themes

// ThemeChangedMsg is sent when the theme is changed.
// The TUI should handle this to start/stop animation as needed.
type ThemeChangedMsg struct {
	ThemeName string
}

// AnimationTickMsg is sent when theme animation should advance.
type AnimationTickMsg struct{}

// AnimationTickInterval is the interval between animation ticks (20hz = 50ms).
// Uses atomic color snapshots for lock-free reads to prevent flicker.
const AnimationTickInterval = 50 * time.Millisecond

// AnimationTickCmd returns a command that sends AnimationTickMsg after the tick interval.
func AnimationTickCmd() tea.Cmd {
	return tea.Tick(AnimationTickInterval, func(time.Time) tea.Msg {
		return AnimationTickMsg{}
	})
}

// StartAnimationIfNeeded starts animation for the current theme if it's animated.
// Returns a command to start the animation tick loop.
func (m *Manager) StartAnimationIfNeeded() tea.Cmd {
	if m.current != nil && m.current.IsAnimated() {
		m.current.StartAnimation()
		return AnimationTickCmd()
	}
	return nil
}

// HandleAnimationTick advances the current theme's animation and returns
// a command to continue the animation tick loop if needed.
func (m *Manager) HandleAnimationTick() tea.Cmd {
	if m.current != nil && m.current.IsAnimated() {
		m.current.AdvanceAnimation()
		return AnimationTickCmd()
	}
	return nil
}
