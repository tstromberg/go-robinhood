package strategy

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tstromberg/roho/pkg/roho"
)

// LiveData gathers combined live stock information.
func LiveData(ctx context.Context, r *roho.Client, syms []string) (map[string]*CombinedStock, error) {
	cs := map[string]*CombinedStock{}

	ps, err := r.Positions(ctx)
	if err != nil {
		return nil, fmt.Errorf("positions: %w", err)
	}

	for _, p := range ps {
		// avoid implicit memory aliasing within a for loop
		p := p
		log.Printf("position: %s", p.InstrumentURL)
		cs[p.InstrumentURL] = &CombinedStock{Position: &p}

		f, err := r.InstrumentFromURL(ctx, p.InstrumentURL)
		if err != nil {
			return cs, fmt.Errorf("instrument from %q: %w", p.InstrumentURL, err)
		}

		syms = append(syms, f.Symbol)
	}

	qs, err := r.Quotes(ctx, syms)
	if err != nil {
		return nil, fmt.Errorf("quote: %w", err)
	}

	for _, q := range qs {
		// avoid implicit memory aliasing within a for loop
		q := q
		_, ok := cs[q.InstrumentURL]
		if !ok {
			cs[q.InstrumentURL] = &CombinedStock{}
		}
		cs[q.InstrumentURL].Quote = &q
	}

	fs, err := r.Fundamentals(ctx, syms...)
	if err != nil {
		return nil, fmt.Errorf("fundamentals: %w", err)
	}

	for _, f := range fs {
		// avoid implicit memory aliasing within a for loop
		f := f
		cs[f.InstrumentURL].Fundamentals = &f
	}

	return cs, nil
}

// HistoricalData simulates data at a particular point in the past - NOT YET IMPLEMENTED.
func HistoricalData(_ context.Context, _ *roho.Client, _ []string, _ time.Time) (map[string]*CombinedStock, error) {
	cs := map[string]*CombinedStock{}
	return cs, fmt.Errorf("HistoricalData is NOT YET IMPLEMENTED")
}
