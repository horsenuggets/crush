package styles_test

import (
	"testing"
	"time"

	"github.com/charmbracelet/crush/internal/tui/styles"
	_ "github.com/charmbracelet/crush/internal/tui/styles/themes"
)

func TestChromaAnimation(t *testing.T) {
	mgr := styles.DefaultManager()

	// Set chroma theme
	if err := mgr.SetTheme("chroma"); err != nil {
		t.Fatalf("Error setting theme: %v", err)
	}

	theme := mgr.Current()
	if theme.Name != "chroma" {
		t.Fatalf("Expected chroma theme, got %s", theme.Name)
	}

	if !theme.IsAnimated() {
		t.Fatal("Chroma theme should be animated")
	}

	// Start animation
	theme.StartAnimation()
	initialHue := theme.GetHueOffset()
	if initialHue != 0 {
		t.Errorf("Initial hue should be 0, got %f", initialHue)
	}

	// Wait and advance animation
	time.Sleep(150 * time.Millisecond)
	changed := theme.AdvanceAnimation()

	newHue := theme.GetHueOffset()
	t.Logf("Initial hue: %f, New hue: %f, Changed: %v", initialHue, newHue, changed)

	// At 30 degrees/second, after 150ms we should have moved ~4.5 degrees
	// The threshold for change is 1 degree, so it should have changed
	if !changed {
		t.Error("Animation should have advanced after 150ms")
	}

	if newHue <= initialHue {
		t.Errorf("Hue should have increased, got %f", newHue)
	}

	// Test that styles are cleared after animation advance
	// Get styles, then advance, then check they need rebuilding
	_ = theme.S() // Build styles

	time.Sleep(150 * time.Millisecond)
	theme.AdvanceAnimation()

	// Styles should still be accessible (will rebuild automatically)
	s := theme.S()
	if s == nil {
		t.Error("Styles should be accessible after animation")
	}
}
