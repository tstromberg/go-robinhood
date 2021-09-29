[![Go Reference](https://pkg.go.dev/badge/github.com/tstromberg/roho.svg)](https://pkg.go.dev/github.com/tstromberg/roho)
[![experimental](http://badges.github.io/stability-badges/dist/experimental.svg)](http://github.com/badges/stability-badges)
[![Go Report Card](https://goreportcard.com/badge/github.com/tstromberg/roho)](https://goreportcard.com/report/github.com/tstromberg/roho)

# RoHo

A new Go-language client library for accessing the Robinhood API. Based on https://github.com/andrewstuart/go-robinhood

![Roho logo](images/roho.png?raw=true "Roho")

## Features

* Idiomatic Go API
* Buying & selling shares and cryptocurrencies
* Access to historical quotes & crypto positions
* Common interface for trading strategies
* Designed for use by trading bots

## Library Usage

```go
// Login to Robinhood. Uses $RH_USER and $RH_PASS environment variables by default
r, err := roho.New(&roho.Config{})

// Lookup the SPDR S&P 500 ETF Trust
i, err := r.Instrument("SPY")

// Buy SPY at $100
o, err := r.Buy(i, roho.OrderOpts{Price: 100.0, Quantity: 1})

// Uh oh! Cancel this order ASAP!
o.Cancel()
```

For a runnable example, see [cmd/example](cmd/example).


## Trading Strategies

RoHo now ships with a `pkg/strategy` library to define and execute basic trading strategies.

*NOTE: You will lose money if you use RoHo's trading strategies feature*


Here is a simplified example from the [hilo example strategy](pkg/strategy/hilo.go), to buy low, sell high:

```go
func (cr *HiLoStrategy) Trades(_ context.Context, cs map[string]*CombinedStock) ([]Trade, error) {
	ts := []Trade{}

	for url, s := range cs {
		p := s.Position

		if p == nil {
			ratio := s.Quote.BidPrice / s.Fundamentals.Low52Weeks
			if ratio < 1.01 {
				ts = append(ts, Trade{
					InstrumentURL: url,
					Symbol:        s.Quote.Symbol,
					Order:         roho.OrderOpts{...},
				})
			}
			continue
		}

		ratio := s.Quote.BidPrice / s.Fundamentals.High52Weeks
		if ratio > 0.99 {
			if p.AverageBuyPrice < s.Quote.BidPrice {
				continue
			}

			ts = append(ts, Trade{
				InstrumentURL: url,
				Symbol:        s.Quote.Symbol,
				Order:         roho.OrderOpts{...},
			})
		}
	}

	return ts, nil
}
```


You can run these strategies by using [cmd/matador](cmd/matador)

## Approximate 2021 Roadmap

* Add backtesting support
* Add response caching for offline testing
