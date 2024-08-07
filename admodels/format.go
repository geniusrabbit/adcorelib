//
// @project GeniusRabbit corelib 2016 – 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2018
//

package admodels

import (
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/models"
)

// Formats field set
var (
	FormatVideoFieldSet = []string{
		types.FormatFieldTitle,
		types.FormatFieldBrandname,
		types.FormatFieldURL,
		types.FormatFieldStartDisplay,
	}
	FormatNativeFieldSet = []string{
		types.FormatFieldTitle,
		types.FormatFieldDescription,
		types.FormatFieldBrandname,
		types.FormatFieldURL,
		types.FormatFieldPhone,
	}
)

// FormatFromModel object
func FormatFromModel(format *models.Format) *types.Format {
	var (
		fmtp types.FormatTypeBitset
		conf = format.Config.Data
	)

	switch {
	case format.IsDirect():
		fmtp.Set(types.FormatDirectType)
	case format.IsProxy() && conf.RequiredFieldExcept(models.FormatFieldContent) == nil:
		fmtp.Set(types.FormatProxyType)
	default:
		if asset := conf.SimpleAsset(); asset != nil {
			if asset.IsHTML5Support() {
				if conf.RequiredFieldExcept() == nil {
					fmtp.Set(types.FormatBannerHTML5Type)
				}
			}

			// Detect the simple banner asset
			if asset.IsFixed() && conf.RequiredFieldExcept(models.FormatFieldTitle) == nil {
				if asset.IsImageSupport() || asset.IsVideoSupport() {
					fmtp.Set(types.FormatBannerType)
				}
			}

			// Detect the video integrated into player
			if asset.IsVideoSupport() &&
				conf.RequiredFieldExcept(FormatVideoFieldSet...) == nil &&
				conf.GetField(models.FormatFieldStartDisplay) != nil {
				fmtp.Set(types.FormatVideoType)
			}

			// Detect native advertisement
			if asset.IsImageSupport() || asset.IsVideoSupport() {
				if conf.RequiredFieldExcept(FormatNativeFieldSet...) == nil {
					// Must be required (brandname || title) && description
					if (conf.GetField(models.FormatFieldTitle) != nil ||
						conf.GetField(models.FormatFieldBrandname) != nil) &&
						conf.GetField(models.FormatFieldDescription) != nil {
						fmtp.Set(types.FormatNativeType)
					}
				}
			} // end if
		}
	}

	return &types.Format{
		ID:        format.ID,
		Codename:  format.Codename,
		Types:     fmtp,
		Width:     format.Width,
		Height:    format.Height,
		MinWidth:  format.MinWidth,
		MinHeight: format.MinHeight,
		Config:    conf,
	}
}
