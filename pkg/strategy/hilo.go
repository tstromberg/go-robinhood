package strategy

import (
	"context"
	"fmt"

	"github.com/tstromberg/roho/pkg/roho"
	"k8s.io/klog/v2"
)

// HiLoStrategy is a demonstration strategy to buy/sell stocks when near their 52-week hi/low averages.
type HiLoStrategy struct {
	c Config
}

func (cr *HiLoStrategy) String() string {
	return "HiLo"
}

func (cr *HiLoStrategy) Trades(_ context.Context, cs []*CombinedStock) ([]Trade, error) {
	ts := []Trade{}

	for _, s := range cs {
		p := s.Position

		// Buy stock only if we do not yet own it
		if p == nil {
			perc := percentDiff(s.Fundamentals.Low52Weeks, s.Quote.AskPrice)
			if perc < 2 {
				klog.Infof("%s: ask price of %.2f is %.2f%% away from 52-week low of %.2f", s.Instrument.Symbol, s.Quote.AskPrice, perc, s.Fundamentals.Low52Weeks)
			}

			if perc < 0 {
				klog.Warningf("%s buy perc=%.2f - ask price is below 52 weeks?", s.Instrument.Symbol, perc)
				continue
			}

			if perc <= 0.9 {
				ts = append(ts, Trade{
					Instrument: s.Instrument,
					Order:      roho.OrderOpts{Price: s.Quote.AskPrice, Quantity: 1, Side: roho.Buy},
					Reason:     fmt.Sprintf("%.1f%% away from 52wk low of %.2f", perc, s.Fundamentals.Low52Weeks),
				})
			}
			continue
		}

		perc := percentDiff(s.Quote.BidPrice, s.Fundamentals.High52Weeks)

		if perc < 2 {
			klog.Infof("%s: bid price of %.2f is %.2f%% away from 52-week high of %.2f", s.Instrument.Symbol, s.Quote.BidPrice, perc, s.Fundamentals.High52Weeks)
		}

		if perc < 0 {
			klog.Warningf("%s sell perc=%.2f - bid price is above 52 weeks?", s.Instrument.Symbol, perc)
			continue
		}

		if perc <= 0.9 {
			// Only sell if we make a profit
			if p.AverageBuyPrice > s.Quote.BidPrice {
				klog.Infof("would sell %s for %.2f but we paid %.2f for it", s.Instrument.Symbol, s.Quote.BidPrice, p.AverageBuyPrice)
				continue
			}

			ts = append(ts, Trade{
				Instrument: s.Instrument,
				Order:      roho.OrderOpts{Price: s.Quote.BidPrice, Quantity: uint64(p.Quantity), Side: roho.Sell},
				Reason:     fmt.Sprintf("%.1f%% away from 52-week high of %.2f", perc, s.Fundamentals.High52Weeks),
			})
			continue
		}
	}

	return ts, nil
}

func percentDiff(old, n float64) float64 {
	return ((n - old) / old) * 100
}
