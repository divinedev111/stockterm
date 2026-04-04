package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/divinedev111/stockterm/internal/api"
	"github.com/divinedev111/stockterm/internal/chart"
	"github.com/divinedev111/stockterm/internal/quote"
	"github.com/divinedev111/stockterm/internal/ui"
	"github.com/divinedev111/stockterm/internal/watchlist"
)

func main() {
	symbols := flag.String("symbols", "", "Comma-separated symbols to watch (e.g., AAPL,GOOGL,BTC-USD)")
	interval := flag.Duration("interval", 30*time.Second, "Refresh interval")
	oneShot := flag.Bool("once", false, "Fetch once and exit (no live updates)")
	chartSym := flag.String("chart", "", "Show chart for a specific symbol")
	chartRange := flag.String("range", "1mo", "Chart time range (1d, 5d, 1mo, 3mo, 6mo, 1y)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: stockterm [options]\n\n")
		fmt.Fprintf(os.Stderr, "Bloomberg terminal in your terminal. Real-time stocks & crypto.\n\n")
		fmt.Fprintf(os.Stderr, "  stockterm                          # use default watchlist\n")
		fmt.Fprintf(os.Stderr, "  stockterm -symbols AAPL,GOOGL,TSLA # watch specific symbols\n")
		fmt.Fprintf(os.Stderr, "  stockterm -chart AAPL -range 3mo   # show chart\n")
		fmt.Fprintf(os.Stderr, "  stockterm -once                    # fetch and exit\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	client := api.NewClient()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Chart mode
	if *chartSym != "" {
		showChart(ctx, client, *chartSym, *chartRange)
		return
	}

	// Determine symbols
	var syms []string
	if *symbols != "" {
		for _, s := range strings.Split(*symbols, ",") {
			s = strings.TrimSpace(strings.ToUpper(s))
			if s != "" {
				syms = append(syms, s)
			}
		}
	} else {
		wl, err := watchlist.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading watchlist: %v\n", err)
			os.Exit(1)
		}
		syms = wl.Symbols()
	}

	if len(syms) == 0 {
		fmt.Fprintln(os.Stderr, "No symbols to watch. Use -symbols or add to your watchlist.")
		os.Exit(1)
	}

	// Fetch and display
	quotes := fetchQuotes(ctx, client, syms)
	fmt.Print(ui.RenderQuotes(quotes))

	if *oneShot {
		return
	}

	// Live refresh
	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Print("\033[?25h") // show cursor
			return
		case <-ticker.C:
			quotes = fetchQuotes(ctx, client, syms)
			fmt.Print(ui.RenderQuotes(quotes))
		}
	}
}

func fetchQuotes(ctx context.Context, client *api.Client, symbols []string) []*quote.Quote {
	quotes := make([]*quote.Quote, 0, len(symbols))
	for _, sym := range symbols {
		q, err := client.GetQuote(ctx, sym)
		if err != nil {
			quotes = append(quotes, &quote.Quote{
				Symbol: sym,
				Name:   fmt.Sprintf("(error: %s)", truncate(err.Error(), 30)),
			})
			continue
		}
		quotes = append(quotes, q)
	}
	return quotes
}

func showChart(ctx context.Context, client *api.Client, symbol, timeRange string) {
	q, err := client.GetQuote(ctx, symbol)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching quote: %v\n", err)
		os.Exit(1)
	}

	chartInterval := "1d"
	switch timeRange {
	case "1d":
		chartInterval = "5m"
	case "5d":
		chartInterval = "15m"
	case "1mo":
		chartInterval = "1d"
	case "3mo", "6mo":
		chartInterval = "1d"
	case "1y", "2y":
		chartInterval = "1wk"
	}

	candles, err := client.GetChart(ctx, symbol, chartInterval, timeRange)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching chart: %v\n", err)
		os.Exit(1)
	}

	chartStr := chart.Render(candles, 60, 15)
	fmt.Print(ui.RenderDetail(q, chartStr))
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
