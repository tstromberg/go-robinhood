package strategy

import (
	"context"
	"fmt"
	"math"
	"regexp"

	"github.com/tstromberg/roho/pkg/roho"
	"k8s.io/klog/v2"
)

var (
	sevenRe = regexp.MustCompile(`^[7\.]+[70]$`)
	eightRe = regexp.MustCompile(`^[8\.]+[80]$`)
)

// LuckySevensStrategy is a demonstration strategy to buy stocks at 7.77 and sell them at 8.88.
type LuckySevensStrategy struct {
	c Config
}

func (cr *LuckySevensStrategy) String() string {
	return "Lucky Sevens"
}

func (cr *LuckySevensStrategy) Trades(_ context.Context, cs []*CombinedStock) ([]Trade, error) {
	ts := []Trade{}

	for _, s := range cs {
		p := s.Position
		if p != nil {
			bid := fmt.Sprintf("%.2f", s.Quote.BidPrice)
			klog.Infof("%q bid=%q", s.Instrument.Symbol, bid)
			if eightRe.MatchString(bid) {
				ts = append(ts, Trade{Instrument: s.Instrument, Order: roho.OrderOpts{Price: s.Quote.BidPrice, Quantity: uint64(p.Quantity), Side: roho.Sell}})
				continue
			}
		}

		ask := fmt.Sprintf("%.2f", s.Quote.AskPrice)
		klog.Infof("%q ask=%q", s.Instrument.Symbol, ask)
		if sevenRe.MatchString(ask) {
			q := uint64(math.Round(777.77 / s.Quote.AskPrice))
			ts = append(ts, Trade{Instrument: s.Instrument, Order: roho.OrderOpts{Price: s.Quote.AskPrice, Quantity: q, Side: roho.Buy}})
		}
	}

	return ts, nil
}
