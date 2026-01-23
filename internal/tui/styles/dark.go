package styles

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
)

// NewDarkTheme creates a minimal dark theme with neutral colors.
func NewDarkTheme() *Theme {
	t := &Theme{
		Name:   "dark",
		IsDark: true,

		Primary:   charmtone.Iron,
		Secondary: charmtone.Smoke,
		Tertiary:  charmtone.Ash,
		Accent:    charmtone.Malibu,

		// Backgrounds
		BgBase:        charmtone.Pepper,
		BgBaseLighter: charmtone.BBQ,
		BgSubtle:      charmtone.Charcoal,
		BgOverlay:     charmtone.Iron,

		// Foregrounds
		FgBase:      charmtone.Ash,
		FgMuted:     charmtone.Oyster,
		FgHalfMuted: charmtone.Smoke,
		FgSubtle:    charmtone.Squid,
		FgSelected:  charmtone.Salt,

		// Borders
		Border:      charmtone.Charcoal,
		BorderFocus: charmtone.Malibu,

		// Status
		Success: charmtone.Guac,
		Error:   charmtone.Coral,
		Warning: charmtone.Citron,
		Info:    charmtone.Malibu,

		// Colors
		White: charmtone.Salt,

		BlueLight: charmtone.Sardine,
		BlueDark:  charmtone.Damson,
		Blue:      charmtone.Malibu,

		Yellow: charmtone.Mustard,
		Citron: charmtone.Citron,

		Green:      charmtone.Julep,
		GreenDark:  charmtone.Guac,
		GreenLight: charmtone.Bok,

		Red:      charmtone.Coral,
		RedDark:  charmtone.Sriracha,
		RedLight: charmtone.Salmon,
		Cherry:   charmtone.Cherry,
	}

	// Text selection.
	t.TextSelection = lipgloss.NewStyle().Foreground(charmtone.Salt).Background(charmtone.Iron)

	// LSP and MCP status.
	t.ItemOfflineIcon = lipgloss.NewStyle().Foreground(charmtone.Squid).SetString("●")
	t.ItemBusyIcon = t.ItemOfflineIcon.Foreground(charmtone.Citron)
	t.ItemErrorIcon = t.ItemOfflineIcon.Foreground(charmtone.Coral)
	t.ItemOnlineIcon = t.ItemOfflineIcon.Foreground(charmtone.Guac)

	// Editor: Yolo Mode.
	t.YoloIconFocused = lipgloss.NewStyle().Foreground(charmtone.Pepper).Background(charmtone.Citron).Bold(true).SetString(" ! ")
	t.YoloIconBlurred = t.YoloIconFocused.Foreground(charmtone.Pepper).Background(charmtone.Squid)
	t.YoloDotsFocused = lipgloss.NewStyle().Foreground(charmtone.Citron).SetString(":::")
	t.YoloDotsBlurred = t.YoloDotsFocused.Foreground(charmtone.Squid)

	// oAuth Chooser.
	t.AuthBorderSelected = lipgloss.NewStyle().BorderForeground(charmtone.Malibu)
	t.AuthTextSelected = lipgloss.NewStyle().Foreground(charmtone.Malibu)
	t.AuthBorderUnselected = lipgloss.NewStyle().BorderForeground(charmtone.Iron)
	t.AuthTextUnselected = lipgloss.NewStyle().Foreground(charmtone.Squid)

	return t
}
