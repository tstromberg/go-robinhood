package strategy

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/tstromberg/roho/pkg/roho"
)

var (
	sevenRe = regexp.MustCompile(`^[7\.]+[70]$`)
	eightRe = regexp.MustCompile(`^[8\.]+[80]$`)
)

// LuckySevens is a simple strategy to buy stocks at 7.77 and sell them at 8.88
type LuckySevensStrategy struct {
	c Config
}

func (cr *LuckySevensStrategy) String() string {
	return "Lucky Sevens"
}

func (cr *LuckySevensStrategy) Trades(ctx context.Context, ps []roho.Position, qs []roho.Quote) ([]Trade, error) {
	ts := []Trade{}

	for _, p := range ps {
		bid := fmt.Sprintf("%.2f", p.AverageBuyPrice)
		if eightRe.MatchString(bid) {
			ts = append(ts, Trade{InstrumentURL: p.Instrument, Order: roho.OrderOpts{Price: p.AverageBuyPrice, Quantity: uint64(p.Quantity), Side: roho.Sell}})
		}
	}

	for _, q := range qs {
		ask := fmt.Sprintf("%.2f", q.AskPrice)
		if sevenRe.MatchString(ask) {
			ts = append(ts, Trade{Symbol: q.Symbol, Order: roho.OrderOpts{Price: q.AskPrice, Quantity: 7, Side: roho.Buy}})
		}
	}

	return ts, nil
}

func (cr *LuckySevensStrategy) SetTime(ctx context.Context, _ time.Time) error {
	return nil
}
