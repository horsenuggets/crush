package main

import (
	"fmt"
	"math"

	"charm.land/lipgloss/v2"
	"github.com/lucasb-eyer/go-colorful"
)

// adjustBlueViolet applies saturation reduction for blue/violet range
func adjustBlueViolet(h, maxSatReduce float64) float64 {
	h = math.Mod(h+360, 360)

	// Flat-top curve covering blue (240°) and violet (280°)
	blueCenter := 240.0
	violetCenter := 280.0
	midpoint := (blueCenter + violetCenter) / 2 // 260°

	dist := math.Abs(h - midpoint)
	if dist > 180 {
		dist = 360 - dist
	}

	flatWidth := (violetCenter - blueCenter) / 2 // 20°
	taperWidth := 50.0

	var factor float64
	if dist <= flatWidth {
		factor = 1.0
	} else if dist <= flatWidth+taperWidth {
		taperDist := dist - flatWidth
		factor = (1 + math.Cos(math.Pi*taperDist/taperWidth)) / 2
	} else {
		return 0
	}

	return maxSatReduce * factor
}

// adjustRed applies saturation reduction for red range (around 0°/360°)
func adjustRed(h, maxSatReduce float64) float64 {
	h = math.Mod(h+360, 360)

	// Red is at 0°/360°, so we need to handle wrap-around
	// Distance from red (0°)
	dist := h
	if dist > 180 {
		dist = 360 - dist
	}

	taperWidth := 40.0 // How far the adjustment extends from red

	var factor float64
	if dist <= taperWidth {
		factor = (1 + math.Cos(math.Pi*dist/taperWidth)) / 2
	} else {
		return 0
	}

	return maxSatReduce * factor
}

// combinedAdjust combines both blue/violet and red adjustments
func combinedAdjust(h, blueSatReduce, redSatReduce float64) float64 {
	blueAdj := adjustBlueViolet(h, blueSatReduce)
	redAdj := adjustRed(h, redSatReduce)
	// Take the max of both adjustments (they shouldn't overlap)
	return math.Max(blueAdj, redAdj)
}

func main() {
	block := "██"
	saturation := 1.0
	value := 1.0

	fmt.Println("Chroma Color Debug - Hue 0° to 360°")
	fmt.Println()

	// Print hue labels
	fmt.Print("Hue:              ")
	for h := 0; h < 360; h += 10 {
		fmt.Printf("%3d ", h)
	}
	fmt.Println()
	fmt.Println()

	// Row 1: Raw colors (no adjustment)
	fmt.Print("Raw:              ")
	for h := 0; h < 360; h += 10 {
		c := colorful.Hsv(float64(h), saturation, value)
		style := lipgloss.NewStyle().Foreground(c)
		fmt.Print(style.Render(block) + "  ")
	}
	fmt.Println()

	// Final chosen settings: Blue/Violet 25%, Red 15%
	fmt.Print("Adjusted:         ")
	for h := 0; h < 360; h += 10 {
		hf := float64(h)
		satReduce := combinedAdjust(hf, 0.25, 0.15)
		adjSat := math.Max(0.0, saturation-satReduce)
		c := colorful.Hsv(hf, adjSat, value)
		style := lipgloss.NewStyle().Foreground(c)
		fmt.Print(style.Render(block) + "  ")
	}
	fmt.Println()

	fmt.Println()
	fmt.Println("Settings: Blue/Violet -25%, Red -15%")
}
