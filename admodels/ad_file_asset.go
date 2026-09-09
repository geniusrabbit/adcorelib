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
	URL         string                `json:"url,omitempty"`         // In case of HTML5, here must be the path to directory on CDN
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

// ClosestThumbBy returns the thumb nearest to (w, h). Thumbs that fit the
// size borders (IsSuits) are preferred; if none fit, the closest of all
// thumbs is returned. w/h <= 0 means that axis is unbounded (same as ThumbBy).
func (f *AdFileAsset) ClosestThumbBy(w, h, wmin, hmin int) *AdFileAssetThumb {
	if len(f.Thumbs) == 0 {
		return nil
	}
	tw, th := w, h
	if w <= 0 {
		w = 0x0fffffff
	}
	if h <= 0 {
		h = 0x0fffffff
	}

	var best *AdFileAssetThumb
	bestDist := int64(0)
	bestSuits := false
	for i := range f.Thumbs {
		t := &f.Thumbs[i]
		suits := t.IsSuits(w, h, wmin, hmin)
		dist := thumbSizeDist2(t, tw, th)
		if best == nil ||
			(suits && !bestSuits) ||
			(suits == bestSuits && dist < bestDist) {
			best = t
			bestDist = dist
			bestSuits = suits
		}
	}
	return best
}

func thumbSizeDist2(t *AdFileAssetThumb, tw, th int) int64 {
	var dw, dh int64
	if tw > 0 {
		dw = int64(t.Width - tw)
	}
	if th > 0 {
		dh = int64(t.Height - th)
	}
	return dw*dw + dh*dh
}

// IsImage file type
func (f *AdFileAsset) IsImage() bool {
	return f.Type.IsImage()
}

// IsVideo file type
func (f *AdFileAsset) IsVideo() bool {
	return f.Type.IsVideo()
}

// IsHTML5 file type
func (f *AdFileAsset) IsHTML5() bool {
	return f.Type.IsHTML5()
}

// IsAudio file type
func (f *AdFileAsset) IsAudio() bool {
	return f.Type.IsAudio()
}
