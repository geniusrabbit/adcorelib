package adtype

import (
	"time"
)

type SourceInfo struct {
	ID          string         `json:"id"`
	Protocol    string         `json:"protocol"`
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Domain      string         `json:"domain,omitempty"`
	IconURL     string         `json:"icon_url,omitempty"`
	LogoURL     string         `json:"logo_url,omitempty"`
	URL         string         `json:"url,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// SourceMinimal contains only minimal set of methods
type SourceMinimal interface {
	// Bid request for standart system filter
	Bid(request BidRequester) Response

	// ProcessResponseItem result or error
	ProcessResponseItem(Response, ResponseItem)
}

// SourceTesteChecker checker
type SourceTesteChecker interface {
	// Test current request for compatibility
	Test(request BidRequester) bool
}

// SourceTimeoutSetter interface
type SourceTimeoutSetter interface {
	// SetTimeout for sourcer
	SetTimeout(timeout time.Duration)
}

// Source of advertisement and where will be selled the traffic
type Source interface {
	SourceMinimal
	SourceTesteChecker

	// ID of the source driver
	ID() uint64

	// ObjectKey of the source driver
	ObjectKey() uint64

	// Protocol of the source driver
	Protocol() string

	// Info returns information about the source platform and the source protocol
	Info() *SourceInfo

	// PriceCorrectionReduceFactor which is a potential
	// Returns percent from 0 to 1 for reducing of the value
	// If there is 10% of price correction, it means that 10% of the final price must be ignored
	PriceCorrectionReduceFactor() float64

	// RequestStrategy description
	RequestStrategy() RequestStrategy
}

// SourceTester interface
type SourceTester interface {
	Source
	SourceTesteChecker
}
