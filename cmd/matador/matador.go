// Matador is a demonstration program to run trading strategies
package main

// usage:
//
// RH_USER=email@example.org RH_PASS=password go run .

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/tstromberg/roho/pkg/index"
	"github.com/tstromberg/roho/pkg/roho"
	"github.com/tstromberg/roho/pkg/strategy"
	"k8s.io/klog/v2"
)

var (
	dryRunFlag          = flag.Bool("dry-run", false, "dry-run mode (don't buy/sell anything)")
	strategyFlag        = flag.String("strategy", "", fmt.Sprintf("strategy to use. Choices: %v", strategy.List()))
	minPollFlag         = flag.Duration("min-poll", 5*time.Second, "minimum time to poll (even if errors happen)")
	maxPollFlag         = flag.Duration("max-poll", 60*time.Second, "maximum time to poll")
	maxBuysFlag         = flag.Int("max-buys", 5, "maximum buys before exiting")
	maxBuysPerPollFlag  = flag.Int("max-buys-per-poll", 1, "maximum buys per polling period")
	maxSalesFlag        = flag.Int("max-sales", 5, "maximum sales before exiting")
	maxSalesPerPollFlag = flag.Int("max-sales-per-poll", 1, "maximum sales per polling period")
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	if !*dryRunFlag {
		klog.Warningf("matador is not in dry-run mode. You will lose money (sleeping for 10s)")
		time.Sleep(10 * time.Second)
	}

	ctx := context.Background()
	r, err := roho.New(ctx, &roho.Config{})
	if err != nil {
		klog.Fatalf("new failed: %v", err)
	}

	klog.Infof("args=%v (dry-run=%v, strategy=%v)", os.Args, *dryRunFlag, *strategyFlag)

	if len(flag.Args()) < 1 {
		klog.Fatalf("usage: matador --strategy=X [symbols]")
	}

	syms, err := index.Resolve(ctx, flag.Args())
	if err != nil {
		klog.Fatalf("failed to resolve symbols: %v", err)
	}

	if len(syms) == 0 {
		klog.Errorf("no symbols were resolved. usage: matador --strategy=X [symbols]")
		os.Exit(1)
	}

	st, err := strategy.New(strategy.Config{Client: r, Kind: *strategyFlag})
	if err != nil {
		klog.Errorf("strategy failed: %v", err)
		os.Exit(1)
	}

	loop(ctx, r, st, syms)
}

func trade(ctx context.Context, r *roho.Client, t strategy.Trade, dryRun bool) error {
	act := "Selling"
	if t.Order.Side == roho.Buy {
		act = "Buying"
	}

	if dryRun {
		act = "[DRY RUN] " + act
	}

	sym := t.Symbol
	if sym == "" && t.InstrumentURL != "" {
		i, err := r.InstrumentFromURL(ctx, t.InstrumentURL)
		if err != nil {
			return fmt.Errorf("instrument from URL: %w", err)
		}
		sym = i.Symbol
	}

	klog.Infof("%s %d shares of %q at %.2f ...", act, t.Order.Quantity, sym, t.Order.Price)
	if dryRun {
		return nil
	}
	out, err := r.Order(ctx, t.InstrumentURL, t.Symbol, t.Order)
	klog.Infof("order result: %+v", out)
	return err
}

func loop(ctx context.Context, r *roho.Client, st strategy.Strategy, syms []string) {
	totalBuys := 0
	totalSales := 0
	klog.Infof("%q loop has begun with %d symbols!", st, len(syms))

	maxSleep := *maxPollFlag - *minPollFlag

	for {
		klog.Infof("sleeping for %s ...", *minPollFlag)
		time.Sleep(*minPollFlag)

		// TODO: update instead of rebuild
		combined, err := strategy.LiveData(ctx, r, syms)
		if err != nil {
			klog.Errorf("live data failed: %v", err)
			continue
		}

		ts, err := st.Trades(ctx, combined)
		if err != nil {
			klog.Errorf("trades failed: %v", err)
			continue
		}

		sales := 0
		buys := 0
		for _, t := range ts {
			if t.Order.Side == roho.Buy {
				buys++
				if buys > *maxBuysPerPollFlag {
					klog.Warningf(" -> BUY %s (ignoring, over max-buys-per-poll=%d): %+v", t.Symbol, *maxBuysPerPollFlag, t)
					continue
				}

				totalBuys++
				if totalBuys > *maxBuysFlag {
					klog.Errorf("hit maximum buys (%d) - exiting", totalBuys)
					return
				}
			}

			if t.Order.Side == roho.Sell {
				sales++

				if sales > *maxSalesPerPollFlag {
					klog.Warningf(" -> SELL %s (ignoring, over max-sales-per-poll=%d): %+v", t.Symbol, *maxSalesPerPollFlag, t)
					continue
				}

				totalSales++
				if totalSales > *maxSalesFlag {
					klog.Errorf("hit maximum sales (%d) - exiting", totalSales)
					return
				}
			}

			if err := trade(ctx, r, t, *dryRunFlag); err != nil {
				klog.Errorf("trade failed: %v", err)
			}
		}

		klog.Infof("Sleeping for %s...", maxSleep)
		time.Sleep(maxSleep)
	}
}
