package styles

import (
	"image/color"

	"github.com/charmbracelet/crush/internal/shared/colors"
)

// ForegroundGrad returns a slice of strings representing the input string
// rendered with a horizontal gradient foreground from color1 to color2. Each
// string in the returned slice corresponds to a grapheme cluster in the input
// string. If bold is true, the rendered strings will be bolded.
func ForegroundGrad(t *Styles, input string, bold bool, color1, color2 color.Color) []string {
	return colors.ForegroundGrad(t.Base, input, bold, color1, color2)
}

// ApplyForegroundGrad renders a given string with a horizontal gradient
// foreground.
func ApplyForegroundGrad(t *Styles, input string, color1, color2 color.Color) string {
	return colors.ApplyForegroundGrad(t.Base, input, color1, color2)
}

// ApplyBoldForegroundGrad renders a given string with a horizontal gradient
// foreground.
func ApplyBoldForegroundGrad(t *Styles, input string, color1, color2 color.Color) string {
	return colors.ApplyBoldForegroundGrad(t.Base, input, color1, color2)
}
