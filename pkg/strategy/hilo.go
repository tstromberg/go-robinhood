package strategy

import (
	"context"
	"fmt"
	"log"

	"github.com/tstromberg/roho/pkg/roho"
)

// HiLoStrategy is a demonstration strategy to buy/sell stocks when near their 52-week hi/low averages.
type HiLoStrategy struct {
	c Config
}

func (cr *HiLoStrategy) String() string {
	return "HiLo"
}

func (cr *HiLoStrategy) Trades(_ context.Context, cs map[string]*CombinedStock) ([]Trade, error) {
	ts := []Trade{}

	for url, s := range cs {
		p := s.Position

		// We don't already own this stock
		if p == nil {
			ratio := s.Quote.BidPrice / s.Fundamentals.Low52Weeks
			if ratio < 1.01 {
				ts = append(ts, Trade{
					InstrumentURL: url,
					Symbol:        s.Quote.Symbol,
					Order:         roho.OrderOpts{Price: s.Quote.AskPrice, Quantity: 1, Side: roho.Buy},
					Reason:        fmt.Sprintf("Near 52-week low of %.2f", s.Fundamentals.Low52Weeks),
				})
			}
			continue
		}

		ratio := s.Quote.BidPrice / s.Fundamentals.High52Weeks
		log.Printf("%s: ratio between buy %.2f and high %.2f is %.2f", s.Quote.Symbol, s.Quote.BidPrice, s.Fundamentals.High52Weeks, ratio)
		if ratio > 0.99 {
			// Only sell if we make a profit
			if p.AverageBuyPrice < s.Quote.BidPrice {
				log.Printf("would sell %s for %.2f but we paid %.2f for it", s.Quote.Symbol, s.Quote.BidPrice, p.AverageBuyPrice)
				continue
			}

			ts = append(ts, Trade{
				InstrumentURL: url,
				Symbol:        s.Quote.Symbol,
				Order:         roho.OrderOpts{Price: p.AverageBuyPrice, Quantity: uint64(p.Quantity), Side: roho.Sell},
				Reason:        fmt.Sprintf("Near 52-week high of %.2f", s.Fundamentals.High52Weeks),
			})
			continue
		}
	}

	return ts, nil
}
