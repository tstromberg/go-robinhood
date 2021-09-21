package strategy

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tstromberg/roho/pkg/roho"
)

// LuckySevens is a simple strategy to sell stocks at 7.77, and buy at 8.88.
type LuckySevens struct {
	c Config
}

func (cr *LuckySevens) String() string {
	return "Lucky Sevens"
}

func (cr *LuckySevens) Simulate(ctx context.Context, _ time.Time) error {
	return nil
}

func (cr *LuckySevens) Symbols(ctx context.Context) ([]string, error) {
	return []string{"SPY"}, nil
}

func (cr *LuckySevens) Trades(ctx context.Context, symbols []string) ([]Trade, error) {
	r := cr.c.Client
	ts := []Trade{}
	qs, err := r.Quote(ctx, symbols...)
	if err != nil {
		return ts, fmt.Errorf("quote for %q: %w", symbols, err)
	}

	for _, q := range qs {
		ask := fmt.Sprintf("%.2f", q.AskPrice)
		bid := fmt.Sprintf("%.2f", q.BidPrice)
		log.Printf("%q ask=%s bid=%s", q.Symbol, ask, bid)

		switch {
		case ask == "7.77" || ask == "77.70":
			i, err := r.Lookup(ctx, q.Symbol)
			if err != nil {
				return ts, fmt.Errorf("lookup for %q: %w", q.Symbol, err)
			}
			ts = append(ts, Trade{Instrument: i, Order: roho.OrderOpts{Price: q.AskPrice, Quantity: 7, Side: roho.Buy}})
		case bid == "8.88" || ask == "88.80":
			i, err := r.Lookup(ctx, q.Symbol)
			if err != nil {
				return ts, fmt.Errorf("lookup for %q: %w", q.Symbol, err)
			}
			ts = append(ts, Trade{Instrument: i, Order: roho.OrderOpts{Price: q.BidPrice, Quantity: 7, Side: roho.Sell}})
		}
	}

	return ts, nil
}
