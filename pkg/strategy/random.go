package strategy

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/tstromberg/roho/pkg/roho"
)

var myLuckyNumber = int64(4)

// Randob is a simple strategy to buy stocks at 7.77 and sell them at 8.88
type RandomStrategy struct {
	c Config
}

func (cr *RandomStrategy) String() string {
	return "Random"
}

func (cr *RandomStrategy) Trades(ctx context.Context, cs map[string]*CombinedStock) ([]Trade, error) {
	maxRand := int64(len(cs)) * myLuckyNumber
	ts := []Trade{}

	for sym, c := range cs {
		if c.Position == nil {
			continue
		}

		nb, err := rand.Int(rand.Reader, big.NewInt(maxRand))
		if err != nil {
			return ts, fmt.Errorf("rand int: %w", err)
		}
		if nb.Int64() != myLuckyNumber {
			continue
		}
		ts = append(ts, Trade{Symbol: sym, Order: roho.OrderOpts{Price: c.Position.AverageBuyPrice, Quantity: uint64(c.Position.Quantity), Side: roho.Sell}})

	}

	for sym, c := range cs {
		nb, err := rand.Int(rand.Reader, big.NewInt(maxRand))
		if err != nil {
			return ts, fmt.Errorf("rand int: %w", err)
		}
		if nb.Int64() != myLuckyNumber {
			continue
		}
		ts = append(ts, Trade{Symbol: sym, Order: roho.OrderOpts{Price: c.Quote.AskPrice, Quantity: uint64(myLuckyNumber), Side: roho.Buy}})
	}

	return ts, nil
}

func (cr *RandomStrategy) SetTime(ctx context.Context, _ time.Time) error {
	return nil
}
