package roho

import (
	"context"
	"strings"
)

// Fundamental represents the JSON struct returned by the Robinhood fundamentals API.
type Fundamental struct {
	Open          float64 `json:"open,string"`
	High          float64 `json:"high,string"`
	Low           float64 `json:"low,string"`
	Volume        float64 `json:"volume,string"`
	AverageVolume float64 `json:"average_volume,string"`
	High52Weeks   float64 `json:"high_52_weeks,string"`
	DividendYield float64 `json:"dividend_yield,string"`
	Low52Weeks    float64 `json:"low_52_weeks,string"`
	MarketCap     float64 `json:"market_cap,string"`
	PERatio       float64 `json:"pe_ratio,string"`
	Description   string  `json:"description"`
	InstrumentURL string  `json:"instrument"`
}

// Fundamentals returns fundamental data for the list of stock symbols provided.
func (c *Client) Fundamentals(ctx context.Context, syms ...string) ([]Fundamental, error) {
	fs := []Fundamental{}

	// Robinhood Fundamentals API only allows 100 symbols at a time
	for _, ck := range chunkStrings(syms, 100) {
		url := baseURL("fundamentals") + "?symbols=" + strings.Join(ck, ",")
		var r struct{ Results []Fundamental }
		if err := c.get(ctx, url, &r); err != nil {
			return fs, err
		}
		fs = append(fs, r.Results...)
	}

	return fs, nil
}

// chunkStrings divides a slice of strings into chunks of a specified size.
func chunkStrings(input []string, size int) [][]string {
	var chunks [][]string

	for i := 0; i < len(input); i += size {
		end := i + size

		if end > len(input) {
			end = len(input)
		}

		chunks = append(chunks, input[i:end])
	}
	return chunks
}
