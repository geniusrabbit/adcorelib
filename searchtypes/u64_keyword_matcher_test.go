package searchtypes

import (
	"sync"
	"testing"
)

func TestU64KeywordMatcher_UserExample(t *testing.T) {
	m := NewU64KeywordMatcher().
		Set("automobile", 1000).
		Set("motocicles", 1000).
		Set("bike racing", 1000).
		Freeze()

	cases := []struct {
		in   string
		want uint64
	}{
		{"cars", 0},
		{"motosport", 0},
		{"world bike racings", 1000},
		{"bike racing", 1000},
		{"BIKE RACING", 1000},
		{"automobile", 1000},
		{"motocicles", 1000},
	}
	for _, tc := range cases {
		if got := m.Match(tc.in); got != tc.want {
			t.Fatalf("Match(%q)=%d want %d", tc.in, got, tc.want)
		}
	}
}

func TestU64KeywordMatcher_NoMidWord(t *testing.T) {
	m := NewU64KeywordMatcher().
		Set("motocicles", 7).
		Set("icycle", 8).
		Freeze()

	if got := m.Match("xmotocicles"); got != 0 {
		t.Fatalf("mid-word start Match=%d want 0", got)
	}
	if got := m.Match("motocicles"); got != 7 {
		t.Fatalf("exact Match=%d want 7", got)
	}
	// "icycle" as its own word matches phrase "icycle", not as mid of motocicles.
	if got := m.Match("big icycle"); got != 8 {
		t.Fatalf("word-start Match=%d want 8", got)
	}
}

func TestU64KeywordMatcher_LongestWins(t *testing.T) {
	m := NewU64KeywordMatcher().
		Set("bike", 1).
		Set("bike racing", 2).
		Freeze()

	if got := m.Match("world bike racings"); got != 2 {
		t.Fatalf("longest Match=%d want 2", got)
	}
	if got := m.Match("bike"); got != 1 {
		t.Fatalf("short Match=%d want 1", got)
	}
}

func TestU64KeywordMatcher_CaseFoldAndWhitespace(t *testing.T) {
	m := NewU64KeywordMatcher().
		Set("  Bike   Racing  ", 9).
		Freeze()

	if got := m.Match("\tWORLD  BIKE   RACINGS\n"); got != 9 {
		t.Fatalf("Match=%d want 9", got)
	}
	if got := m.Match("   "); got != 0 {
		t.Fatalf("empty Match=%d want 0", got)
	}
	if got := m.Match(""); got != 0 {
		t.Fatalf("empty string Match=%d want 0", got)
	}
}

func TestU64KeywordMatcher_SkipEmptySet(t *testing.T) {
	m := NewU64KeywordMatcher().
		Set("", 1).
		Set("   ", 2).
		Set("ok", 3).
		Freeze()
	if got := m.Match("ok"); got != 3 {
		t.Fatalf("Match=%d want 3", got)
	}
}

func TestU64KeywordMatcher_OverwritePhrase(t *testing.T) {
	m := NewU64KeywordMatcher().
		Set("cars", 1).
		Set("cars", 2).
		Freeze()
	if got := m.Match("cars"); got != 2 {
		t.Fatalf("Match=%d want 2", got)
	}
}

func TestU64KeywordMatcher_NilAndUnfrozen(t *testing.T) {
	var m *U64KeywordMatcher
	if got := m.Match("x"); got != 0 {
		t.Fatalf("nil Match=%d want 0", got)
	}
	u := NewU64KeywordMatcher().Set("x", 1)
	if got := u.Match("x"); got != 0 {
		t.Fatalf("unfrozen Match=%d want 0", got)
	}
}

func TestU64KeywordMatcher_FreezeIdempotent(t *testing.T) {
	m := NewU64KeywordMatcher().Set("a", 1).Freeze().Freeze()
	if got := m.Match("a"); got != 1 {
		t.Fatalf("Match=%d want 1", got)
	}
}

func TestU64KeywordMatcher_FreezePanicOnSet(t *testing.T) {
	m := NewU64KeywordMatcher().Set("a", 1).Freeze()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on Set after Freeze")
		}
	}()
	m.Set("b", 2)
}

func TestU64KeywordMatcher_ConcurrentMatch(t *testing.T) {
	m := NewU64KeywordMatcher().
		Set("automobile", 1000).
		Set("bike racing", 1000).
		Freeze()

	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				if got := m.Match("world bike racings"); got != 1000 {
					t.Errorf("Match=%d want 1000", got)
					return
				}
			}
		}()
	}
	wg.Wait()
}

func TestNormalizeKeywordPhrase(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"  Bike   Racing  ", "bike racing"},
		{"A", "a"},
		{"", ""},
		{" \t ", ""},
	}
	for _, tc := range cases {
		if got := normalizeKeywordPhrase(tc.in); got != tc.want {
			t.Fatalf("normalize(%q)=%q want %q", tc.in, got, tc.want)
		}
	}
}
