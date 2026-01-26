package themes

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/tui/styles"
)

// Hacker theme - classic green on black, Matrix-inspired.
var (
	hackerBgBase        = styles.ParseHex("#000000") // Pure black
	hackerBgBaseLighter = styles.ParseHex("#0a0a0a") // Near black
	hackerBgSubtle      = styles.ParseHex("#121212") // Very dark gray
	hackerBgOverlay     = styles.ParseHex("#1c1c1c") // Dark gray

	hackerFgBase      = styles.ParseHex("#00ff41") // Matrix green
	hackerFgMuted     = styles.ParseHex("#00b432") // Darker green
	hackerFgHalfMuted = styles.ParseHex("#00c837") // Medium green
	hackerFgSubtle    = styles.ParseHex("#008c28") // Subtle green
	hackerFgSelected  = styles.ParseHex("#000000") // Black for contrast on green selection

	hackerBorder      = styles.ParseHex("#005019") // Dark green border
	hackerBorderFocus = styles.ParseHex("#00ff41") // Bright green

	hackerPrimary   = styles.ParseHex("#00ff41") // Matrix green
	hackerSecondary = styles.ParseHex("#00c837") // Medium green
	hackerTertiary  = styles.ParseHex("#00a02d") // Darker green
	hackerAccent    = styles.ParseHex("#32ff96") // Cyan-green

	hackerSuccess = styles.ParseHex("#00ff41") // Green
	hackerError   = styles.ParseHex("#ff3232") // Red (stands out)
	hackerWarning = styles.ParseHex("#ffc800") // Amber
	hackerInfo    = styles.ParseHex("#00c8ff") // Cyan

	hackerWhite     = styles.ParseHex("#b4ffb4") // Green-white
	hackerBlueLight = styles.ParseHex("#64ffda") // Cyan
	hackerBlueDark  = styles.ParseHex("#006450") // Dark cyan
	hackerBlue      = styles.ParseHex("#00c8ff") // Cyan

	hackerYellow = styles.ParseHex("#ffc800") // Amber
	hackerCitron = styles.ParseHex("#c8ff00") // Lime

	hackerGreen      = styles.ParseHex("#00ff41") // Matrix green
	hackerGreenDark  = styles.ParseHex("#00b432") // Dark green
	hackerGreenLight = styles.ParseHex("#64ff96") // Light green

	hackerRed      = styles.ParseHex("#ff3232") // Red
	hackerRedDark  = styles.ParseHex("#b41e1e") // Dark red
	hackerRedLight = styles.ParseHex("#ff6464") // Light red
	hackerCherry   = styles.ParseHex("#ff0064") // Magenta
)

// NewHackerTheme creates a Matrix-inspired green on black theme.
func NewHackerTheme() *styles.Theme {
	t := &styles.Theme{
		Name:   "hacker",
		IsDark: true,

		Primary:   hackerPrimary,
		Secondary: hackerSecondary,
		Tertiary:  hackerTertiary,
		Accent:    hackerAccent,

		BgBase:        hackerBgBase,
		BgBaseLighter: hackerBgBaseLighter,
		BgSubtle:      hackerBgSubtle,
		BgOverlay:     hackerBgOverlay,

		FgBase:      hackerFgBase,
		FgMuted:     hackerFgMuted,
		FgHalfMuted: hackerFgHalfMuted,
		FgSubtle:    hackerFgSubtle,
		FgSelected: hackerFgSelected,

		Border:      hackerBorder,
		BorderFocus: hackerBorderFocus,

		Success: hackerSuccess,
		Error:   hackerError,
		Warning: hackerWarning,
		Info:    hackerInfo,

		White:     hackerWhite,
		BlueLight: hackerBlueLight,
		BlueDark:  hackerBlueDark,
		Blue:      hackerBlue,

		Yellow: hackerYellow,
		Citron: hackerCitron,

		Green:      hackerGreen,
		GreenDark:  hackerGreenDark,
		GreenLight: hackerGreenLight,

		Red:      hackerRed,
		RedDark:  hackerRedDark,
		RedLight: hackerRedLight,
		Cherry:   hackerCherry,
	}

	t.TextSelection = lipgloss.NewStyle().Foreground(hackerBgBase).Background(hackerPrimary)

	t.ItemOfflineIcon = lipgloss.NewStyle().Foreground(hackerFgSubtle).SetString("●")
	t.ItemBusyIcon = t.ItemOfflineIcon.Foreground(hackerWarning)
	t.ItemErrorIcon = t.ItemOfflineIcon.Foreground(hackerError)
	t.ItemOnlineIcon = t.ItemOfflineIcon.Foreground(hackerSuccess)

	// Editor: Yolo Mode - using lime green to fit the hacker aesthetic.
	t.YoloIconFocused = lipgloss.NewStyle().Foreground(hackerBgBase).Background(hackerCitron).Bold(true).SetString(" ! ")
	t.YoloIconBlurred = t.YoloIconFocused.Foreground(hackerFgBase).Background(hackerBgOverlay)
	t.YoloDotsFocused = lipgloss.NewStyle().Foreground(hackerCitron).SetString(":::")
	t.YoloDotsBlurred = t.YoloDotsFocused.Foreground(hackerFgSubtle)

	t.AuthBorderSelected = lipgloss.NewStyle().BorderForeground(hackerSuccess)
	t.AuthTextSelected = lipgloss.NewStyle().Foreground(hackerSuccess)
	t.AuthBorderUnselected = lipgloss.NewStyle().BorderForeground(hackerBorder)
	t.AuthTextUnselected = lipgloss.NewStyle().Foreground(hackerFgMuted)

	return t
}
