package strategy

import (
	"context"
	"fmt"
	"regexp"

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

func (cr *LuckySevensStrategy) Trades(ctx context.Context, cs map[string]*CombinedStock) ([]Trade, error) {
	ts := []Trade{}

	for url, s := range cs {
		p := s.Position
		if p != nil {
			bid := fmt.Sprintf("%.2f", p.AverageBuyPrice)
			if eightRe.MatchString(bid) {
				ts = append(ts, Trade{InstrumentURL: url, Order: roho.OrderOpts{Price: p.AverageBuyPrice, Quantity: uint64(p.Quantity), Side: roho.Sell}})
				continue
			}
		}

		ask := fmt.Sprintf("%.2f", s.Quote.AskPrice)
		if sevenRe.MatchString(ask) {
			ts = append(ts, Trade{Symbol: s.Quote.Symbol, Order: roho.OrderOpts{Price: s.Quote.AskPrice, Quantity: 7, Side: roho.Buy}})
		}
	}

	return ts, nil
}
