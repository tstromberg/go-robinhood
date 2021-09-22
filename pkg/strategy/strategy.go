package strategy

import (
	"context"
	"fmt"
	"time"

	"github.com/tstromberg/roho/pkg/roho"
)

var (
	LuckySevens = "lucky-sevens"
	Random      = "random"
	strategies  = []string{LuckySevens, Random}
)

type Trade struct {
	Symbol        string
	InstrumentURL string
	Order         roho.OrderOpts
	Reason        string
}

// Strategy is an interface for executing stock strategies
type Strategy interface {
	Trades(ctx context.Context, ps []roho.Position, qs []roho.Quote) ([]Trade, error)
	SetTime(ctx context.Context, t time.Time) error
	String() string
}

type Config struct {
	Client   *roho.Client
	Kind     string
	Holdings []string
}

// New returns a new strategy manager
func New(ctx context.Context, c Config) (Strategy, error) {
	switch c.Kind {
	case LuckySevens:
		l := &LuckySevensStrategy{c: c}
		return l, nil
	case Random:
		l := &RandomStrategy{c: c}
		return l, nil
	default:
		return nil, fmt.Errorf("no strategy named %q exists", c.Kind)
	}
}

// List returns a list of strategies
func List() []string {
	return strategies
}
