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

	klog.Infof("%s %d shares of %q at %.2f: %q ...", act, t.Order.Quantity, t.Instrument.Symbol, t.Order.Price, t.Reason)
	if dryRun {
		return nil
	}
	out, err := r.Order(ctx, t.Instrument.URL, t.Instrument.Symbol, t.Order)
	klog.Infof("order result: %+v", out)
	return err
}

type Counter struct {
	TotalBuys  int
	TotalSales int
	PollBuys   int
	PollSales  int
	Polls      int
}

func loop(ctx context.Context, r *roho.Client, st strategy.Strategy, syms []string) {
	klog.Infof("%q loop has begun with %d symbols!", st, len(syms))

	maxSleep := *maxPollFlag - *minPollFlag
	counter := &Counter{}

	klog.Infof("Gathering live data for %d symbols ...", len(syms))
	tctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	combined, err := strategy.LiveData(tctx, r, syms)
	if err != nil {
		klog.Errorf("live data: %v", err)
		return
	}

	for {
		counter.Polls++

		if counter.Polls > 1 {
			klog.Infof("%d buys, %d sells. Sleeping for %s...", counter.TotalBuys, counter.TotalSales, maxSleep)
			time.Sleep(maxSleep)
			klog.Infof("Updating data for %s symbols ...", len(combined))
			combined, err = strategy.UpdateData(ctx, r, combined)
			if err != nil {
				klog.Errorf("failed to update data: %v", err)
				continue
			}
		}

		cont, err := check(ctx, r, st, combined, *dryRunFlag, counter)
		if err != nil {
			klog.Errorf("check failed: %v", err)
			continue
		}

		if !cont {
			klog.Infof("loop has completed with %d buys and %d sells", counter.TotalBuys, counter.TotalSales)
			return
		}
	}
}

func check(ctx context.Context, r *roho.Client, st strategy.Strategy, combined []*strategy.CombinedStock, dryRun bool, count *Counter) (bool, error) {
	klog.Infof("Calculating trades for %d stocks ...", len(combined))
	ts, err := st.Trades(ctx, combined)
	if err != nil {
		return true, fmt.Errorf("trades: %w", err)
	}

	if len(ts) == 0 {
		return true, nil
	}

	if len(ts) > 0 {
		klog.Infof("%d possible trades found ...", len(ts))
	}

	count.PollSales = 0
	count.PollBuys = 0

	for _, t := range ts {
		if t.Order.Side == roho.Buy {
			if count.PollBuys+1 > *maxBuysPerPollFlag {
				klog.Warningf(" -> BUY %s (ignoring, over max-buys-per-poll=%d): %+v", t.Instrument.Symbol, *maxBuysPerPollFlag, t)
				continue
			}

			if count.TotalBuys+1 > *maxBuysFlag {
				return false, fmt.Errorf("hit maximum buys (%d)", *maxBuysFlag)
			}

			count.PollBuys++
			count.TotalBuys++
		}

		if t.Order.Side == roho.Sell {
			if count.PollSales+1 > *maxSalesPerPollFlag {
				klog.Warningf(" -> SELL %s (ignoring, over max-sales-per-poll=%d): %+v", t.Instrument.Symbol, *maxSalesPerPollFlag, t)
				continue
			}

			if count.TotalSales+1 > *maxSalesFlag {
				return false, fmt.Errorf("hit maximum sales (%d)", *maxSalesFlag)
			}
		}

		if err := trade(ctx, r, t, dryRun); err != nil {
			return true, fmt.Errorf("trade failed: %w", err)
		}
	}

	return true, nil
}
