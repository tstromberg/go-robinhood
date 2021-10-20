package strategy

import (
	"context"
	"fmt"
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tstromberg/roho/pkg/roho"
)

func TestBounce(t *testing.T) {
	s, err := New(Config{Kind: Bounce})
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	cs := []*CombinedStock{
		{
			Instrument:   &roho.Instrument{Symbol: "hold-stalled-position"},
			Position:     &roho.Position{Quantity: 3, AverageBuyPrice: 6.00},
			Quote:        &roho.Quote{BidPrice: 8.88},
			Fundamentals: &roho.Fundamental{High52Weeks: 8.88},
			Historical: &roho.Historical{
				Records: []roho.HistoricalRecord{
					{OpenPrice: 8.88, ClosePrice: 8.88, LowPrice: 8.88, HighPrice: 8.88},
					{OpenPrice: 8.88, ClosePrice: 8.88, LowPrice: 8.88, HighPrice: 8.88},
					{OpenPrice: 8.88, ClosePrice: 8.88, LowPrice: 8.88, HighPrice: 8.88},
					{OpenPrice: 8.88, ClosePrice: 8.88, LowPrice: 8.88, HighPrice: 8.88},
				},
			},
		},
		{
			Instrument:   &roho.Instrument{Symbol: "sell-downward-position"},
			Position:     &roho.Position{Quantity: 2, AverageBuyPrice: 7.00},
			Quote:        &roho.Quote{BidPrice: 8.84},
			Fundamentals: &roho.Fundamental{High52Weeks: 8.88},
			Historical: &roho.Historical{
				Records: []roho.HistoricalRecord{
					{OpenPrice: 8.88, HighPrice: 8.88, LowPrice: 8.87, ClosePrice: 8.87},
					{OpenPrice: 8.87, HighPrice: 8.87, LowPrice: 8.86, ClosePrice: 8.86},
					{OpenPrice: 8.86, HighPrice: 8.86, LowPrice: 8.85, ClosePrice: 8.85},
					{OpenPrice: 8.85, HighPrice: 8.85, LowPrice: 8.84, ClosePrice: 8.84},
					{OpenPrice: 8.84, HighPrice: 8.84, LowPrice: 8.83, ClosePrice: 8.83},
					{OpenPrice: 8.83, HighPrice: 8.83, LowPrice: 8.82, ClosePrice: 8.82},
				},
			},
		},
		{
			Instrument:   &roho.Instrument{Symbol: "hold-upward-position"},
			Position:     &roho.Position{Quantity: 2, AverageBuyPrice: 7.00},
			Quote:        &roho.Quote{BidPrice: 8.88},
			Fundamentals: &roho.Fundamental{High52Weeks: 8.88},
			Historical: &roho.Historical{
				Records: []roho.HistoricalRecord{
					{OpenPrice: 8.84, LowPrice: 8.84, HighPrice: 8.85, ClosePrice: 8.85},
					{OpenPrice: 8.85, LowPrice: 8.85, HighPrice: 8.86, ClosePrice: 8.86},
					{OpenPrice: 8.86, LowPrice: 8.86, HighPrice: 8.87, ClosePrice: 8.87},
				},
			},
		},
		{
			Instrument:   &roho.Instrument{Symbol: "wait-stalled-option"},
			Quote:        &roho.Quote{AskPrice: 8.88},
			Fundamentals: &roho.Fundamental{Low52Weeks: 8.88},
			Historical: &roho.Historical{
				Records: []roho.HistoricalRecord{
					{OpenPrice: 8.88, ClosePrice: 8.88, LowPrice: 8.88, HighPrice: 8.88},
					{OpenPrice: 8.88, ClosePrice: 8.88, LowPrice: 8.88, HighPrice: 8.88},
					{OpenPrice: 8.88, ClosePrice: 8.88, LowPrice: 8.88, HighPrice: 8.88},
					{OpenPrice: 8.88, ClosePrice: 8.88, LowPrice: 8.88, HighPrice: 8.88},
					{OpenPrice: 8.88, ClosePrice: 8.88, LowPrice: 8.88, HighPrice: 8.88},
				},
			},
		},
		{
			Instrument:   &roho.Instrument{Symbol: "wait-downward-option"},
			Quote:        &roho.Quote{AskPrice: 8.83},
			Fundamentals: &roho.Fundamental{Low52Weeks: 8.88},
			Historical: &roho.Historical{
				Records: []roho.HistoricalRecord{
					{OpenPrice: 8.88, HighPrice: 8.88, LowPrice: 8.87, ClosePrice: 8.87},
					{OpenPrice: 8.87, HighPrice: 8.87, LowPrice: 8.86, ClosePrice: 8.86},
					{OpenPrice: 8.86, HighPrice: 8.86, LowPrice: 8.85, ClosePrice: 8.85},
					{OpenPrice: 8.85, HighPrice: 8.85, LowPrice: 8.84, ClosePrice: 8.84},
				},
			},
		},
		{
			Instrument:   &roho.Instrument{Symbol: "buy-upward-option"},
			Quote:        &roho.Quote{AskPrice: 8.94},
			Fundamentals: &roho.Fundamental{Low52Weeks: 8.88},
			Historical: &roho.Historical{
				Records: []roho.HistoricalRecord{
					{OpenPrice: 8.84, LowPrice: 8.84, HighPrice: 8.85, ClosePrice: 8.85},
					{OpenPrice: 8.85, LowPrice: 8.85, HighPrice: 8.86, ClosePrice: 8.86},
					{OpenPrice: 8.86, LowPrice: 8.86, HighPrice: 8.87, ClosePrice: 8.87},
					{OpenPrice: 8.87, LowPrice: 8.87, HighPrice: 8.88, ClosePrice: 8.88},
					{OpenPrice: 8.88, LowPrice: 8.88, HighPrice: 8.89, ClosePrice: 8.89},
					{OpenPrice: 8.89, LowPrice: 8.89, HighPrice: 8.90, ClosePrice: 8.90},
				},
			},
		},
		{
			Instrument:   &roho.Instrument{Symbol: "ignore-far"},
			Quote:        &roho.Quote{BidPrice: 9.15, AskPrice: 9.15},
			Fundamentals: &roho.Fundamental{Low52Weeks: 8.88, High52Weeks: 9.99},
			Historical: &roho.Historical{
				Records: []roho.HistoricalRecord{
					{OpenPrice: 8.88},
				},
			},
		},
	}

	want := []Trade{
		{Instrument: &roho.Instrument{Symbol: "sell-downward-position"}, Order: roho.OrderOpts{Price: 8.84, Quantity: 2, Side: roho.Sell}, Reason: `0.5% away from 52-week high of 8.88, -0.11% bounce`},
		{Instrument: &roho.Instrument{Symbol: "buy-upward-option"}, Order: roho.OrderOpts{Price: 8.94, Quantity: 1, Side: roho.Buy}, Reason: "0.7% away from 52wk low of 8.88, 0.79% bounce"},
	}
	got, err := s.Trades(context.Background(), cs)
	if err != nil {
		t.Errorf("Trades() returned unexpected error: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unexpected trade diff: %s", diff)
	}
}

func TestUpward(t *testing.T) {
	tests := []struct {
		input    []float64
		wantOK   bool
		wantPerc float64
	}{
		{input: []float64{2, 2, 3, 4}, wantOK: true, wantPerc: 100.0},
		{input: []float64{2, 2, 1, 4}, wantOK: false, wantPerc: 100.0},
		{input: []float64{0.1, 0.2, 0.2, 0.3}, wantOK: true, wantPerc: 200},
		{input: []float64{100, 100, 101, 101}, wantOK: true, wantPerc: 1},
	}

	for _, tc := range tests {
		ok, perc := upward(tc.input)
		if ok != tc.wantOK {
			t.Errorf("upward(%v).ok = %v, want %v", tc.input, ok, tc.wantOK)
		}
		if math.Abs(percentDiff(perc, tc.wantPerc)) > 0.001 && fmt.Sprintf("%.3f", perc) != fmt.Sprintf("%.3f", tc.wantPerc) {
			t.Errorf("upward(%v).perc = %.3f, want %.3f", tc.input, perc, tc.wantPerc)
		}
	}
}
