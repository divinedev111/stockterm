// Package api provides financial data from Yahoo Finance.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/divinedev111/stockterm/internal/quote"
)

const baseURL = "https://query1.finance.yahoo.com"

// Client fetches quotes and chart data from Yahoo Finance.
type Client struct {
	http *http.Client
}

// NewClient creates a new Yahoo Finance client.
func NewClient() *Client {
	return &Client{
		http: &http.Client{Timeout: 10 * time.Second},
	}
}

// GetQuote fetches a real-time quote for the given symbol.
func (c *Client) GetQuote(ctx context.Context, symbol string) (*quote.Quote, error) {
	u := fmt.Sprintf("%s/v7/finance/quote?symbols=%s", baseURL, url.QueryEscape(symbol))
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "stockterm/0.1")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch quote: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("yahoo API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		QuoteResponse struct {
			Result []struct {
				Symbol                     string  `json:"symbol"`
				ShortName                  string  `json:"shortName"`
				RegularMarketPrice         float64 `json:"regularMarketPrice"`
				RegularMarketChange        float64 `json:"regularMarketChange"`
				RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
				RegularMarketOpen          float64 `json:"regularMarketOpen"`
				RegularMarketDayHigh       float64 `json:"regularMarketDayHigh"`
				RegularMarketDayLow        float64 `json:"regularMarketDayLow"`
				RegularMarketVolume        int64   `json:"regularMarketVolume"`
				MarketCap                  int64   `json:"marketCap"`
				RegularMarketTime          int64   `json:"regularMarketTime"`
			} `json:"result"`
			Error *struct {
				Code string `json:"code"`
			} `json:"error"`
		} `json:"quoteResponse"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode quote: %w", err)
	}

	if result.QuoteResponse.Error != nil {
		return nil, fmt.Errorf("quote error: %s", result.QuoteResponse.Error.Code)
	}

	if len(result.QuoteResponse.Result) == 0 {
		return nil, fmt.Errorf("no results for symbol %q", symbol)
	}

	r := result.QuoteResponse.Result[0]
	return &quote.Quote{
		Symbol:        r.Symbol,
		Name:          r.ShortName,
		Price:         r.RegularMarketPrice,
		Change:        r.RegularMarketChange,
		ChangePercent: r.RegularMarketChangePercent,
		Open:          r.RegularMarketOpen,
		High:          r.RegularMarketDayHigh,
		Low:           r.RegularMarketDayLow,
		Volume:        r.RegularMarketVolume,
		MarketCap:     r.MarketCap,
		Timestamp:     time.Unix(r.RegularMarketTime, 0),
	}, nil
}

// GetChart fetches historical OHLC data for charting.
func (c *Client) GetChart(ctx context.Context, symbol, interval, timeRange string) ([]quote.OHLC, error) {
	u := fmt.Sprintf("%s/v8/finance/chart/%s?interval=%s&range=%s",
		baseURL, url.QueryEscape(symbol), interval, timeRange)
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "stockterm/0.1")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch chart: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("chart API error (status %d)", resp.StatusCode)
	}

	var result struct {
		Chart struct {
			Result []struct {
				Timestamp  []int64 `json:"timestamp"`
				Indicators struct {
					Quote []struct {
						Open   []float64 `json:"open"`
						High   []float64 `json:"high"`
						Low    []float64 `json:"low"`
						Close  []float64 `json:"close"`
						Volume []int64   `json:"volume"`
					} `json:"quote"`
				} `json:"indicators"`
			} `json:"result"`
		} `json:"chart"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode chart: %w", err)
	}

	if len(result.Chart.Result) == 0 || len(result.Chart.Result[0].Timestamp) == 0 {
		return nil, fmt.Errorf("no chart data for %s", symbol)
	}

	r := result.Chart.Result[0]
	q := r.Indicators.Quote[0]
	n := len(r.Timestamp)

	candles := make([]quote.OHLC, 0, n)
	for i := 0; i < n; i++ {
		if i >= len(q.Close) || q.Close[i] == 0 {
			continue
		}
		candles = append(candles, quote.OHLC{
			Time:   time.Unix(r.Timestamp[i], 0),
			Open:   q.Open[i],
			High:   q.High[i],
			Low:    q.Low[i],
			Close:  q.Close[i],
			Volume: safeVolume(q.Volume, i),
		})
	}

	return candles, nil
}

func safeVolume(volumes []int64, i int) int64 {
	if i < len(volumes) {
		return volumes[i]
	}
	return 0
}
