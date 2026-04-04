// Package ui renders the terminal dashboard.
package ui

import (
	"fmt"
	"strings"

	"github.com/divinedev111/stockterm/internal/quote"
)

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	dim    = "\033[2m"
	green  = "\033[32m"
	red    = "\033[31m"
	cyan   = "\033[36m"
	yellow = "\033[33m"
	clear  = "\033[2J\033[H"
)

// RenderQuotes formats a list of quotes as a watchlist table.
func RenderQuotes(quotes []*quote.Quote) string {
	var b strings.Builder

	b.WriteString(clear)
	b.WriteString(bold + cyan + "  stockterm" + reset + dim + "  —  press q to quit, a to add, d to delete, c to chart" + reset + "\n")
	b.WriteString(dim + "  " + strings.Repeat("─", 78) + reset + "\n")

	// Header
	b.WriteString(fmt.Sprintf("  %-8s %-24s %10s %10s %8s %12s\n",
		dim+"SYMBOL", "NAME", "PRICE", "CHANGE", "CHG%", "VOLUME"+reset))
	b.WriteString(dim + "  " + strings.Repeat("─", 78) + reset + "\n")

	for _, q := range quotes {
		color := red
		arrow := "▼"
		if q.IsUp() {
			color = green
			arrow = "▲"
		}

		name := q.Name
		if len(name) > 22 {
			name = name[:22]
		}

		change := fmt.Sprintf("%s %.2f", arrow, q.Change)
		pct := fmt.Sprintf("%.2f%%", q.ChangePercent)
		vol := formatVolume(q.Volume)

		b.WriteString(fmt.Sprintf("  "+bold+"%-8s"+reset+" %-24s %s%10s %10s %8s"+reset+" %12s\n",
			q.Symbol, name, color, formatPrice(q.Price), change, pct, vol))
	}

	b.WriteString(dim + "\n  " + strings.Repeat("─", 78) + reset + "\n")

	return b.String()
}

// RenderDetail formats a detailed view for a single quote.
func RenderDetail(q *quote.Quote, chartStr string) string {
	var b strings.Builder

	color := red
	arrow := "▼"
	if q.IsUp() {
		color = green
		arrow = "▲"
	}

	b.WriteString(clear)
	b.WriteString(bold + cyan + "  " + q.Symbol + reset + "  " + dim + q.Name + reset + "\n")
	b.WriteString(dim + "  " + strings.Repeat("─", 60) + reset + "\n\n")

	b.WriteString(fmt.Sprintf("  Price:   %s%s %s\n", color+bold, formatPrice(q.Price), reset))
	b.WriteString(fmt.Sprintf("  Change:  %s%s %.2f (%.2f%%)%s\n", color, arrow, q.Change, q.ChangePercent, reset))
	b.WriteString(fmt.Sprintf("  Open:    %s\n", formatPrice(q.Open)))
	b.WriteString(fmt.Sprintf("  High:    %s\n", formatPrice(q.High)))
	b.WriteString(fmt.Sprintf("  Low:     %s\n", formatPrice(q.Low)))
	b.WriteString(fmt.Sprintf("  Volume:  %s\n", formatVolume(q.Volume)))

	if q.MarketCap > 0 {
		b.WriteString(fmt.Sprintf("  MktCap:  %s\n", formatVolume(q.MarketCap)))
	}

	b.WriteString(fmt.Sprintf("  As of:   %s\n", dim+q.Timestamp.Format("Jan 2, 2006 3:04 PM")+reset))

	if chartStr != "" {
		b.WriteString("\n")
		b.WriteString(chartStr)
	}

	b.WriteString(dim + "\n  Press ESC to go back" + reset + "\n")

	return b.String()
}

func formatPrice(p float64) string {
	if p >= 1000 {
		return fmt.Sprintf("$%.2f", p)
	}
	if p >= 1 {
		return fmt.Sprintf("$%.2f", p)
	}
	return fmt.Sprintf("$%.6f", p)
}

func formatVolume(v int64) string {
	switch {
	case v >= 1_000_000_000_000:
		return fmt.Sprintf("%.1fT", float64(v)/1e12)
	case v >= 1_000_000_000:
		return fmt.Sprintf("%.1fB", float64(v)/1e9)
	case v >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(v)/1e6)
	case v >= 1_000:
		return fmt.Sprintf("%.1fK", float64(v)/1e3)
	default:
		return fmt.Sprintf("%d", v)
	}
}
