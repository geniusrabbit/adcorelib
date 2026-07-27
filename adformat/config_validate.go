package adformat

import "fmt"

// Reserved NavigateTo commands (§3.2.11) — not screen names.
const (
	NavigateBack  = "$back"
	NavigateClose = "$close"
)

// Validate walks the whole Config tree once and checks everything that
// could otherwise fail silently at first use (§3.2.6):
//  1. Every Field.RegExp compiles.
//  2. Every Condition.Field (on a Field/AssetRequirement/nested Config/
//     Param) references an existing sibling name at the same tree level.
//  3. No cycles in Condition dependencies at a given level.
//  4. Every NavigateTo (other than the reserved "$back"/"$close") and
//     every screen_ref Field.Options value references an existing
//     screen (named nested group) — the set of screen names is the
//     nearest enclosing level of group siblings: a field/asset nested
//     directly inside one named screen (e.g. "intro") targets the
//     screen's own siblings ("offer", ...), not something inside
//     "intro" itself — see the worked example in multiscreen.yaml
//     (§3.2.11).
//  5. Every AssetRequirement.AspectRatio parses as "W:H", and
//     BitrateMin does not exceed BitrateMax when both are set
//     (§3.2.1 gaps #1/#4).
//
// Call this once when a Format/Config is loaded (codec.go does this
// automatically for LoadYAML/LoadJSON), not on every use.
func (c Config) Validate() error {
	return c.validate("<root>", nil)
}

func (c Config) validate(path string, inheritedScreenNames map[string]struct{}) error {
	names := make(map[string]struct{}, len(c.Children))
	ownGroupNames := make(map[string]struct{})
	for _, n := range c.Children {
		if name := n.Name(); name != "" {
			names[name] = struct{}{}
			if n.Kind == NodeGroup {
				ownGroupNames[name] = struct{}{}
			}
		}
	}

	// The screen namespace for this level's own fields/assets/nested
	// groups: this level's own group children if it has any (this level
	// IS the set of screens), otherwise whatever was inherited from an
	// enclosing level (this level is the content of one screen, whose
	// siblings were defined further up the tree).
	groupNames := inheritedScreenNames
	if len(ownGroupNames) > 0 {
		groupNames = ownGroupNames
	}

	for _, n := range c.Children {
		childPath := path + "." + n.Name()
		switch n.Kind {
		case NodeAsset:
			if n.Asset == nil {
				return fmt.Errorf("adformat: %s: asset node has a nil AssetRequirement", childPath)
			}
			if err := validateConditionRefs(n.Asset.Condition, names, childPath); err != nil {
				return err
			}
			if err := validateNavigateTo(n.Asset.NavigateTo, groupNames, childPath); err != nil {
				return err
			}
			if n.Asset.AspectRatio != "" {
				if _, _, err := ParseAspectRatio(n.Asset.AspectRatio); err != nil {
					return fmt.Errorf("adformat: %s: %w", childPath, err)
				}
			}
			if n.Asset.BitrateMin > 0 && n.Asset.BitrateMax > 0 && n.Asset.BitrateMin > n.Asset.BitrateMax {
				return fmt.Errorf("adformat: %s: bitrate_min (%d) is greater than bitrate_max (%d)", childPath, n.Asset.BitrateMin, n.Asset.BitrateMax)
			}
		case NodeField:
			if n.Field == nil {
				return fmt.Errorf("adformat: %s: field node has a nil Field", childPath)
			}
			if n.Field.RegExp != "" {
				if _, err := compileRegexp(n.Field.RegExp); err != nil {
					return fmt.Errorf("adformat: %s: invalid regexp: %w", childPath, err)
				}
			}
			if err := validateConditionRefs(n.Field.Condition, names, childPath); err != nil {
				return err
			}
			if err := validateNavigateTo(n.Field.NavigateTo, groupNames, childPath); err != nil {
				return err
			}
			if n.Field.Type == FieldScreenRefType {
				for _, opt := range n.Field.Options {
					screen := fmt.Sprint(opt.Value)
					if _, ok := groupNames[screen]; !ok {
						return fmt.Errorf("adformat: %s: screen_ref option %q does not reference an existing screen", childPath, screen)
					}
				}
			}
		case NodeGroup:
			if n.Group == nil {
				return fmt.Errorf("adformat: %s: group node has a nil Config", childPath)
			}
			if err := validateConditionRefs(n.Group.Condition, names, childPath); err != nil {
				return err
			}
			if err := n.Group.validate(childPath, groupNames); err != nil {
				return err
			}
		case NodeParam:
			if n.Param == nil {
				return fmt.Errorf("adformat: %s: param node has a nil Param", childPath)
			}
			if err := validateConditionRefs(n.Param.Condition, names, childPath); err != nil {
				return err
			}
		default:
			return fmt.Errorf("adformat: %s: unknown node kind %q", childPath, n.Kind)
		}
	}

	if err := detectConditionCycles(c.Children); err != nil {
		return fmt.Errorf("adformat: %s: %w", path, err)
	}
	return nil
}

