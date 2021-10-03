package roho

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Historical struct {
	Symbol   string             `json:"symbol"`
	Interval string             `json:"interval"`
	Bounds   string             `json:"bounds"`
	Span     string             `json:"span"`
	Records  []HistoricalRecord `json:"historicals"`
}

type HistoricalRecord struct {
	BeginsAt     time.Time `json:"begins_at"`
	OpenPrice    float64   `json:"open_price,string"`
	ClosePrice   float64   `json:"close_price,string"`
	HighPrice    float64   `json:"high_price,string"`
	LowPrice     float64   `json:"low_price,string"`
	Volume       int64     `json:"volume"`
	Session      string    `json:"session"`
	Interpolated bool      `json:"interpolated"`
}

const (
	// Valid interval values.
	FiveMinute   = "5minute"
	TenMinute    = "10minute"
	ThirtyMinute = "30minute"
	Hour         = "hour"

	// Valid for both intervals and spans.
	Day  = "day"
	Week = "week"

	// Valid span values.
	Month    = "month"
	Year     = "year"
	FiveYear = "5year"
)

// Historicals returns historical data for the list of stocks provided. See the interval/span constants for hints.
func (c *Client) Historicals(ctx context.Context, interval string, span string, symbols []string) ([]Historical, error) {
	url := fmt.Sprintf("%s?interval=%s&span=%s&symbols=%s", baseURL("quotes/historicals"), interval, span, strings.Join(symbols, ","))
	var r struct{ Results []Historical }
	err := c.get(ctx, url, &r)
	return r.Results, err
}

// Historicals returns historical data for the list of stocks provided. See the interval/span constants for hints.
func (c *Client) Historical(ctx context.Context, interval string, span string, symbol string) (Historical, error) {
	hs, err := c.Historicals(ctx, interval, span, []string{symbol})
	if err != nil {
		return Historical{}, err
	}
	return hs[0], nil
}
