package admodels

import (
	"github.com/geniusrabbit/adcorelib/admodels/types"
)

// AdAssets contains the list of file assets
type AdAssets []*AdAsset

// Main asset of the list
func (assets AdAssets) Main() *AdAsset {
	return assets.Asset(types.FormatAssetMain)
}

// Asset by name
func (assets AdAssets) Asset(name string) *AdAsset {
	isMain := name == types.FormatAssetMain
	for _, it := range assets {
		if (isMain && (it.Name == "" || it.Name == types.FormatAssetMain)) || it.Name == name {
			return it
		}
	}
	return nil
}

// AssetByID returns asset with specific ID
func (assets AdAssets) AssetByID(id uint64) *AdAsset {
	for _, it := range assets {
		if it.ID == id {
			return it
		}
	}
	return nil
}

// AssetBaner by fixed size
func (assets AdAssets) AssetBanner(w, h int) *AdAsset {
	for i, it := range assets {
		if it.Name == types.FormatAssetBanner && it.Width == w && it.Height == h {
			return assets[i]
		}
	}
	return nil
}

// Len of the assets list
func (assets AdAssets) Len() int {
	return len(assets)
}
