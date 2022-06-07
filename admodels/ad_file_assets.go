package admodels

import (
	"geniusrabbit.dev/corelib/admodels/types"
)

// AdFileAssets contains the list of file assets
type AdFileAssets []*AdFile

// Main asset of the list
func (assets AdFileAssets) Main() *AdFile {
	return assets.Asset(types.FormatAssetMain)
}

// Asset by name
func (assets AdFileAssets) Asset(name string) *AdFile {
	isMain := name == types.FormatAssetMain
	for _, it := range assets {
		if (isMain && (it.Name == "" || it.Name == types.FormatAssetMain)) || it.Name == name {
			return it
		}
	}
	return nil
}

// AssetByID returns asset with specific ID
func (assets AdFileAssets) AssetByID(id uint64) *AdFile {
	for _, it := range assets {
		if it.ID == id {
			return it
		}
	}
	return nil
}

// AssetBaner by fixed size
func (assets AdFileAssets) AssetBanner(w, h int) *AdFile {
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
