package watchlist

import "testing"

func TestWatchlistAddRemove(t *testing.T) {
	w := &Watchlist{symbols: []string{"AAPL"}, path: "/dev/null"}

	w.Add("GOOGL")
	if len(w.Symbols()) != 2 {
		t.Errorf("after add, len = %d", len(w.Symbols()))
	}

	// No duplicates
	w.Add("GOOGL")
	if len(w.Symbols()) != 2 {
		t.Error("should not add duplicate")
	}

	w.Remove("AAPL")
	if len(w.Symbols()) != 1 {
		t.Errorf("after remove, len = %d", len(w.Symbols()))
	}
	if w.Symbols()[0] != "GOOGL" {
		t.Error("remaining symbol should be GOOGL")
	}
}

func TestWatchlistHas(t *testing.T) {
	w := &Watchlist{symbols: []string{"AAPL", "TSLA"}, path: "/dev/null"}

	if !w.Has("AAPL") {
		t.Error("should have AAPL")
	}
	if w.Has("MSFT") {
		t.Error("should not have MSFT")
	}
}

func TestWatchlistSymbolsCopy(t *testing.T) {
	w := &Watchlist{symbols: []string{"AAPL"}, path: "/dev/null"}
	syms := w.Symbols()
	syms[0] = "MODIFIED"

	if w.Symbols()[0] != "AAPL" {
		t.Error("Symbols() should return a copy")
	}
}

func TestWatchlistRemoveNonExistent(t *testing.T) {
	w := &Watchlist{symbols: []string{"AAPL"}, path: "/dev/null"}
	w.Remove("MISSING")
	if len(w.Symbols()) != 1 {
		t.Error("removing non-existent should be no-op")
	}
}
