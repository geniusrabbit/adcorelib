//
// @project GeniusRabbit corelib
// @author Dmitry Ponomarev <demdxx@gmail.com>
//

package searchtypes

import "sort"

// buildKwNode is used only while constructing the trie (before Freeze).
type buildKwNode struct {
	id       uint64
	phraseLen int // length of normalized phrase ending here (0 if not accepting)
	children map[byte]*buildKwNode
}

type kwNode struct {
	id        uint64
	phraseLen uint16 // normalized phrase length for longest-match; 0 if not accepting
	edgeOff   uint32
	edgeCount uint16
}

type kwEdge struct {
	b     byte
	child uint32
}

// U64KeywordMatcher maps keyword phrases to uint64 IDs with word-start prefix matching.
// Lifecycle: Set → Freeze → concurrent Match (no locks).
//
// Matching: from each whitespace-delimited word start in the input, walk the trie.
// Intermediate phrase words must match exactly (then whitespace); the final phrase
// word may be a prefix of the input word (e.g. "racing" matches "racings").
// Among hits, the longest normalized phrase wins; ties prefer the first word-start.
type U64KeywordMatcher struct {
	root   *buildKwNode
	nodes  []kwNode
	edges  []kwEdge
	frozen bool
}

// NewU64KeywordMatcher creates an empty matcher.
func NewU64KeywordMatcher() *U64KeywordMatcher {
	return &U64KeywordMatcher{
		root: &buildKwNode{children: make(map[byte]*buildKwNode)},
	}
}

// Set registers or overwrites a phrase → id mapping. Empty / whitespace-only phrases are skipped.
func (m *U64KeywordMatcher) Set(phrase string, id uint64) *U64KeywordMatcher {
	if m.frozen {
		panic("searchtypes: U64KeywordMatcher.Set after Freeze")
	}
	if id == 0 {
		return m
	}
	norm := normalizeKeywordPhrase(phrase)
	if norm == "" {
		return m
	}
	n := m.root
	for i := 0; i < len(norm); i++ {
		b := norm[i]
		child := n.children[b]
		if child == nil {
			child = &buildKwNode{children: make(map[byte]*buildKwNode)}
			n.children[b] = child
		}
		n = child
	}
	n.id = id
	n.phraseLen = len(norm)
	return m
}

// Freeze seals the matcher into contiguous node/edge arenas for Match.
func (m *U64KeywordMatcher) Freeze() *U64KeywordMatcher {
	if m == nil || m.frozen {
		return m
	}
	m.nodes = m.nodes[:0]
	m.edges = m.edges[:0]
	m.nodes = append(m.nodes, kwNode{}) // index 0 = root
	m.freezeNode(m.root, 0)
	m.root = nil
	m.frozen = true
	return m
}

func (m *U64KeywordMatcher) freezeNode(bn *buildKwNode, idx uint32) {
	n := &m.nodes[idx]
	n.id = bn.id
	if bn.phraseLen > 0 {
		n.phraseLen = uint16(bn.phraseLen)
	}
	if len(bn.children) == 0 {
		return
	}
	keys := make([]byte, 0, len(bn.children))
	for b := range bn.children {
		keys = append(keys, b)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	n.edgeOff = uint32(len(m.edges))
	n.edgeCount = uint16(len(keys))
	base := len(m.edges)
	for range keys {
		m.edges = append(m.edges, kwEdge{})
	}
	for i, b := range keys {
		childIdx := uint32(len(m.nodes))
		m.nodes = append(m.nodes, kwNode{})
		m.edges[base+i] = kwEdge{b: b, child: childIdx}
		m.freezeNode(bn.children[b], childIdx)
	}
}

// Match finds the best category ID for input, or 0.
func (m *U64KeywordMatcher) Match(input string) uint64 {
	if m == nil || !m.frozen || len(m.nodes) == 0 {
		return 0
	}
	folded := asciiFoldLookup(input)
	// Trim leading/trailing ASCII space without allocating when possible.
	start, end := 0, len(folded)
	for start < end && isASCIISpace(folded[start]) {
		start++
	}
	for end > start && isASCIISpace(folded[end-1]) {
		end--
	}
	if start >= end {
		return 0
	}
	s := folded[start:end]

	var bestID uint64
	var bestLen uint16

	for i := 0; i < len(s); {
		if i == 0 || isASCIISpace(s[i-1]) {
			if !isASCIISpace(s[i]) {
				id, plen := m.matchFrom(s, i)
				if id != 0 && (bestID == 0 || plen > bestLen) {
					bestID = id
					bestLen = plen
				}
			}
		}
		i++
	}
	return bestID
}

func (m *U64KeywordMatcher) matchFrom(s string, pos int) (uint64, uint16) {
	nodeIdx := uint32(0)
	i := pos
	var bestID uint64
	var bestLen uint16

	for {
		n := &m.nodes[nodeIdx]
		if n.id != 0 && (bestID == 0 || n.phraseLen > bestLen) {
			bestID = n.id
			bestLen = n.phraseLen
		}
		if n.edgeCount == 0 || i >= len(s) {
			return bestID, bestLen
		}

		// Prefer consuming a space edge when input has whitespace (multi-word phrases).
		if isASCIISpace(s[i]) {
			if child, ok := m.findChild(n, ' '); ok {
				for i < len(s) && isASCIISpace(s[i]) {
					i++
				}
				nodeIdx = child
				continue
			}
			return bestID, bestLen
		}

		child, ok := m.findChild(n, s[i])
		if !ok {
			return bestID, bestLen
		}
		i++
		nodeIdx = child
	}
}

func (m *U64KeywordMatcher) findChild(n *kwNode, b byte) (uint32, bool) {
	edges := m.edges[n.edgeOff : n.edgeOff+uint32(n.edgeCount)]
	// Binary search over sorted edges.
	lo, hi := 0, len(edges)
	for lo < hi {
		mid := (lo + hi) / 2
		eb := edges[mid].b
		if eb < b {
			lo = mid + 1
		} else if eb > b {
			hi = mid
		} else {
			return edges[mid].child, true
		}
	}
	return 0, false
}

func isASCIISpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == '\f' || b == '\v'
}

// normalizeKeywordPhrase ASCII-folds, trims, and collapses internal whitespace to a single space.
func normalizeKeywordPhrase(s string) string {
	folded := asciiFoldAlloc(s)
	if folded == "" {
		return ""
	}
	b := make([]byte, 0, len(folded))
	spacePending := false
	started := false
	for i := 0; i < len(folded); i++ {
		c := folded[i]
		if isASCIISpace(c) {
			if started {
				spacePending = true
			}
			continue
		}
		if spacePending {
			b = append(b, ' ')
			spacePending = false
		}
		b = append(b, c)
		started = true
	}
	return string(b)
}
