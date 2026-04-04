// Package chart renders ASCII price charts in the terminal.
package chart

import (
	"fmt"
	"math"
	"strings"

	"github.com/divinedev111/stockterm/internal/quote"
)

const (
	green = "\033[32m"
	red   = "\033[31m"
	dim   = "\033[2m"
	reset = "\033[0m"
)

// Render draws an ASCII chart from OHLC data.
// width and height are in terminal characters.
func Render(candles []quote.OHLC, width, height int) string {
	if len(candles) == 0 || width < 10 || height < 5 {
		return ""
	}

	// Use closing prices for the line chart
	prices := make([]float64, len(candles))
	for i, c := range candles {
		prices[i] = c.Close
	}

	// Resample to fit width
	if len(prices) > width {
		prices = resample(prices, width)
	}

	minPrice, maxPrice := minMax(prices)
	priceRange := maxPrice - minPrice
	if priceRange == 0 {
		priceRange = 1
	}

	// Build the chart grid
	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, len(prices))
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	// Plot prices
	for col, price := range prices {
		row := height - 1 - int(math.Round(float64(height-1)*(price-minPrice)/priceRange))
		if row < 0 {
			row = 0
		}
		if row >= height {
			row = height - 1
		}
		grid[row][col] = '█'

		// Fill below for area chart effect
		for r := row + 1; r < height; r++ {
			grid[r][col] = '░'
		}
	}

	// Determine if overall trend is up or down
	isUp := len(prices) >= 2 && prices[len(prices)-1] >= prices[0]
	lineColor := red
	if isUp {
		lineColor = green
	}

	// Render to string with Y-axis labels
	var b strings.Builder
	labelWidth := priceLabel(maxPrice)
	for row := 0; row < height; row++ {
		// Y-axis price label (top, middle, bottom)
		price := maxPrice - (float64(row)/float64(height-1))*priceRange
		switch {
		case row == 0:
			b.WriteString(dim + fmt.Sprintf("%*s", labelWidth, formatPrice(maxPrice)) + " │" + reset)
		case row == height/2:
			mid := (maxPrice + minPrice) / 2
			b.WriteString(dim + fmt.Sprintf("%*s", labelWidth, formatPrice(mid)) + " │" + reset)
		case row == height-1:
			b.WriteString(dim + fmt.Sprintf("%*s", labelWidth, formatPrice(minPrice)) + " │" + reset)
		default:
			_ = price
			b.WriteString(strings.Repeat(" ", labelWidth) + dim + " │" + reset)
		}

		// Chart row
		b.WriteString(lineColor)
		for _, ch := range grid[row] {
			b.WriteRune(ch)
		}
		b.WriteString(reset + "\n")
	}

	// X-axis
	b.WriteString(strings.Repeat(" ", labelWidth) + dim + " └" + strings.Repeat("─", len(prices)) + reset + "\n")

	// Time labels
	if len(candles) > 0 {
		first := candles[0].Time.Format("Jan 2")
		last := candles[len(candles)-1].Time.Format("Jan 2")
		gap := len(prices) - len(first) - len(last)
		if gap > 0 {
			b.WriteString(strings.Repeat(" ", labelWidth+2) + dim + first + strings.Repeat(" ", gap) + last + reset + "\n")
		}
	}

	return b.String()
}

func resample(data []float64, target int) []float64 {
	result := make([]float64, target)
	ratio := float64(len(data)) / float64(target)
	for i := 0; i < target; i++ {
		idx := int(float64(i) * ratio)
		if idx >= len(data) {
			idx = len(data) - 1
		}
		result[i] = data[idx]
	}
	return result
}

func minMax(data []float64) (float64, float64) {
	min, max := data[0], data[0]
	for _, v := range data[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}

func formatPrice(p float64) string {
	if p >= 1000 {
		return fmt.Sprintf("%.0f", p)
	}
	if p >= 1 {
		return fmt.Sprintf("%.2f", p)
	}
	return fmt.Sprintf("%.4f", p)
}

func priceLabel(maxPrice float64) int {
	return len(formatPrice(maxPrice)) + 1
}
