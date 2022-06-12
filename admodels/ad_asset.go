//
// @project GeniusRabbit rotator 2018 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019
//

package admodels

import "geniusrabbit.dev/corelib/admodels/types"

// AdAssetThumb of the file
type AdAssetThumb struct {
	Path        string // Path to image or video
	Type        types.AdAssetType
	Width       int
	Height      int
	ContentType string
	Ext         map[string]any
}

// IsSuits thumb by size
func (th AdAssetThumb) IsSuits(w, h, wmin, hmin int) bool {
	return th.Width <= w && th.Width >= wmin && th.Height <= h && th.Height >= hmin
}

// IsImage file type
func (th AdAssetThumb) IsImage() bool {
	return th.Type.IsImage()
}

// IsVideo file type
func (th AdAssetThumb) IsVideo() bool {
	return th.Type.IsVideo()
}

// AdAsset information
type AdAsset struct {
	ID          uint64
	Name        string
	Path        string // In case of HTML5, hare must be the path to directory on CDN
	Type        types.AdAssetType
	ContentType string
	Width       int
	Height      int
	Thumbs      []AdAssetThumb
}

// // AdAssetByModel original
// func AdAssetByModel(file *models.AdAsset) *AdAsset {
// 	var (
// 		newThumbs []AdAssetThumb
// 		meta      = file.ObjectMeta()
// 	)

// 	// Prepare thumb list
// 	for _, thumb := range meta.Items {
// 		newThumbs = append(newThumbs, AdAssetThumb{
// 			Path:        urlPathJoin(file.Path, thumb.Name),
// 			Type:        AdAssetTypeByObjectType(thumb.Type),
// 			Width:       thumb.Width,
// 			Height:      thumb.Height,
// 			ContentType: thumb.ContentType,
// 		})
// 	}

// 	return &AdAsset{
// 		ID:          file.ID,
// 		Name:        file.Name.String,
// 		Path:        urlPathJoin(file.Path, meta.Main.Name),
// 		Type:        AdAssetTypeByObjectType(file.Type),
// 		ContentType: file.ContentType,
// 		Width:       meta.Main.Width,
// 		Height:      meta.Main.Height,
// 		Thumbs:      newThumbs,
// 	}
// }

// ThumbBy size borders and specific type
func (f *AdAsset) ThumbBy(w, h, wmin, hmin int, adType types.AdAssetType) (th *AdAssetThumb) {
	for i := 0; i < len(f.Thumbs); i++ {
		if f.Thumbs[i].Type == adType && f.Thumbs[i].IsSuits(w, h, wmin, hmin) {
			return &f.Thumbs[i]
		}
	}
	return nil
}

// IsImage file type
func (f *AdAsset) IsImage() bool {
	return f.Type.IsImage()
}

// IsVideo file type
func (f *AdAsset) IsVideo() bool {
	return f.Type.IsVideo()
}

// func urlPathJoin(urlBase, name string) string {
// 	if strings.HasSuffix(urlBase, "/") != strings.HasPrefix(name, "/") {
// 		return urlBase + name
// 	}
// 	return strings.TrimRight(urlBase, "/") + "/" + strings.TrimLeft(name, "/")
// }
