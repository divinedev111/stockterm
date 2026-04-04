package chart

import (
	"strings"
	"testing"
	"time"

	"github.com/divinedev111/stockterm/internal/quote"
)

func TestRenderEmpty(t *testing.T) {
	result := Render(nil, 50, 15)
	if result != "" {
		t.Error("empty candles should return empty string")
	}
}

func TestRenderTooSmall(t *testing.T) {
	candles := []quote.OHLC{{Close: 100}}
	if Render(candles, 5, 3) != "" {
		t.Error("too small dimensions should return empty")
	}
}

func TestRenderBasic(t *testing.T) {
	now := time.Now()
	candles := make([]quote.OHLC, 20)
	for i := range candles {
		candles[i] = quote.OHLC{
			Time:  now.Add(time.Duration(i) * time.Hour),
			Close: 100 + float64(i),
		}
	}

	result := Render(candles, 50, 10)
	if result == "" {
		t.Fatal("expected non-empty chart")
	}
	if !strings.Contains(result, "█") {
		t.Error("chart should contain block characters")
	}
	if !strings.Contains(result, "│") {
		t.Error("chart should contain Y-axis")
	}
}

func TestResample(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result := resample(data, 5)
	if len(result) != 5 {
		t.Errorf("len = %d, want 5", len(result))
	}
}

func TestMinMax(t *testing.T) {
	min, max := minMax([]float64{3, 1, 4, 1, 5, 9})
	if min != 1 {
		t.Errorf("min = %f, want 1", min)
	}
	if max != 9 {
		t.Errorf("max = %f, want 9", max)
	}
}

func TestFormatPrice(t *testing.T) {
	tests := []struct {
		price float64
		want  string
	}{
		{1500, "1500"},
		{42.50, "42.50"},
		{0.0035, "0.0035"},
	}
	for _, tt := range tests {
		got := formatPrice(tt.price)
		if got != tt.want {
			t.Errorf("formatPrice(%f) = %q, want %q", tt.price, got, tt.want)
		}
	}
}
