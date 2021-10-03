package strategy

import (
	"context"
	"fmt"

	"github.com/tstromberg/roho/pkg/roho"
)

var (
	LuckySevens = "lucky-sevens"
	Random      = "random"
	HiLo        = "hilo"
	Bounce      = "bounce"
	strategies  = []string{LuckySevens, Random, HiLo}
)

type Trade struct {
	Instrument *roho.Instrument
	Order      roho.OrderOpts
	Reason     string
}

type CombinedStock struct {
	Quote        *roho.Quote
	Instrument   *roho.Instrument
	Fundamentals *roho.Fundamental
	Position     *roho.Position
	Historical   *roho.Historical
}

// Strategy is an interface for executing stock strategies.
type Strategy interface {
	Trades(ctx context.Context, cs []*CombinedStock) ([]Trade, error)
	String() string
}

type Config struct {
	Client *roho.Client
	Kind   string
	Values map[string]int64
}

// New returns a new strategy manager.
func New(c Config) (Strategy, error) {
	switch c.Kind {
	case LuckySevens:
		l := &LuckySevensStrategy{c: c}
		return l, nil
	case Random:
		l := &RandomStrategy{c: c}
		return l, nil
	case HiLo:
		l := &HiLoStrategy{c: c}
		return l, nil
	case Bounce:
		l := &BounceStrategy{c: c}
		return l, nil
	default:
		return nil, fmt.Errorf("no strategy named %q exists", c.Kind)
	}
}

// List returns a list of strategies.
func List() []string {
	return strategies
}
