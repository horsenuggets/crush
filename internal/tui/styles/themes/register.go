package themes

import "github.com/charmbracelet/crush/internal/tui/styles"

func init() {
	styles.RegisterTheme(NewDarkTheme())
	styles.RegisterTheme(NewCharmtoneTheme())
	styles.RegisterTheme(NewLightTheme())
}
