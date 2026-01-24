package themes

import "github.com/charmbracelet/crush/internal/tui/styles"

func init() {
	styles.RegisterTheme(NewCharmtoneTheme())
	styles.RegisterTheme(NewChromaTheme())
	styles.RegisterTheme(NewDarkTheme())
	styles.RegisterTheme(NewGitHubTheme())
	styles.RegisterTheme(NewHackerTheme())
	styles.RegisterTheme(NewLightTheme())
	styles.RegisterTheme(NewMonokaiTheme())
	styles.RegisterTheme(NewOceanTheme())
	styles.RegisterTheme(NewSolarizedTheme())
	styles.RegisterTheme(NewSugarcookieTheme())
}
