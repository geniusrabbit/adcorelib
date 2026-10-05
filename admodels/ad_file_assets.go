package admodels

import (
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/rand"
)

// AdFileAssets contains the list of file assets
type AdFileAssets []*AdFileAsset

// RandomMain returns a random main asset
func (assets AdFileAssets) RandomMain(assetType types.AdFileAssetType) *AdFileAsset {
	return assets.RandomAsset(types.FormatAssetMain, assetType)
}

// RandomMainAsset returns a random main asset
func (assets AdFileAssets) RandomMainAsset(assetType types.AdFileAssetType) *AdFileAsset {
	return assets.RandomAsset(types.FormatAssetMain, assetType)
}

// RandomAsset returns a random asset by name (main or other)
func (assets AdFileAssets) RandomAsset(name string, assetType types.AdFileAssetType) *AdFileAsset {
	if len(assets) == 0 {
		return nil
	}
	isMain := name == types.FormatAssetMain
	n := len(assets)
	start := rand.FastPositiveIntn(n)
	for i := 0; i < n; i++ {
		asset := assets[(start+i)%n]
		if asset != nil && (isMain && (asset.Name == "" || asset.Name == types.FormatAssetMain)) || asset.Name == name {
			if assetType != types.AdFileAssetUndefinedType && asset.Type != assetType {
				continue
			}
			return asset
		}
	}
	return nil
}

// RandomAssetWithSize returns a random asset by name and size (main or other)
func (assets AdFileAssets) RandomAssetWithSize(name string, assetType types.AdFileAssetType, w, h int) *AdFileAsset {
	if len(assets) == 0 {
		return nil
	}
	n := len(assets)
	start := rand.FastPositiveIntn(n)
	for i := 0; i < n; i++ {
		asset := assets[(start+i)%n]
		if asset != nil && asset.Name == name && asset.Width == w && asset.Height == h {
			if assetType != types.AdFileAssetUndefinedType && asset.Type != assetType {
				continue
			}
			return asset
		}
	}
	return nil
}

// HasAsset returns true if the assets list has an asset with the given name
func (assets AdFileAssets) HasAsset(name string, assetType types.AdFileAssetType) bool {
	return assets.RandomAsset(name, assetType) != nil
}

// Len returns the length of the assets list
func (assets AdFileAssets) Len() int {
	return len(assets)
}
