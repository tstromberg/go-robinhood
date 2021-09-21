// Matador runs trading strategies
package main

// usage:
//
// RH_USER=email@example.org RH_PASS=password go run .

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/tstromberg/roho/pkg/roho"
	"github.com/tstromberg/roho/pkg/strategy"
)

func main() {
	ctx := context.Background()
	r, err := roho.New(ctx, &roho.Config{})
	if err != nil {
		log.Fatalf("new failed: %v", err)
	}

	if len(os.Args) < 3 {
		log.Fatalf("syntax: matador [strategy] [symbols]")
	}

	symbols := os.Args[2:]

	cr, err := strategy.New(ctx, strategy.Config{Client: r, Kind: os.Args[1]})
	if err != nil {
		log.Fatalf("strategy failed: %v", err)
	}

	for {
		ts, err := cr.Trades(ctx, symbols)
		if err != nil {
			log.Fatalf("trades failed: %v", err)
		}

		for _, t := range ts {
			log.Printf("TRADE: %+v", t)
		}

		log.Printf("Sleeping for 1 minute")
		time.Sleep(1 * time.Minute)
	}
}
