[![CI](https://github.com/divinedev111/stockterm/actions/workflows/ci.yml/badge.svg)](https://github.com/divinedev111/stockterm/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/divinedev111/stockterm)](https://goreportcard.com/report/github.com/divinedev111/stockterm)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

# stockterm

Bloomberg terminal in your terminal. Real-time stock and crypto quotes, ASCII charts, and watchlists — powered by Yahoo Finance.

```
  stockterm  —  press q to quit, a to add, d to delete, c to chart
  ──────────────────────────────────────────────────────────────────────────────
  SYMBOL   NAME                        PRICE     CHANGE    CHG%       VOLUME
  ──────────────────────────────────────────────────────────────────────────────
  AAPL     Apple Inc.                $195.50   ▲ 2.30     1.19%       55.0M
  GOOGL    Alphabet Inc.             $178.25   ▲ 1.15     0.65%       28.3M
  MSFT     Microsoft Corporation     $425.80   ▼ -3.20    -0.75%      32.1M
  AMZN     Amazon.com, Inc.          $189.30   ▲ 4.50     2.43%       41.7M
  TSLA     Tesla, Inc.               $245.00   ▼ -5.00    -2.00%     120.0M
  ──────────────────────────────────────────────────────────────────────────────
```

## Install

```bash
go install github.com/divinedev111/stockterm/cmd/stockterm@latest
```

## Usage

```bash
# Default watchlist (AAPL, GOOGL, MSFT, AMZN, TSLA)
stockterm

# Watch specific symbols
stockterm -symbols AAPL,GOOGL,BTC-USD,ETH-USD

# Show chart for a symbol
stockterm -chart AAPL -range 3mo

# Fetch once and exit (no live refresh)
stockterm -once

# Custom refresh interval
stockterm -interval 10s
```

### Chart View

```bash
stockterm -chart TSLA -range 1mo
```

```
  TSLA  Tesla, Inc.
  ────────────────────────────────────────────────────────

  Price:   $245.00
  Change:  ▼ -5.00 (-2.00%)
  Open:    $250.00
  High:    $252.30
  Low:     $244.10
  Volume:  120.0M
  MktCap:  780.5B

   252 │██
       │  ██                          ██
   248 │    ████                    ██  ██
       │        ██░░░░░░░░░░░░░░██░░░░░░██
   244 │░░░░░░░░░░██░░░░░░░░░░██░░░░░░░░░░██
       │░░░░░░░░░░░░██░░░░░░██░░░░░░░░░░░░░░
   240 │░░░░░░░░░░░░░░████░░░░░░░░░░░░░░░░░░
       └──────────────────────────────────────────
       Mar 1                              Mar 31
```

### Crypto

Works with crypto symbols using Yahoo Finance tickers:

```bash
stockterm -symbols BTC-USD,ETH-USD,SOL-USD,DOGE-USD
```

## Options

| Flag | Description | Default |
|------|-------------|---------|
| `-symbols` | Comma-separated symbols | Default watchlist |
| `-interval` | Refresh interval | `30s` |
| `-once` | Fetch once and exit | `false` |
| `-chart` | Show chart for symbol | — |
| `-range` | Chart range: `1d`, `5d`, `1mo`, `3mo`, `6mo`, `1y` | `1mo` |

## Watchlist

Saved to `~/.config/stockterm/watchlist.json`. Default: `AAPL, GOOGL, MSFT, AMZN, TSLA`.

## Architecture

```
cmd/stockterm/
  main.go              CLI and refresh loop
internal/
  api/
    yahoo.go           Yahoo Finance API client (quotes + chart data)
  chart/
    ascii.go           ASCII chart renderer with area fill
  quote/
    quote.go           Quote and OHLC data types
  ui/
    render.go          Terminal dashboard renderer (ANSI colors)
  watchlist/
    watchlist.go       Persistent watchlist with JSON storage
```

## License

MIT
