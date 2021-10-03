package strategy

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tstromberg/roho/pkg/roho"
)

func TestHiLo(t *testing.T) {
	s, err := New(Config{Kind: HiLo})
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	cs := []*CombinedStock{
		{
			Instrument:   &roho.Instrument{Symbol: "sell-exact"},
			Position:     &roho.Position{Quantity: 3, AverageBuyPrice: 6.00},
			Quote:        &roho.Quote{BidPrice: 8.88},
			Fundamentals: &roho.Fundamental{High52Weeks: 8.88},
		},
		{
			Instrument:   &roho.Instrument{Symbol: "sell-close"},
			Position:     &roho.Position{Quantity: 2, AverageBuyPrice: 7.00},
			Quote:        &roho.Quote{BidPrice: 8.82},
			Fundamentals: &roho.Fundamental{High52Weeks: 8.88},
		},
		{
			Instrument:   &roho.Instrument{Symbol: "buy-exact"},
			Quote:        &roho.Quote{AskPrice: 8.88},
			Fundamentals: &roho.Fundamental{Low52Weeks: 8.88},
		},
		{
			Instrument:   &roho.Instrument{Symbol: "buy-close"},
			Quote:        &roho.Quote{AskPrice: 8.94},
			Fundamentals: &roho.Fundamental{Low52Weeks: 8.88},
		},
		{
			Instrument:   &roho.Instrument{Symbol: "ignore-far"},
			Quote:        &roho.Quote{BidPrice: 9.15, AskPrice: 9.15},
			Fundamentals: &roho.Fundamental{Low52Weeks: 8.88, High52Weeks: 9.99},
		},
		{
			Instrument:   &roho.Instrument{Symbol: "ignore-lower-than-paid"},
			Quote:        &roho.Quote{BidPrice: 9.99, AskPrice: 9.99},
			Position:     &roho.Position{Quantity: 10, AverageBuyPrice: 100.00},
			Fundamentals: &roho.Fundamental{Low52Weeks: 8.88, High52Weeks: 9.99},
		},
	}

	want := []Trade{
		{Instrument: &roho.Instrument{Symbol: "sell-exact"}, Order: roho.OrderOpts{Price: 8.88, Quantity: 3, Side: roho.Sell}, Reason: "0.0% away from 52-week high of 8.88"},
		{Instrument: &roho.Instrument{Symbol: "sell-close"}, Order: roho.OrderOpts{Price: 8.82, Quantity: 2, Side: roho.Sell}, Reason: "0.7% away from 52-week high of 8.88"},
		{Instrument: &roho.Instrument{Symbol: "buy-exact"}, Order: roho.OrderOpts{Price: 8.88, Quantity: 1, Side: roho.Buy}, Reason: "0.0% away from 52wk low of 8.88"},
		{Instrument: &roho.Instrument{Symbol: "buy-close"}, Order: roho.OrderOpts{Price: 8.94, Quantity: 1, Side: roho.Buy}, Reason: "0.7% away from 52wk low of 8.88"},
	}
	got, err := s.Trades(context.Background(), cs)
	if err != nil {
		t.Errorf("Trades() returned unexpected error: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unexpected trade diff: %s", diff)
	}
}

func TestPercentDiff(t *testing.T) {
	tests := []struct {
		a    float64
		b    float64
		want float64
	}{
		{100, 80, -20},
		{4, 5, 25},
	}

	for _, tc := range tests {
		got := percentDiff(tc.a, tc.b)
		if fmt.Sprintf("%.3f", got) != fmt.Sprintf("%.3f", tc.want) {
			t.Errorf("percentDiff(%v, %v) = %.3f, want %.3f", tc.a, tc.b, got, tc.want)
		}
	}
}
