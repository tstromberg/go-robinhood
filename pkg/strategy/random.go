package strategy

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/tstromberg/roho/pkg/roho"
)

// RandomStrategy is a demonstration strategy to buy/sell stocks at random.
type RandomStrategy struct {
	c Config
}

func (cr *RandomStrategy) String() string {
	return "Random"
}

func (cr *RandomStrategy) Trades(_ context.Context, cs []*CombinedStock) ([]Trade, error) {
	luckyNumber, ok := cr.c.Values["lucky-number"]
	if !ok {
		luckyNumber = int64(4)
	}

	maxRand := int64(len(cs)) * luckyNumber
	ts := []Trade{}

	// Sell first
	for _, s := range cs {
		if s.Position == nil {
			continue
		}

		nb, err := rand.Int(rand.Reader, big.NewInt(maxRand))
		if err != nil {
			return ts, fmt.Errorf("rand int: %w", err)
		}
		if nb.Int64() != luckyNumber {
			continue
		}
		ts = append(ts, Trade{Instrument: s.Instrument, Order: roho.OrderOpts{Price: s.Quote.BidPrice, Quantity: uint64(s.Position.Quantity), Side: roho.Sell}})
	}

	// Now buy
	for _, s := range cs {
		nb, err := rand.Int(rand.Reader, big.NewInt(maxRand))
		if err != nil {
			return ts, fmt.Errorf("rand int: %w", err)
		}
		if nb.Int64() != luckyNumber {
			continue
		}
		ts = append(ts, Trade{Instrument: s.Instrument, Order: roho.OrderOpts{Price: s.Quote.AskPrice, Quantity: uint64(luckyNumber), Side: roho.Buy}})
	}

	return ts, nil
}
