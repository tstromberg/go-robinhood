package strategy

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/tstromberg/roho/pkg/roho"
)

// LiveData gathers combined live stock information.
func LiveData(ctx context.Context, r *roho.Client, syms []string) ([]*CombinedStock, error) {
	cu := map[string]*CombinedStock{}

	// Using the instruments API is exceptionally slow
	is, err := r.Instruments(ctx, syms)
	if err != nil {
		return nil, fmt.Errorf("instruments: %w", err)
	}

	for _, i := range is {
		i := i
		cu[i.URL] = &CombinedStock{Instrument: &i}
	}

	return updateMapData(ctx, r, cu)
}

func updateMapData(ctx context.Context, r *roho.Client, cu map[string]*CombinedStock) ([]*CombinedStock, error) {
	ps, err := r.Positions(ctx)
	if err != nil {
		return nil, fmt.Errorf("positions: %w", err)
	}

	for _, p := range ps {
		p := p
		_, ok := cu[p.InstrumentURL]
		if !ok {
			i, err := r.InstrumentFromURL(ctx, p.InstrumentURL)
			if err != nil {
				return nil, fmt.Errorf("instruments: %w", err)
			}
			cu[p.InstrumentURL] = &CombinedStock{Instrument: &i}
		}
		cu[p.InstrumentURL].Position = &p
	}

	syms := []string{}
	for _, s := range cu {
		syms = append(syms, s.Instrument.Symbol)
	}

	qs, err := r.Quotes(ctx, syms)
	if err != nil {
		return nil, fmt.Errorf("quote: %w", err)
	}

	for _, q := range qs {
		// avoid implicit memory aliasing within a for loop
		q := q
		_, ok := cu[q.InstrumentURL]
		if !ok {
			cu[q.InstrumentURL] = &CombinedStock{}
		}
		cu[q.InstrumentURL].Quote = &q
	}

	fs, err := r.Fundamentals(ctx, syms...)
	if err != nil {
		return nil, fmt.Errorf("fundamentals: %w", err)
	}

	for _, f := range fs {
		// avoid implicit memory aliasing within a for loop
		f := f
		cu[f.InstrumentURL].Fundamentals = &f
	}

	cs := []*CombinedStock{}
	for _, c := range cu {
		cs = append(cs, c)
	}

	sort.Slice(cs, func(i, j int) bool {
		if cs[i].Position != nil && cs[j].Position == nil {
			return true
		}
		return cs[i].Instrument.URL < cs[j].Instrument.URL
	})
	return cs, nil
}

// UpdateData updates stock information.
func UpdateData(ctx context.Context, r *roho.Client, cs []*CombinedStock) ([]*CombinedStock, error) {
	cu := map[string]*CombinedStock{}
	for _, s := range cs {
		s := s
		cu[s.Instrument.URL] = s
	}
	return updateMapData(ctx, r, cu)
}

// HistoricalData simulates data at a particular point in the past - NOT YET IMPLEMENTED.
func HistoricalData(_ context.Context, _ *roho.Client, _ []string, _ time.Time) ([]*CombinedStock, error) {
	cs := []*CombinedStock{}
	return cs, fmt.Errorf("HistoricalData is NOT YET IMPLEMENTED")
}
