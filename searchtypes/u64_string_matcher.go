//
// @project GeniusRabbit corelib
// @author Dmitry Ponomarev <demdxx@gmail.com>
//

package searchtypes

const asciiFoldStackMax = 64

// Node is a tree entry for U64StringMatcher.
type Node struct {
	ParentID uint64 // 0 = root
	ID       uint64
	Value    string // IAB code and/or name
}

type u64StringNode struct {
	parentID    uint64
	id          uint64
	value       string
	childIDs    []uint64
	allChildIDs []uint64
}

// U64StringMatcher is a build-once string↔uint64 tree matcher (e.g. IAB categories).
// Lifecycle: Set / SetNodes → Freeze → concurrent Match* (no locks).
type U64StringMatcher struct {
	nodes      map[uint64]*u64StringNode
	byParent   map[uint64]map[string]uint64
	childOrder map[uint64][]uint64 // parentID → direct child IDs in Set order
	byKey      map[string]uint64
	byID       map[uint64]string
	frozen     bool
}

// NewU64StringMatcher creates an empty matcher.
func NewU64StringMatcher() *U64StringMatcher {
	return &U64StringMatcher{
		nodes:      make(map[uint64]*u64StringNode),
		byParent:   make(map[uint64]map[string]uint64),
		childOrder: make(map[uint64][]uint64),
		byKey:      make(map[string]uint64),
		byID:       make(map[uint64]string),
	}
}

// Set registers or overwrites a node under parentID (0 = root).
func (m *U64StringMatcher) Set(parentID, id uint64, value string) *U64StringMatcher {
	if m.frozen {
		panic("searchtypes: U64StringMatcher.Set after Freeze")
	}
	if id == 0 {
		panic("searchtypes: U64StringMatcher.Set id must be non-zero")
	}

	lower := asciiFoldAlloc(value)

	if prev, ok := m.nodes[id]; ok {
		// Remove from previous parent's child list / byParent index.
		m.removeChild(prev.parentID, id, asciiFoldAlloc(prev.value))
	}

	n := &u64StringNode{
		parentID: parentID,
		id:       id,
		value:    value,
	}
	m.nodes[id] = n
	m.byID[id] = value
	m.byKey[lower] = id

	pm := m.byParent[parentID]
	if pm == nil {
		pm = make(map[string]uint64)
		m.byParent[parentID] = pm
	}
	pm[lower] = id

	if !containsUint64(m.childOrder[parentID], id) {
		m.childOrder[parentID] = append(m.childOrder[parentID], id)
	}

	return m
}

// SetNodes registers multiple nodes.
func (m *U64StringMatcher) SetNodes(nodes ...Node) *U64StringMatcher {
	for _, n := range nodes {
		m.Set(n.ParentID, n.ID, n.Value)
	}
	return m
}

// Freeze seals the matcher and materializes full descendant ID lists.
func (m *U64StringMatcher) Freeze() *U64StringMatcher {
	if m.frozen {
		return m
	}
	for id, n := range m.nodes {
		n.childIDs = m.childOrder[id]
		n.allChildIDs = m.collectDescendants(id)
	}
	m.frozen = true
	return m
}

// MatchID finds key in the subtree rooted at id (id itself + all descendants).
// id == 0 searches the whole forest. Returns 0 if missing / outside subtree.
func (m *U64StringMatcher) MatchID(id uint64, key string) uint64 {
	if m == nil {
		return 0
	}
	candidate := m.byKey[asciiFoldLookup(key)]
	if candidate == 0 {
		return 0
	}
	if id == 0 {
		return candidate
	}
	for cur := candidate; cur != 0; {
		if cur == id {
			return candidate
		}
		n := m.nodes[cur]
		if n == nil {
			return 0
		}
		cur = n.parentID
	}
	return 0
}

// MatchCode returns the canonical Value for id, or "".
func (m *U64StringMatcher) MatchCode(id uint64) string {
	if m == nil {
		return ""
	}
	return m.byID[id]
}

// ParentID returns the parent of id (0 for roots / unknown).
func (m *U64StringMatcher) ParentID(id uint64) uint64 {
	if m == nil {
		return 0
	}
	if n := m.nodes[id]; n != nil {
		return n.parentID
	}
	return 0
}

// AllChildrenID returns all descendant IDs of id (self not included).
// After Freeze the slice is internal/shared — do not mutate.
func (m *U64StringMatcher) AllChildrenID(id uint64) []uint64 {
	if m == nil {
		return nil
	}
	n := m.nodes[id]
	if n == nil {
		return nil
	}
	return n.allChildIDs
}

// MatchAllChildrenID resolves key globally, then returns AllChildrenID of that node.
func (m *U64StringMatcher) MatchAllChildrenID(key string) []uint64 {
	if m == nil {
		return nil
	}
	id := m.byKey[asciiFoldLookup(key)]
	if id == 0 {
		return nil
	}
	return m.AllChildrenID(id)
}

func (m *U64StringMatcher) collectDescendants(id uint64) []uint64 {
	direct := m.childOrder[id]
	if len(direct) == 0 {
		return nil
	}
	out := make([]uint64, 0, len(direct))
	var walk func(uint64)
	walk = func(cid uint64) {
		out = append(out, cid)
		for _, next := range m.childOrder[cid] {
			walk(next)
		}
	}
	for _, cid := range direct {
		walk(cid)
	}
	return out
}

func (m *U64StringMatcher) removeChild(parentID, id uint64, lower string) {
	if pm := m.byParent[parentID]; pm != nil {
		if cur, ok := pm[lower]; ok && cur == id {
			delete(pm, lower)
		}
	}
	m.childOrder[parentID] = removeUint64(m.childOrder[parentID], id)
	// Drop stale byKey only if it still points at this id.
	if m.byKey[lower] == id {
		delete(m.byKey, lower)
	}
}

func containsUint64(s []uint64, v uint64) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}

func removeUint64(s []uint64, v uint64) []uint64 {
	for i, x := range s {
		if x == v {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

// asciiFoldAlloc lowercases ASCII for storage (heap string).
func asciiFoldAlloc(s string) string {
	need := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			need = true
			break
		}
	}
	if !need {
		return s
	}
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

// asciiFoldLookup lowercases ASCII for map lookup; ≤64 bytes uses stack buffer.
func asciiFoldLookup(s string) string {
	need := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			need = true
			break
		}
	}
	if !need {
		return s
	}
	if len(s) <= asciiFoldStackMax {
		var buf [asciiFoldStackMax]byte
		for i := 0; i < len(s); i++ {
			c := s[i]
			if c >= 'A' && c <= 'Z' {
				c += 'a' - 'A'
			}
			buf[i] = c
		}
		return string(buf[:len(s)])
	}
	return asciiFoldAlloc(s)
}
