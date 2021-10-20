package strategy

import (
	"context"
	"fmt"
	"time"

	"github.com/tstromberg/roho/pkg/roho"
	"k8s.io/klog/v2"
)

// BounceStrategy is a demonstration strategy to buy/sell stocks when they change direction near their 52-week hi/low.
type BounceStrategy struct {
	c Config
}

func (cr *BounceStrategy) String() string {
	return "Bounce"
}

func (cr *BounceStrategy) Trades(ctx context.Context, cs []*CombinedStock) ([]Trade, error) {
	ts := []Trade{}

	for _, s := range cs {
		if s.Position == nil {
			t := cr.determineBuy(ctx, s)
			if t != nil {
				ts = append(ts, *t)
			}
			continue
		}

		t := cr.determineSell(ctx, s)
		if t != nil {
			ts = append(ts, *t)
		}
	}

	return ts, nil
}

func upward(fs []float64) (bool, float64) {
	upward := true

	for i, f := range fs {
		if i == 0 {
			continue
		}
		if f < fs[0] {
			upward = false
		}
	}

	//	klog.Infof("upward values=%v, trend is %.2f", fs, percentDiff(fs[0], fs[len(fs)-1]))
	return upward, percentDiff(fs[0], fs[len(fs)-1])
}

const (
	minTrend       = 8
	bounceNearness = 1.2
	tooFarOff      = -10
)

func downward(fs []float64) (bool, float64) {
	downward := true

	for i, f := range fs {
		if i == 0 {
			continue
		}
		if f > fs[0] {
			downward = false
		}
	}

	//	klog.Infof("downward values=%v, trend is %.2f", fs, percentDiff(fs[0], fs[len(fs)-1]))
	return downward, percentDiff(fs[0], fs[len(fs)-1])
}

func priceTrend(rs []roho.HistoricalRecord) []float64 {
	prices := []float64{}
	for _, r := range rs {
		prices = append(prices, r.OpenPrice)
		prices = append(prices, r.ClosePrice)
	}
	return prices
}

func (cr *BounceStrategy) determineBuy(ctx context.Context, s *CombinedStock) *Trade {
	perc := percentDiff(s.Fundamentals.Low52Weeks, s.Quote.AskPrice)
	if perc > bounceNearness {
		return nil
	}

	klog.Infof("%s: ask price of %.2f is %.2f%% away from 52-week low of %.2f", s.Instrument.Symbol, s.Quote.AskPrice, perc, s.Fundamentals.Low52Weeks)

	if perc < tooFarOff {
		klog.Warningf("%s: %.2f is far off of historical averages; something fishy is going on: %+v", s.Instrument.Symbol, s.Fundamentals)
		return nil
	}

	var hs roho.Historical
	var err error
	if s.Historical != nil {
		hs = *s.Historical
	} else {
		hs, err = cr.c.Client.Historical(ctx, "5minute", "day", s.Instrument.Symbol)
		if err != nil {
			klog.Errorf("get historicals failed: %v", err)
			return nil
		}
	}

	prices := append(priceTrend(hs.Records), s.Quote.AskPrice)

	if len(prices) < minTrend {
		klog.Warningf("%s: not enough historical data: %v", s.Instrument.Symbol, prices)
		return nil
	}

	recent := lastFloats(prices, minTrend)

	ok, bounce := upward(recent)
	if ok && bounce > 0.01 {
		klog.Infof("%s: buy now: upward=%v, bounce=%.2f: %v", s.Instrument.Symbol, ok, bounce, recent)
		return &Trade{
			Instrument: s.Instrument,
			Order:      roho.OrderOpts{Price: s.Quote.AskPrice, Quantity: 1, Side: roho.Buy},
			Reason:     fmt.Sprintf("%.1f%% away from 52wk low of %.2f, %.2f%% bounce", perc, s.Fundamentals.Low52Weeks, bounce),
		}
	}

	klog.Infof("%s: wait to buy: upward=%v, bounce=%.2f: %v", s.Instrument.Symbol, ok, bounce, recent)
	return nil
}

func lastFloats(nums []float64, size int) []float64 {
	if len(nums) <= size {
		return nums
	}
	return nums[len(nums)-size:]
}

func (cr *BounceStrategy) determineSell(ctx context.Context, s *CombinedStock) *Trade {
	p := s.Position
	perc := percentDiff(s.Quote.BidPrice, s.Fundamentals.High52Weeks)

	if perc < (bounceNearness * 1.5) {
		klog.Infof("%s: bid price of %.2f is %.2f%% away from 52-week high of %.2f", s.Instrument.Symbol, s.Quote.BidPrice, perc, s.Fundamentals.High52Weeks)
	}

	if perc < 0 {
		klog.Warningf("%s: %.2f is far off of historical averages; something fishy is going on: %+v", s.Instrument.Symbol, s.Fundamentals)
		return nil
	}

	if perc > bounceNearness {
		return nil
	}

	if p.AverageBuyPrice > s.Quote.BidPrice {
		klog.Infof("would sell %s for %.2f but we paid %.2f for it", s.Instrument.Symbol, s.Quote.BidPrice, p.AverageBuyPrice)
		return nil
	}

	age := time.Since(p.CreatedAt)
	if age < time.Hour*24*365 {
		klog.Infof("would sell %s for %.2f but it's been held for less than a year (%s)", s.Instrument.Symbol, s.Quote.BidPrice, age)
	}

	var hs roho.Historical
	var err error
	if s.Historical != nil {
		hs = *s.Historical
	} else {
		hs, err = cr.c.Client.Historical(ctx, "5minute", "day", s.Instrument.Symbol)
		if err != nil {
			klog.Errorf("get historicals failed: %v", err)
			return nil
		}
	}

	prices := append(priceTrend(hs.Records), s.Quote.BidPrice)
	if len(prices) < minTrend {
		klog.Warningf("%s: not enough historical data: %v", s.Instrument.Symbol, prices)
		return nil
	}

	recent := lastFloats(prices, minTrend)
	ok, bounce := downward(recent)

	if ok && bounce < -0.01 {
		klog.Infof("%s: sell now: upward=%v, bounce=%.2f: %v", s.Instrument.Symbol, ok, bounce, recent)

		return &Trade{
			Instrument: s.Instrument,
			Order:      roho.OrderOpts{Price: s.Quote.BidPrice, Quantity: uint64(p.Quantity), Side: roho.Sell},
			Reason:     fmt.Sprintf("%.1f%% away from 52-week high of %.2f, %.2f%% bounce", perc, s.Fundamentals.High52Weeks, bounce),
		}
	}

	klog.Infof("%s: wait to sell: downward=%v, bounce=%.2f: %v", s.Instrument.Symbol, ok, bounce, recent)
	return nil
}
