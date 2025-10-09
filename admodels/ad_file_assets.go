package admodels

import (
	"github.com/geniusrabbit/adcorelib/admodels/types"
)

// AdFileAssets contains the list of file assets
type AdFileAssets []*AdFileAsset

// Main asset of the list
func (assets AdFileAssets) Main() *AdFileAsset {
	return assets.Asset(types.FormatAssetMain)
}

// Asset by name
func (assets AdFileAssets) Asset(name string) *AdFileAsset {
	isMain := name == types.FormatAssetMain
	for _, it := range assets {
		if (isMain && (it.Name == "" || it.Name == types.FormatAssetMain)) || it.Name == name {
			return it
		}
	}
	return nil
}

// AssetBanner by fixed size
func (assets AdFileAssets) AssetBanner(w, h int) *AdFileAsset {
	for i, it := range assets {
		if it.Name == types.FormatAssetBanner && it.Width == w && it.Height == h {
			return assets[i]
		}
	}
	return nil
}

// Len of the assets list
func (assets AdFileAssets) Len() int {
	return len(assets)
}
