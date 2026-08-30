package localecontent

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func nestedFixture() map[string]any {
	return map[string]any{
		"en": map[string]any{
			DefaultFlag: true,
			"title":     "Hello",
			"cta":       "Buy",
		},
		"ru": map[string]any{
			"title": "Привет",
		},
	}
}

func TestSplitNested(t *testing.T) {
	r := Split(nestedFixture())
	assert.Equal(t, "en", r.DefaultLocale)
	assert.Equal(t, "Hello", r.DefaultFields["title"])
	assert.Equal(t, "Buy", r.DefaultFields["cta"])
	require.NotNil(t, r.Locales["en"])
	require.NotNil(t, r.Locales["ru"])
	_, hasFlag := r.Locales["en"][DefaultFlag]
	assert.False(t, hasFlag)
	assert.Equal(t, "Привет", r.Locales["ru"]["title"])
}

func TestSplitFlatLegacy(t *testing.T) {
	flat := map[string]any{"title": "Hello", "description": "Sale"}
	r := Split(flat)
	assert.Empty(t, r.DefaultLocale)
	assert.Nil(t, r.Locales)
	assert.Equal(t, flat, r.DefaultFields)
}

func TestSplitAnyLanguageFallback(t *testing.T) {
	r := Split(map[string]any{
		AnyLanguageCode: map[string]any{"title": "Any"},
		"de":            map[string]any{"title": "Hallo"},
	})
	assert.Equal(t, AnyLanguageCode, r.DefaultLocale)
	assert.Equal(t, "Any", r.DefaultFields["title"])
}

func TestSplitPrefersEnWithoutFlag(t *testing.T) {
	r := Split(map[string]any{
		"de": map[string]any{"title": "Hallo"},
		"en": map[string]any{"title": "Hello"},
	})
	assert.Equal(t, "en", r.DefaultLocale)
}

func TestSplitEmpty(t *testing.T) {
	r := Split(nil)
	assert.Empty(t, r.DefaultLocale)
	assert.Nil(t, r.DefaultFields)
	r = Split(map[string]any{})
	assert.Nil(t, r.DefaultFields)
}

func TestFieldsLocaleMerge(t *testing.T) {
	ctx := nestedFixture()
	en := Fields(ctx, "en")
	assert.Equal(t, "Hello", en["title"])
	assert.Equal(t, "Buy", en["cta"])
	_, hasFlag := en[DefaultFlag]
	assert.False(t, hasFlag)

	ru := Fields(ctx, "ru")
	assert.Equal(t, "Привет", ru["title"])
	assert.Equal(t, "Buy", ru["cta"])

	de := Fields(ctx, "de")
	assert.Equal(t, "Hello", de["title"])
	assert.Equal(t, "Buy", de["cta"])
}

func TestFieldsLanguagePrefix(t *testing.T) {
	ctx := nestedFixture()
	assert.Equal(t, "Hello", Fields(ctx, "en-US")["title"])
	assert.Equal(t, "Привет", Fields(ctx, "ru_RU")["title"])
}

func TestFieldsEmptyLangUsesDefault(t *testing.T) {
	ctx := nestedFixture()
	assert.Equal(t, "Hello", Fields(ctx, "")["title"])
	assert.Equal(t, "Hello", Item(ctx, "", "title"))
}

func TestFieldsFlatLegacy(t *testing.T) {
	flat := map[string]any{"title": "Hello"}
	assert.Equal(t, "Hello", Fields(flat, "ru")["title"])
	assert.Equal(t, "Hello", Item(flat, "en", "title"))
}

func TestItemTopLevelPrimitive(t *testing.T) {
	ctx := map[string]any{
		"en": map[string]any{
			DefaultFlag: true,
			"title":     "Hello",
		},
		"iframe_url": "https://example.com/frame",
	}
	assert.Equal(t, "https://example.com/frame", Item(ctx, "en", "iframe_url"))
	assert.Equal(t, "https://example.com/frame", Item(ctx, "ru", "iframe_url"))
	assert.Equal(t, "https://example.com/frame", Fields(ctx, "en")["iframe_url"])
}

func TestItemMissing(t *testing.T) {
	assert.Nil(t, Item(nestedFixture(), "en", "missing"))
	assert.Nil(t, Item(nil, "en", "title"))
}
