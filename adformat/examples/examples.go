// Package examples embeds the adformat/examples/*.yaml sample formats
// (§4 of the design plan) via embed.FS, so they can be loaded without a
// fragile runtime.Caller + os.Open dance (the old
// admodels/types/format_mock.go approach) and reused both by tests and
// by adformat/adformattest.
package examples

import (
	"embed"
	"fmt"

	"github.com/geniusrabbit/adcorelib/adformat"
)

//go:embed *.yaml
var FS embed.FS

// Names of every embedded example file, in the same order as presented
// in §4 of the design plan.
var Names = []string{
	"banner.yaml",
	"banner_fixed.yaml",
	"popunder.yaml",
	"native.yaml",
	"video_preroll.yaml",
	"push_ad.yaml",
	"playable.yaml",
	"multiscreen.yaml",
}

// Load reads and decodes (via adformat.DecodeAny, which also runs
// Config.Validate) the named embedded example.
func Load(name string) (*adformat.Format, error) {
	data, err := FS.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("adformat/examples: %s: %w", name, err)
	}
	f, err := adformat.DecodeAny(name, data)
	if err != nil {
		return nil, fmt.Errorf("adformat/examples: %s: %w", name, err)
	}
	return f, nil
}

// All decodes every embedded example listed in Names, in order.
func All() ([]*adformat.Format, error) {
	formats := make([]*adformat.Format, 0, len(Names))
	for _, name := range Names {
		f, err := Load(name)
		if err != nil {
			return nil, err
		}
		formats = append(formats, f)
	}
	return formats, nil
}
