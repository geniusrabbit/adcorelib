//
// @project GeniusRabbit corelib
// @author Dmitry Ponomarev <demdxx@gmail.com>
//

package localecontent

import (
	"strings"
)

const (
	// DefaultFlag marks the fallback locale bucket inside AdAsset.context.
	DefaultFlag = "_default"
	// AnyLanguageCode is the legacy locale key for "any language".
	AnyLanguageCode = "**"
)

// Result is a pre-split AdAsset.context: per-locale field maps plus the
// default locale extracted once so serving does not scan for _default.
type Result struct {
	Locales       map[string]map[string]any
	DefaultLocale string
	DefaultFields map[string]any
}

// Split extracts locale buckets and the default locale from a raw context map.
// Nested context is detected when at least one top-level value is an object.
// Legacy flat maps ({title, description}) become DefaultFields with no locales.
func Split(context map[string]any) Result {
	if len(context) == 0 {
		return Result{}
	}
	if !isNested(context) {
		return Result{DefaultFields: context}
	}

	locales := make(map[string]map[string]any, len(context))
	primitives := make(map[string]any)
	var flaggedDefault, anyKey, firstKey string

	for key, value := range context {
		bucket, ok := asBucket(value)
		if !ok {
			if key != DefaultFlag {
				primitives[key] = value
			}
			continue
		}
		norm := strings.ToLower(key)
		cleaned := stripDefault(bucket)
		locales[norm] = cleaned
		if firstKey == "" || norm < firstKey {
			firstKey = norm
		}
		if hasDefaultFlag(bucket) {
			flaggedDefault = norm
		}
		if norm == AnyLanguageCode {
			anyKey = norm
		}
	}

	defaultLocale := flaggedDefault
	if defaultLocale == "" {
		defaultLocale = anyKey
	}
	if defaultLocale == "" {
		if _, ok := locales["en"]; ok {
			defaultLocale = "en"
		} else {
			defaultLocale = firstKey
		}
	}

	defaultFields := locales[defaultLocale]
	if len(primitives) > 0 {
		defaultFields = merge(primitives, defaultFields)
	}

	return Result{
		Locales:       locales,
		DefaultLocale: defaultLocale,
		DefaultFields: defaultFields,
	}
}

// Fields returns the flattened field map for lang (merged with default when
// a non-default locale matches). Empty lang uses the default locale.
func Fields(context map[string]any, lang string) map[string]any {
	return FieldsOf(Split(context), lang)
}

// Item returns a single field from the resolved locale map, then a top-level
// primitive on the raw context (legacy iframe_url / content).
func Item(context map[string]any, lang, name string) any {
	return ItemOf(Split(context), context, lang, name)
}

// FieldsOf resolves lang against an already-split Result.
func FieldsOf(r Result, lang string) map[string]any {
	key := langKey(lang)
	if key == "" || len(r.Locales) == 0 {
		return r.DefaultFields
	}
	matched := key
	bucket, ok := r.Locales[matched]
	if !ok {
		if prefix := languagePrefix(key); prefix != key {
			matched = prefix
			bucket, ok = r.Locales[matched]
		}
	}
	if !ok || matched == r.DefaultLocale {
		return r.DefaultFields
	}
	return merge(r.DefaultFields, bucket)
}

// ItemOf looks up name in FieldsOf, then in raw top-level primitives.
func ItemOf(r Result, context map[string]any, lang, name string) any {
	if fields := FieldsOf(r, lang); fields != nil {
		if v, ok := fields[name]; ok {
			return v
		}
	}
	if context != nil {
		if v, ok := context[name]; ok {
			if _, isBucket := asBucket(v); !isBucket {
				return v
			}
		}
	}
	return nil
}

func langKey(lang string) string {
	return strings.ToLower(strings.TrimSpace(lang))
}

func languagePrefix(lang string) string {
	if i := strings.IndexAny(lang, "-_"); i > 0 {
		return lang[:i]
	}
	return lang
}

func isNested(context map[string]any) bool {
	for _, v := range context {
		if _, ok := asBucket(v); ok {
			return true
		}
	}
	return false
}

func asBucket(v any) (map[string]any, bool) {
	m, ok := v.(map[string]any)
	if !ok || m == nil {
		return nil, false
	}
	return m, true
}

func hasDefaultFlag(bucket map[string]any) bool {
	v, ok := bucket[DefaultFlag]
	if !ok {
		return false
	}
	b, ok := v.(bool)
	return ok && b
}

func stripDefault(bucket map[string]any) map[string]any {
	if _, has := bucket[DefaultFlag]; !has {
		return bucket
	}
	out := make(map[string]any, len(bucket)-1)
	for k, v := range bucket {
		if k == DefaultFlag {
			continue
		}
		out[k] = v
	}
	return out
}

func merge(base, overlay map[string]any) map[string]any {
	if len(overlay) == 0 {
		return base
	}
	if len(base) == 0 {
		return overlay
	}
	out := make(map[string]any, len(base)+len(overlay))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range overlay {
		out[k] = v
	}
	return out
}
