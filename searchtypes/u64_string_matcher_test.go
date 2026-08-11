package searchtypes

import (
	"reflect"
	"strings"
	"testing"
)

func sampleCategories() *U64StringMatcher {
	return NewU64StringMatcher().
		SetNodes(
			Node{ID: 1, Value: "IAB1"},
			Node{ParentID: 1, ID: 2, Value: "IAB1-1"},
			Node{ParentID: 2, ID: 4, Value: "IAB1-1-1"},
			Node{ParentID: 1, ID: 3, Value: "IAB1-2"},
		).Freeze()
}

func TestU64StringMatcher_MatchID(t *testing.T) {
	m := sampleCategories()

	cases := []struct {
		name string
		id   uint64
		key  string
		want uint64
	}{
		{"forest nested", 0, "iab1-1-1", 4},
		{"self", 1, "iab1", 1},
		{"nested under root", 1, "IAB1-1-1", 4},
		{"direct child", 1, "iab1-2", 3},
		{"outside subtree", 3, "iab1-1-1", 0},
		{"miss", 0, "missing", 0},
		{"self on leaf", 4, "iab1-1-1", 4},
		{"parent of leaf under mid", 2, "iab1-1", 2},
		{"nested under mid", 2, "iab1-1-1", 4},
		{"sibling not under mid", 2, "iab1-2", 0},
		{"unknown start id", 99, "iab1", 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := m.MatchID(tc.id, tc.key); got != tc.want {
				t.Fatalf("MatchID(%d, %q)=%d want %d", tc.id, tc.key, got, tc.want)
			}
		})
	}
}

func TestU64StringMatcher_MatchID_NilReceiver(t *testing.T) {
	var m *U64StringMatcher
	if got := m.MatchID(0, "x"); got != 0 {
		t.Fatalf("nil MatchID=%d want 0", got)
	}
	if got := m.MatchCode(1); got != "" {
		t.Fatalf("nil MatchCode=%q want empty", got)
	}
	if got := m.ParentID(1); got != 0 {
		t.Fatalf("nil ParentID=%d want 0", got)
	}
	if got := m.AllChildrenID(1); got != nil {
		t.Fatalf("nil AllChildrenID=%v want nil", got)
	}
	if got := m.MatchAllChildrenID("x"); got != nil {
		t.Fatalf("nil MatchAllChildrenID=%v want nil", got)
	}
}

func TestU64StringMatcher_MatchCode_ParentID(t *testing.T) {
	m := sampleCategories()

	if got := m.MatchCode(2); got != "IAB1-1" {
		t.Fatalf("MatchCode(2)=%q want IAB1-1", got)
	}
	if got := m.MatchCode(99); got != "" {
		t.Fatalf("MatchCode unknown=%q want empty", got)
	}
	if got := m.ParentID(4); got != 2 {
		t.Fatalf("ParentID(4)=%d want 2", got)
	}
	if got := m.ParentID(1); got != 0 {
		t.Fatalf("ParentID(root)=%d want 0", got)
	}
	if got := m.ParentID(99); got != 0 {
		t.Fatalf("ParentID unknown=%d want 0", got)
	}
}

func TestU64StringMatcher_AllChildrenID(t *testing.T) {
	m := sampleCategories()

	want := []uint64{2, 4, 3}
	if got := m.AllChildrenID(1); !reflect.DeepEqual(got, want) {
		t.Fatalf("AllChildrenID(1)=%v want %v", got, want)
	}
	if got := m.AllChildrenID(2); !reflect.DeepEqual(got, []uint64{4}) {
		t.Fatalf("AllChildrenID(2)=%v want [4]", got)
	}
	if got := m.AllChildrenID(4); got != nil && len(got) != 0 {
		t.Fatalf("AllChildrenID(leaf)=%v want empty", got)
	}
	if got := m.AllChildrenID(99); got != nil {
		t.Fatalf("AllChildrenID unknown=%v want nil", got)
	}
	if got := m.MatchAllChildrenID("iab1"); !reflect.DeepEqual(got, want) {
		t.Fatalf("MatchAllChildrenID(iab1)=%v want %v", got, want)
	}
	if got := m.MatchAllChildrenID("IAB1-1"); !reflect.DeepEqual(got, []uint64{4}) {
		t.Fatalf("MatchAllChildrenID(IAB1-1)=%v want [4]", got)
	}
	if got := m.MatchAllChildrenID("nope"); got != nil {
		t.Fatalf("MatchAllChildrenID miss=%v want nil", got)
	}
}

func TestU64StringMatcher_ForestMultipleRoots(t *testing.T) {
	m := NewU64StringMatcher().
		SetNodes(
			Node{ID: 10, Value: "RootA"},
			Node{ParentID: 10, ID: 11, Value: "ChildA"},
			Node{ID: 20, Value: "RootB"},
			Node{ParentID: 20, ID: 21, Value: "ChildB"},
		).Freeze()

	if got := m.MatchID(0, "childb"); got != 21 {
		t.Fatalf("forest MatchID=%d want 21", got)
	}
	if got := m.MatchID(10, "childb"); got != 0 {
		t.Fatalf("cross-tree MatchID=%d want 0", got)
	}
	if got := m.MatchID(20, "ChildB"); got != 21 {
		t.Fatalf("under RootB=%d want 21", got)
	}
	if got := m.AllChildrenID(10); !reflect.DeepEqual(got, []uint64{11}) {
		t.Fatalf("AllChildrenID(RootA)=%v want [11]", got)
	}
}

