// Package quote represents stock/crypto price data.
package quote

import "time"

// Quote holds a single price snapshot.
type Quote struct {
	Symbol        string
	Name          string
	Price         float64
	Change        float64
	ChangePercent float64
	Open          float64
	High          float64
	Low           float64
	Volume        int64
	MarketCap     int64
	Timestamp     time.Time
}

// IsUp returns true if the quote shows positive change.
func (q *Quote) IsUp() bool { return q.Change >= 0 }

// OHLC holds a single candlestick.
type OHLC struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}
