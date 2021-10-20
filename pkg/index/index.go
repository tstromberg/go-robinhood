// index package supplies a list of sources for symbols
package index

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// TODO: NDX, Russel 2000, DJI, Nasdaq Composite
// TODO: Cache results for 24h

// SP500 returns the symbols on the S&P 500.
func SP500(ctx context.Context) ([]string, error) {
	symbols := []string{}

	req, err := http.NewRequestWithContext(ctx, "GET", "https://en.wikipedia.org/wiki/List_of_S%26P_500_companies", nil)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	doc.Find("#constituents tr").Each(func(i int, s *goquery.Selection) {
		sym := s.Find("td a").First().Text()
		if len(sym) > 0 && len(sym) < 5 {
			symbols = append(symbols, sym)
		}
	})

	return symbols, nil
}

// SP50 returns the top-50 from the S&P 500.
func SP50(ctx context.Context) ([]string, error) {
	symbols, err := SP500(ctx)
	if err != nil {
		return nil, err
	}
	return symbols[0:50], nil
}

func Resolve(ctx context.Context, syms []string) ([]string, error) {
	rs := []string{}
	for _, s := range syms {
		if !strings.HasPrefix(s, "^") {
			rs = append(rs, s)
			continue
		}

		switch s {
		case "^SP500":
			sp, err := SP500(ctx)
			if err != nil {
				return rs, err
			}
			rs = append(rs, sp...)
		case "^SP50":
			sp, err := SP50(ctx)
			if err != nil {
				return rs, err
			}
			rs = append(rs, sp...)
		default:
			return rs, fmt.Errorf("unknown index: %q", s)
		}
	}

	return rs, nil
}
