package roho

import (
	"context"
	"fmt"
)

// Instrument is a type to represent the "instrument" API type in the
// unofficial robinhood API.
type Instrument struct {
	BloombergUnique       string      `json:"bloomberg_unique"`
	Country               string      `json:"country"`
	DayTradeRatio         string      `json:"day_trade_ratio"`
	DefaultCollarFraction string      `json:"default_collar_fraction"`
	FractionalTradability string      `json:"fractional_tradability"`
	Fundamentals          string      `json:"fundamentals"`
	ID                    string      `json:"id"`
	ListDate              string      `json:"list_date"`
	MaintenanceRatio      string      `json:"maintenance_ratio"`
	MarginInitialRatio    string      `json:"margin_initial_ratio"`
	Market                string      `json:"market"`
	MinTickSize           interface{} `json:"min_tick_size"`
	Name                  string      `json:"name"`
	Quote                 string      `json:"quote"`
	RHSTRadability        string      `json:"rhs_tradability"`
	SimpleName            interface{} `json:"simple_name"`
	Splits                string      `json:"splits"`
	State                 string      `json:"state"`
	Symbol                string      `json:"symbol"`
	Tradeable             bool        `json:"tradeable"`
	Tradability           string      `json:"tradability"`
	TradableChainID       string      `json:"tradable_chain_id"`
	Type                  string      `json:"type"`
	URL                   string      `json:"url"`
}

// Instruments returns an Instrument for a single stock symbol.
func (c *Client) Instrument(ctx context.Context, symbol string) (Instrument, error) {
	var i struct {
		Results []Instrument
	}

	url := fmt.Sprintf("%s?symbol=%s", baseURL("instruments"), symbol)

	err := c.get(ctx, url, &i)
	if err != nil {
		return Instrument{}, err
	}
	if len(i.Results) < 1 {
		return Instrument{}, fmt.Errorf("no results")
	}
	return i.Results[0], err
}

// Instrument returns Instruments for a set of stock symbols.
func (c *Client) Instruments(ctx context.Context, syms []string) ([]Instrument, error) {
	is := []Instrument{}
	// Unlike quotes, RH has no native way to query for multiple symbols :(
	for _, s := range syms {
		i, err := c.Instrument(ctx, s)
		if err != nil {
			return is, fmt.Errorf("instrument %q: %w", s, err)
		}
		is = append(is, i)
	}
	return is, nil
}

// Instrument returns an Instrument given a URL.
func (c *Client) InstrumentFromURL(ctx context.Context, url string) (Instrument, error) {
	var i Instrument

	if err := c.get(ctx, url, &i); err != nil {
		return Instrument{}, err
	}

	return i, nil
}