func validateConditionRefs(cond *Condition, names map[string]struct{}, path string) error {
	if cond == nil {
		return nil
	}
	if cond.Field != "" {
		if _, ok := names[cond.Field]; !ok {
			return fmt.Errorf("adformat: %s: condition references unknown field/asset %q", path, cond.Field)
		}
	}
	for i := range cond.All {
		if err := validateConditionRefs(&cond.All[i], names, path); err != nil {
			return err
		}
	}
	for i := range cond.Any {
		if err := validateConditionRefs(&cond.Any[i], names, path); err != nil {
			return err
		}
	}
	return validateConditionRefs(cond.Not, names, path)
}

func validateNavigateTo(target string, screenNames map[string]struct{}, path string) error {
	if target == "" || target == NavigateBack || target == NavigateClose {
		return nil
	}
	if _, ok := screenNames[target]; !ok {
		return fmt.Errorf("adformat: %s: navigate_to references unknown screen %q", path, target)
	}
	return nil
}

// detectConditionCycles reports a cyclic Condition dependency among the
// named siblings of one Config level (A depends on B, B depends on A) via
// a plain DFS — the dependency graph of one format's fields is expected
// to be small and shallow.
func detectConditionCycles(children []Node) error {
	graph := make(map[string][]string, len(children))
	for _, n := range children {
		name := n.Name()
		if name == "" {
			continue
		}
		graph[name] = collectConditionFields(n.Condition())
	}

	const (
		stateUnvisited = 0
		stateVisiting  = 1
		stateDone      = 2
	)
	state := make(map[string]int, len(graph))

	var visit func(node string) error
	visit = func(node string) error {
		switch state[node] {
		case stateVisiting:
			return fmt.Errorf("cyclic condition dependency involving %q", node)
		case stateDone:
			return nil
		}
		state[node] = stateVisiting
		for _, dep := range graph[node] {
			if err := visit(dep); err != nil {
				return err
			}
		}
		state[node] = stateDone
		return nil
	}

	for node := range graph {
		if err := visit(node); err != nil {
			return err
		}
	}
	return nil
}

func collectConditionFields(cond *Condition) []string {
	if cond == nil {
		return nil
	}
	var fields []string
	if cond.Field != "" {
		fields = append(fields, cond.Field)
	}
	for i := range cond.All {
		fields = append(fields, collectConditionFields(&cond.All[i])...)
	}
	for i := range cond.Any {
		fields = append(fields, collectConditionFields(&cond.Any[i])...)
	}
	fields = append(fields, collectConditionFields(cond.Not)...)
	return fields
}

// Walk calls fn for every Node in the tree rooted at c, passing a stable
// dotted path (e.g. "panel.headline", "offer.hotspot.target_screen")
// usable for error messages and for addressing a specific node from a UI
// editor. Recurses into nested groups depth-first, in Children order.
func (c Config) Walk(fn func(path string, n Node) error) error {
	return c.walk("", fn)
}

func (c Config) walk(prefix string, fn func(path string, n Node) error) error {
	for _, n := range c.Children {
		path := n.Name()
		if prefix != "" {
			path = prefix + "." + path
		}
		if err := fn(path, n); err != nil {
			return err
		}
		if n.Kind == NodeGroup && n.Group != nil {
			if err := n.Group.walk(path, fn); err != nil {
				return err
			}
		}
	}
	return nil
}
