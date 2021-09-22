package strategy

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/tstromberg/roho/pkg/roho"
)

var (
	myLuckyNumber = int64(4)
)

// Randob is a simple strategy to buy stocks at 7.77 and sell them at 8.88
type RandomStrategy struct {
	c Config
}

func (cr *RandomStrategy) String() string {
	return "Random"
}

func (cr *RandomStrategy) Trades(ctx context.Context, ps []roho.Position, qs []roho.Quote) ([]Trade, error) {
	maxRand := int64(len(ps)+len(qs)) * myLuckyNumber
	ts := []Trade{}

	for _, p := range ps {
		nb, err := rand.Int(rand.Reader, big.NewInt(maxRand))
		if err != nil {
			return ts, fmt.Errorf("rand int: %w", err)
		}
		if nb.Int64() != myLuckyNumber {
			continue
		}
		ts = append(ts, Trade{InstrumentURL: p.Instrument, Order: roho.OrderOpts{Price: p.AverageBuyPrice, Quantity: uint64(p.Quantity), Side: roho.Sell}})

	}

	for _, q := range qs {
		nb, err := rand.Int(rand.Reader, big.NewInt(maxRand))
		if err != nil {
			return ts, fmt.Errorf("rand int: %w", err)
		}
		if nb.Int64() != myLuckyNumber {
			continue
		}
		ts = append(ts, Trade{Symbol: q.Symbol, Order: roho.OrderOpts{Price: q.AskPrice, Quantity: 7, Side: roho.Buy}})
	}

	return ts, nil
}

func (cr *RandomStrategy) SetTime(ctx context.Context, _ time.Time) error {
	return nil
}
