package strategy

import (
	"context"
	"fmt"
	"time"

	"github.com/tstromberg/roho/pkg/roho"
)

type Trade struct {
	Instrument *roho.Instrument
	Order      roho.OrderOpts
	Reason     string
}

// Commander is a common interface for executing strategies
type Commander interface {
	Trades(context.Context, []string) ([]Trade, error)
	Symbols(context.Context) ([]string, error)
	Simulate(context.Context, time.Time) error
	String() string
}

type Config struct {
	Client   *roho.Client
	Kind     string
	Holdings []string
}

// New returns a new strategy manager
func New(ctx context.Context, c Config) (Commander, error) {
	switch c.Kind {
	case "lucky-sevens":
		l := &LuckySevens{c: c}
		return l, nil
	default:
		return nil, fmt.Errorf("no strategy named %q exists", c.Kind)
	}
}
