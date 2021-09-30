package strategy

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tstromberg/roho/pkg/roho"
)

func TestLuckySevens(t *testing.T) {
	s, err := New(Config{Kind: LuckySevens})
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	cs := []*CombinedStock{
		{Instrument: &roho.Instrument{URL: "s1"}, Position: &roho.Position{Quantity: 3}, Quote: &roho.Quote{BidPrice: 8.88}},
		{Instrument: &roho.Instrument{URL: "s2"}, Position: &roho.Position{Quantity: 2}, Quote: &roho.Quote{BidPrice: 88.80}},
		{Instrument: &roho.Instrument{URL: "b1"}, Quote: &roho.Quote{AskPrice: 7.77}},
		{Instrument: &roho.Instrument{URL: "b2"}, Quote: &roho.Quote{AskPrice: 77.70}},
		{Instrument: &roho.Instrument{URL: "ignore"}, Quote: &roho.Quote{AskPrice: 44.40}},
	}

	want := []Trade{
		{Instrument: &roho.Instrument{URL: "s1"}, Order: roho.OrderOpts{Price: 8.88, Quantity: 3, Side: roho.Sell}},
		{Instrument: &roho.Instrument{URL: "s2"}, Order: roho.OrderOpts{Price: 88.80, Quantity: 2, Side: roho.Sell}},
		{Instrument: &roho.Instrument{URL: "b1"}, Order: roho.OrderOpts{Price: 7.77, Quantity: 100, Side: roho.Buy}},
		{Instrument: &roho.Instrument{URL: "b2"}, Order: roho.OrderOpts{Price: 77.70, Quantity: 10, Side: roho.Buy}},
	}
	got, err := s.Trades(context.Background(), cs)
	if err != nil {
		t.Errorf("Trades() returned unexpected error: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unexpected trade diff: %s", diff)
	}
}