func TestU64StringMatcher_ChildrenBeforeParent(t *testing.T) {
	// Insertion order: child first, then parent — Freeze must still wire the tree.
	m := NewU64StringMatcher().
		Set(1, 2, "child").
		Set(0, 1, "parent").
		Freeze()

	if got := m.MatchID(1, "child"); got != 2 {
		t.Fatalf("MatchID after late parent=%d want 2", got)
	}
	if got := m.AllChildrenID(1); !reflect.DeepEqual(got, []uint64{2}) {
		t.Fatalf("AllChildrenID=%v want [2]", got)
	}
	if got := m.ParentID(2); got != 1 {
		t.Fatalf("ParentID=%d want 1", got)
	}
}

func TestU64StringMatcher_MoveNodeReparent(t *testing.T) {
	m := NewU64StringMatcher().
		Set(0, 1, "a").
		Set(0, 2, "b").
		Set(1, 3, "c").
		Set(2, 3, "c"). // reparent 3 under 2
		Freeze()

	if got := m.ParentID(3); got != 2 {
		t.Fatalf("ParentID after move=%d want 2", got)
	}
	if got := m.AllChildrenID(1); len(got) != 0 {
		t.Fatalf("old parent children=%v want empty", got)
	}
	if got := m.AllChildrenID(2); !reflect.DeepEqual(got, []uint64{3}) {
		t.Fatalf("new parent children=%v want [3]", got)
	}
	if got := m.MatchID(1, "c"); got != 0 {
		t.Fatalf("should not match under old parent, got %d", got)
	}
	if got := m.MatchID(2, "c"); got != 3 {
		t.Fatalf("match under new parent=%d want 3", got)
	}
}

func TestU64StringMatcher_DuplicateKeyLastWins(t *testing.T) {
	m := NewU64StringMatcher().
		Set(0, 1, "Same").
		Set(0, 2, "same").
		Freeze()

	if got := m.MatchID(0, "SAME"); got != 2 {
		t.Fatalf("duplicate key last wins=%d want 2", got)
	}
	// Both nodes still exist for reverse lookups.
	if got := m.MatchCode(1); got != "Same" {
		t.Fatalf("MatchCode(1)=%q", got)
	}
	if got := m.MatchCode(2); got != "same" {
		t.Fatalf("MatchCode(2)=%q", got)
	}
}

func TestU64StringMatcher_SetOverwrite(t *testing.T) {
	m := NewU64StringMatcher().
		Set(0, 1, "A").
		Set(0, 1, "B").
		Freeze()
	if got := m.MatchID(0, "a"); got != 0 {
		t.Fatalf("old key should miss, got %d", got)
	}
	if got := m.MatchID(0, "b"); got != 1 {
		t.Fatalf("new key=%d want 1", got)
	}
	if got := m.MatchCode(1); got != "B" {
		t.Fatalf("MatchCode=%q want B", got)
	}
}

func TestU64StringMatcher_LongKeyFold(t *testing.T) {
	long := strings.Repeat("A", 80)
	m := NewU64StringMatcher().Set(0, 1, long).Freeze()
	if got := m.MatchID(0, strings.Repeat("a", 80)); got != 1 {
		t.Fatalf("long key MatchID=%d want 1", got)
	}
	if got := m.MatchID(0, strings.Repeat("a", 79)); got != 0 {
		t.Fatalf("wrong length should miss, got %d", got)
	}
}

func TestU64StringMatcher_FreezeIdempotent(t *testing.T) {
	m := NewU64StringMatcher().Set(0, 1, "x").Freeze().Freeze()
	if got := m.MatchID(0, "X"); got != 1 {
		t.Fatalf("after double Freeze=%d want 1", got)
	}
}

func TestU64StringMatcher_FreezePanicOnSet(t *testing.T) {
	m := NewU64StringMatcher().Set(0, 1, "x").Freeze()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on Set after Freeze")
		}
	}()
	m.Set(0, 2, "y")
}

func TestU64StringMatcher_SetZeroIDPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on id=0")
		}
	}()
	NewU64StringMatcher().Set(0, 0, "x")
}

func TestU64StringMatcher_AllChildrenID_SharedSlice(t *testing.T) {
	m := sampleCategories()
	a := m.AllChildrenID(1)
	b := m.AllChildrenID(1)
	if len(a) == 0 {
		t.Fatal("expected children")
	}
	// Same backing array after Freeze (read-only contract).
	if &a[0] != &b[0] {
		t.Fatal("AllChildrenID should return shared internal slice")
	}
}

func BenchmarkMatchID(b *testing.B) {
	m := sampleCategories()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if m.MatchID(1, "iab1-1-1") == 0 {
			b.Fatal("miss")
		}
	}
}

func BenchmarkMatchID_CaseFold(b *testing.B) {
	m := sampleCategories()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if m.MatchID(1, "IAB1-1-1") == 0 {
			b.Fatal("miss")
		}
	}
}

func BenchmarkMatchAllChildrenID(b *testing.B) {
	m := sampleCategories()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if len(m.MatchAllChildrenID("iab1")) == 0 {
			b.Fatal("empty")
		}
	}
}
