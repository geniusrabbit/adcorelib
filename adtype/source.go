package adtype

import (
	"path/filepath"
	"strings"
	"time"
)

// SourceInfo contains information about the source platform and the source protocol
// It contains the name, description, domain, icon URL, logo URL, URL, DSP domains, and metadata
// DSP domains is a list of domains that are used to serve DSP traffic
// Metadata is a map of metadata about the source platform
// It is used to store additional information about the source platform
type SourceInfo struct {
	ID          string         `json:"id"`
	Protocol    string         `json:"protocol"`
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Domain      string         `json:"domain,omitempty"`
	IconURL     string         `json:"icon_url,omitempty"`
	LogoURL     string         `json:"logo_url,omitempty"`
	URL         string         `json:"url,omitempty"`
	DSPDomains  []string       `json:"dsp_domains,omitempty"` // ["*.dsp.domain.com"]
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// IsDSPDomain checks if the domain is a DSP domain
// It checks if the domain is a suffix of the DSP domain or if the domain matches the DSP domain pattern
func (s *SourceInfo) IsDSPDomain(domain string) bool {
	for _, dspDomain := range s.DSPDomains {
		if strings.HasSuffix(domain, dspDomain) {
			return true
		}
		if matched, _ := filepath.Match(dspDomain, domain); matched {
			return true
		}
	}
	return false
}

// SourceInfoList is a list of SourceInfo
type SourceInfoList []*SourceInfo

// SourceInfoByID returns the SourceInfo by ID
func (l SourceInfoList) SourceInfoByDSPDomain(domain string) *SourceInfo {
	for _, source := range l {
		if source.IsDSPDomain(domain) {
			return source
		}
	}
	return nil
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
	// Test current request for compatibility.
	// Returns a typed cause on rejection, or nil when the request may proceed.
	Test(request BidRequester) error
}

// SourceTimeoutSetter interface
type SourceTimeoutSetter interface {
	// SetTimeout for sourcer
	SetTimeout(timeout time.Duration)
}

// Source of advertisement and where will be selled the traffic
type Source interface {
	SourceMinimal

	// AccountID of the source driver
	AccountID() uint64

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
