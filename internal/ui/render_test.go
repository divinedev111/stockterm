package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/divinedev111/stockterm/internal/quote"
)

func TestRenderQuotes(t *testing.T) {
	quotes := []*quote.Quote{
		{Symbol: "AAPL", Name: "Apple Inc.", Price: 195.50, Change: 2.30, ChangePercent: 1.19, Volume: 55_000_000, Timestamp: time.Now()},
		{Symbol: "TSLA", Name: "Tesla, Inc.", Price: 245.00, Change: -5.00, ChangePercent: -2.00, Volume: 120_000_000, Timestamp: time.Now()},
	}

	result := RenderQuotes(quotes)

	if !strings.Contains(result, "AAPL") {
		t.Error("should contain AAPL")
	}
	if !strings.Contains(result, "TSLA") {
		t.Error("should contain TSLA")
	}
	if !strings.Contains(result, "stockterm") {
		t.Error("should contain header")
	}
}

func TestRenderDetail(t *testing.T) {
	q := &quote.Quote{
		Symbol:        "AAPL",
		Name:          "Apple Inc.",
		Price:         195.50,
		Change:        2.30,
		ChangePercent: 1.19,
		Open:          193.20,
		High:          196.00,
		Low:           192.50,
		Volume:        55_000_000,
		MarketCap:     3_000_000_000_000,
		Timestamp:     time.Now(),
	}

	result := RenderDetail(q, "")

	if !strings.Contains(result, "AAPL") {
		t.Error("should contain symbol")
	}
	if !strings.Contains(result, "$195.50") {
		t.Error("should contain price")
	}
	if !strings.Contains(result, "3.0T") {
		t.Error("should contain market cap")
	}
}

func TestFormatVolume(t *testing.T) {
	tests := []struct {
		input int64
		want  string
	}{
		{500, "500"},
		{5_000, "5.0K"},
		{5_500_000, "5.5M"},
		{2_300_000_000, "2.3B"},
		{1_500_000_000_000, "1.5T"},
	}

	for _, tt := range tests {
		got := formatVolume(tt.input)
		if got != tt.want {
			t.Errorf("formatVolume(%d) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
