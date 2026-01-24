package themes

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/crush/internal/tui/styles"
)

// GitHub Dark theme colors based on GitHub's dark mode.
var (
	ghBgBase        = styles.ParseHex("#0d1117")
	ghBgBaseLighter = styles.ParseHex("#161b22")
	ghBgSubtle      = styles.ParseHex("#21262d")
	ghBgOverlay     = styles.ParseHex("#30363d")

	ghFgBase      = styles.ParseHex("#e6edf3")
	ghFgMuted     = styles.ParseHex("#7d8590")
	ghFgHalfMuted = styles.ParseHex("#8b949e")
	ghFgSubtle    = styles.ParseHex("#6e7681")
	ghFgSelected  = styles.ParseHex("#ffffff")

	ghBorder      = styles.ParseHex("#30363d")
	ghBorderFocus = styles.ParseHex("#2f81f7")

	ghPrimary   = styles.ParseHex("#2f81f7") // blue
	ghSecondary = styles.ParseHex("#8b949e")
	ghTertiary  = styles.ParseHex("#6e7681")
	ghAccent    = styles.ParseHex("#a371f7") // purple

	ghSuccess = styles.ParseHex("#3fb950")
	ghError   = styles.ParseHex("#f85149")
	ghWarning = styles.ParseHex("#d29922")
	ghInfo    = styles.ParseHex("#2f81f7")

	ghWhite     = styles.ParseHex("#ffffff")
	ghBlueLight = styles.ParseHex("#79c0ff")
	ghBlueDark  = styles.ParseHex("#1f6feb")
	ghBlue      = styles.ParseHex("#2f81f7")

	ghYellow = styles.ParseHex("#d29922")
	ghCitron = styles.ParseHex("#bb8009")

	ghGreen      = styles.ParseHex("#3fb950")
	ghGreenDark  = styles.ParseHex("#238636")
	ghGreenLight = styles.ParseHex("#56d364")

	ghRed      = styles.ParseHex("#f85149")
	ghRedDark  = styles.ParseHex("#da3633")
	ghRedLight = styles.ParseHex("#ff7b72")
	ghCherry   = styles.ParseHex("#db61a2")
)

// NewGitHubTheme creates a GitHub dark mode inspired theme.
func NewGitHubTheme() *styles.Theme {
	t := &styles.Theme{
		Name:   "github",
		IsDark: true,

		Primary:   ghPrimary,
		Secondary: ghSecondary,
		Tertiary:  ghTertiary,
		Accent:    ghAccent,

		BgBase:        ghBgBase,
		BgBaseLighter: ghBgBaseLighter,
		BgSubtle:      ghBgSubtle,
		BgOverlay:     ghBgOverlay,

		FgBase:      ghFgBase,
		FgMuted:     ghFgMuted,
		FgHalfMuted: ghFgHalfMuted,
		FgSubtle:    ghFgSubtle,
		FgSelected:  ghFgSelected,

		Border:      ghBorder,
		BorderFocus: ghBorderFocus,

		Success: ghSuccess,
		Error:   ghError,
		Warning: ghWarning,
		Info:    ghInfo,

		White:     ghWhite,
		BlueLight: ghBlueLight,
		BlueDark:  ghBlueDark,
		Blue:      ghBlue,

		Yellow: ghYellow,
		Citron: ghCitron,

		Green:      ghGreen,
		GreenDark:  ghGreenDark,
		GreenLight: ghGreenLight,

		Red:      ghRed,
		RedDark:  ghRedDark,
		RedLight: ghRedLight,
		Cherry:   ghCherry,
	}

	t.TextSelection = lipgloss.NewStyle().Foreground(ghFgSelected).Background(ghPrimary)

	t.ItemOfflineIcon = lipgloss.NewStyle().Foreground(ghFgSubtle).SetString("●")
	t.ItemBusyIcon = t.ItemOfflineIcon.Foreground(ghWarning)
	t.ItemErrorIcon = t.ItemOfflineIcon.Foreground(ghError)
	t.ItemOnlineIcon = t.ItemOfflineIcon.Foreground(ghSuccess)

	t.YoloIconFocused = lipgloss.NewStyle().Foreground(ghBgBase).Background(ghWarning).Bold(true).SetString(" ! ")
	t.YoloIconBlurred = t.YoloIconFocused.Foreground(ghBgBase).Background(ghFgSubtle)
	t.YoloDotsFocused = lipgloss.NewStyle().Foreground(ghWarning).SetString(":::")
	t.YoloDotsBlurred = t.YoloDotsFocused.Foreground(ghFgSubtle)

	t.AuthBorderSelected = lipgloss.NewStyle().BorderForeground(ghSuccess)
	t.AuthTextSelected = lipgloss.NewStyle().Foreground(ghSuccess)
	t.AuthBorderUnselected = lipgloss.NewStyle().BorderForeground(ghBorder)
	t.AuthTextUnselected = lipgloss.NewStyle().Foreground(ghFgMuted)

	return t
}
