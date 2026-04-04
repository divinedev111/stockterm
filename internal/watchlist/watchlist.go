// Package watchlist manages a user's saved stock symbols.
package watchlist

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Watchlist holds a set of symbols the user is tracking.
type Watchlist struct {
	mu      sync.RWMutex
	symbols []string
	path    string
}

// Load reads or creates a watchlist from the default config path.
func Load() (*Watchlist, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = os.TempDir()
	}

	path := filepath.Join(dir, "stockterm", "watchlist.json")
	w := &Watchlist{path: path}

	data, err := os.ReadFile(path)
	if err != nil {
		// Default watchlist
		w.symbols = []string{"AAPL", "GOOGL", "MSFT", "AMZN", "TSLA"}
		return w, nil
	}

	if err := json.Unmarshal(data, &w.symbols); err != nil {
		w.symbols = []string{"AAPL", "GOOGL", "MSFT", "AMZN", "TSLA"}
	}

	return w, nil
}

// Symbols returns the current watchlist.
func (w *Watchlist) Symbols() []string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	out := make([]string, len(w.symbols))
	copy(out, w.symbols)
	return out
}

// Add appends a symbol if not already present.
func (w *Watchlist) Add(symbol string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, s := range w.symbols {
		if s == symbol {
			return
		}
	}
	w.symbols = append(w.symbols, symbol)
	w.save()
}

// Remove deletes a symbol from the watchlist.
func (w *Watchlist) Remove(symbol string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	for i, s := range w.symbols {
		if s == symbol {
			w.symbols = append(w.symbols[:i], w.symbols[i+1:]...)
			w.save()
			return
		}
	}
}

// Has returns true if the symbol is in the watchlist.
func (w *Watchlist) Has(symbol string) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	for _, s := range w.symbols {
		if s == symbol {
			return true
		}
	}
	return false
}

func (w *Watchlist) save() {
	dir := filepath.Dir(w.path)
	os.MkdirAll(dir, 0o755)

	data, _ := json.MarshalIndent(w.symbols, "", "  ")
	os.WriteFile(w.path, data, 0o644)
}
