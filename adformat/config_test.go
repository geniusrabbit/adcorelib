package adformat

import "testing"

// Ported from admodels/types/format_meta_test.go (TestConfigIntersec),
// rewritten against the new builder API.
func TestConfigIntersec(t *testing.T) {
	tests := []struct {
		name   string
		c1, c2 Config
		result bool
	}{
		{
			name:   "empty",
			c1:     Config{},
			c2:     Config{},
			result: true,
		},
		{
			name:   "basic_stretch",
			c1:     Config{}.WithAssets(Asset("main").Require()),
			c2:     Config{}.WithAssets(Asset("main")),
			result: true,
		},
		{
			name: "basic_assets_similar_sizes",
			c1: Config{}.WithAssets(
				Asset("main").Require().Size(100, 1000, 100, 1000),
				Asset("icon").Size(0, 80, 0, 80),
			),
			c2: Config{}.WithAssets(
				Asset("main").Require().Size(150, 1100, 200, 900),
				Asset("icon").Require().Size(50, 100, 50, 100),
			),
			result: true,
		},
		{
			name: "basic_assets_similar_sizes_negative",
			c1: Config{}.WithAssets(
				Asset("main").Require().Size(100, 1000, 100, 1000),
				Asset("icon").Require().Size(50, 100, 50, 100),
			),
			c2: Config{}.WithAssets(
				Asset("main").Require().Size(150, 1100, 200, 900),
				Asset("icon").Size(0, 80, 0, 80),
			),
			result: false,
		},
		{
			name: "basic_ext",
			c1: Config{}.
				WithAssets(Asset("main").Require().ExactSize(300, 100)).
				WithFields(StringField("title").Require()),
			c2: Config{}.
				WithAssets(Asset("main").ExactSize(300, 100), Asset("icon").ExactSize(30, 30)).
				WithFields(StringField("title")),
			result: true,
		},
		{
			name:   "basic_field_negative",
			c1:     Config{}.WithFields(StringField("title").Require()),
			c2:     Config{},
			result: false,
		},
		{
			name: "basic_field_negative2",
			c1:   Config{}.WithFields(StringField("title").Require()),
			c2: Config{}.WithFields(
				StringField("title").Require(),
				StringField("icon").Require(),
			),
			result: false,
		},
		{
			name: "basic_field_positive",
			c1: Config{}.WithFields(
				StringField("title").Require(),
				StringField("icon"),
			),
			c2: Config{}.WithFields(
				StringField("title").Require(),
				StringField("icon").Require(),
			),
			result: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c1.Intersec(tt.c2); got != tt.result {
				t.Errorf("Intersec() = %v, want %v", got, tt.result)
			}
		})
	}
}

func TestConfigAssetAndFieldAccessors(t *testing.T) {
	cfg := Config{}.
		WithAssets(
			Asset("main").Require().ExactSize(300, 100),
			Asset("icon").ExactSize(30, 30),
		).
		WithFields(
			StringField("title").Require(),
			StringField("description"),
		)

	if got := len(cfg.Assets()); got != 2 {
		t.Errorf("len(Assets()) = %d, want 2", got)
	}
	if got := len(cfg.Fields()); got != 2 {
		t.Errorf("len(Fields()) = %d, want 2", got)
	}
	if cfg.AssetByName("icon") == nil {
		t.Error("expected to find asset by name 'icon'")
	}
	if cfg.AssetByName("") == nil {
		t.Error("expected empty name to resolve to the main asset")
	}
	if cfg.MainAsset() == nil {
		t.Error("expected MainAsset to be found")
	}
	if cfg.GetField("title") == nil {
		t.Error("expected to find field by name 'title'")
	}
	if cfg.GetField("does-not-exist") != nil {
		t.Error("expected nil for unknown field name")
	}
	if cfg.RequiredField() == nil {
		t.Error("expected a required field to be found")
	}
	if cfg.RequiredFieldExcept("title") != nil {
		t.Error("expected no required field left after excepting 'title'")
	}
	if cfg.IsEmpty() {
		t.Error("expected a non-empty config")
	}
	if !(Config{}).IsEmpty() {
		t.Error("expected a zero-value config to be empty")
	}
}

func TestConfigSimpleAsset(t *testing.T) {
	simple := Config{}.WithAssets(Asset("main").Require())
	if simple.SimpleAsset() == nil {
		t.Error("expected SimpleAsset to be found when main is the only required asset")
	}

	withAnotherRequired := Config{}.WithAssets(
		Asset("main").Require(),
		Asset("icon").Require(),
	)
	if withAnotherRequired.SimpleAsset() != nil {
		t.Error("expected SimpleAsset to be nil when another asset is also required")
	}
}

func TestConfigGetParam(t *testing.T) {
	cfg := Config{}.WithParams(NewParam("layout", "v2"))
	p, ok := cfg.GetParam("layout")
	if !ok || p.Value != "v2" {
		t.Errorf("GetParam(layout) = (%+v, %v), want v2/true", p, ok)
	}
	if _, ok := cfg.GetParam("missing"); ok {
		t.Error("expected GetParam to report false for a missing param")
	}
}

func TestConfigActiveFieldsAndAssets(t *testing.T) {
	cfg := Config{}.
		WithFields(
			StringField("adtype").WithOptions(Opt("business", ""), Opt("game", "")),
			PhoneField("phone").When(FieldIn("adtype", "business")),
		)

	active := cfg.ActiveFields(map[string]any{"adtype": "business"})
	if len(active) != 2 {
		t.Errorf("expected both fields active for adtype=business, got %d", len(active))
	}

	inactive := cfg.ActiveFields(map[string]any{"adtype": "game"})
	if len(inactive) != 1 {
		t.Errorf("expected only adtype active for adtype=game, got %d", len(inactive))
	}
}

func TestConfigEntryGroup(t *testing.T) {
	cfg := Config{}.WithGroups(
		Group("intro").AsEntry(),
		Group("offer"),
	)
	entry, ok := cfg.EntryGroup()
	if !ok || entry.Name != "intro" {
		t.Errorf("EntryGroup() = (%q, %v), want (intro, true)", entry.Name, ok)
	}

	noExplicitEntry := Config{}.WithGroups(Group("first"), Group("second"))
	entry, ok = noExplicitEntry.EntryGroup()
	if !ok || entry.Name != "first" {
		t.Errorf("EntryGroup() with no explicit Entry = (%q, %v), want (first, true)", entry.Name, ok)
	}

	noGroups := Config{}
	if _, ok := noGroups.EntryGroup(); ok {
		t.Error("expected EntryGroup to report false when there are no nested groups")
	}
}

func TestConfigWalk(t *testing.T) {
	cfg := Config{}.WithGroups(
		Group("panel").WithFields(StringField("headline")),
	)

	var visited []string
	err := cfg.Walk(func(path string, n Node) error {
		visited = append(visited, path)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(visited) != 2 {
		t.Fatalf("expected 2 visited nodes, got %d: %v", len(visited), visited)
	}
	if visited[0] != "panel" || visited[1] != "panel.headline" {
		t.Errorf("unexpected paths: %v", visited)
	}
}
