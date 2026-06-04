//
// @project GeniusRabbit corelib 2018 - 2019, 2025
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019, 2025
//

package admodels

import (
	"github.com/geniusrabbit/adcorelib/admodels/types"
)

// AdFileAssetThumb of the file
type AdFileAssetThumb struct {
	URL    string                `json:"url"`
	Width  int                   `json:"w"`
	Height int                   `json:"h"`
	Type   types.AdFileAssetType `json:"type,omitempty"`
}

// IsSuits thumb by size
func (th AdFileAssetThumb) IsSuits(w, h, wmin, hmin int) bool {
	return th.Width <= w && th.Width >= wmin && th.Height <= h && th.Height >= hmin
}

// IsImage file type
func (th *AdFileAssetThumb) IsImage() bool {
	return th.Type.IsImage()
}

// IsVideo file type
func (th *AdFileAssetThumb) IsVideo() bool {
	return th.Type.IsVideo()
}

// AdFileAsset information
type AdFileAsset struct {
	ID          uint64                `json:"id,omitempty"`
	ExternalID  string                `json:"external_id,omitempty"` // ID of the asset in the source system, like VAST media file ID
	Name        string                `json:"name,omitempty"`        // Name of the asset, like "main", "banner", "icon", etc.
	AltText     string                `json:"alt_text,omitempty"`    // Alternative text for the asset, used for accessibility and as a fallback
	URL         string                `json:"url,omitempty"`         // In case of HTML5, hare must be the path to directory on CDN
	Type        types.AdFileAssetType `json:"type,omitempty"`
	ContentType string                `json:"content_type,omitempty"`
	Width       int                   `json:"width,omitempty"`
	Height      int                   `json:"height,omitempty"`
	Duration    int                   `json:"duration,omitempty"` // Duration in seconds, for video assets
	Thumbs      []AdFileAssetThumb    `json:"thumbs,omitempty"`
}

// ThumbBy size borders and specific type
func (f *AdFileAsset) ThumbBy(w, h, wmin, hmin int) (th *AdFileAssetThumb) {
	if w <= 0 {
		w = 0x0fffffff
	}
	if h <= 0 {
		h = 0x0fffffff
	}
	for i := 0; i < len(f.Thumbs); i++ {
		if f.Thumbs[i].IsSuits(w, h, wmin, hmin) {
			if th == nil || th.Width > f.Thumbs[i].Width {
				th = &f.Thumbs[i]
			}
		}
	}
	return th
}

// IsImage file type
func (f *AdFileAsset) IsImage() bool {
	return f.Type.IsImage()
}

// IsVideo file type
func (f *AdFileAsset) IsVideo() bool {
	return f.Type.IsVideo()
}
