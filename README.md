[![Go Reference](https://pkg.go.dev/badge/github.com/tstromberg/roho.svg)](https://pkg.go.dev/github.com/tstromberg/roho)
[![experimental](http://badges.github.io/stability-badges/dist/experimental.svg)](http://github.com/badges/stability-badges)
[![Go Report Card](https://goreportcard.com/badge/github.com/tstromberg/roho)](https://goreportcard.com/report/github.com/tstromberg/roho)

# RoHo

A new Go-language client library for accessing the Robinhood API. Based on https://github.com/andrewstuart/go-robinhood

![Roho logo](images/roho.png?raw=true "Roho")

## Features

* Idiomatic Go API
* Buying & selling shares and cryptocurrencies
* Access to historical quotes
* Designed for use by trading bots

## Usage

```go
// Login to Robinhood. Uses $RH_USER and $RH_PASS environment variables by default
r, err := roho.New(&roho.Config{})

// Lookup the SPDR S&P 500 ETF Trust
i, err := r.Lookup("SPY")

// Buy SPY at $100
o, err := r.Buy(i, roho.OrderOpts{Price: 100.0, Quantity: 1})

// Uh oh! Cancel this order ASAP!
o.Cancel()
```

For a runnable example, see [cmd/example](cmd/example).

## Approximate 2021 Roadmap

* Clean up codebase
* Add CI
* Add 1st-class support for executing pluggable trading strategies
* Add 1st-class support for backtesting trading strategies
* Add response caching for offline testing
* Profit!
