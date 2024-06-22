package models

// AdMeta information
type AdMeta struct {
	Content string         `json:"content,omitempty"`
	URL     string         `json:"url,omitempty"`    // For compatibility. @Deprecated
	Direct  string         `json:"direct,omitempty"` // URL @Deprecated
	Extra   map[string]any `json:"ext,omitempty"`
}
