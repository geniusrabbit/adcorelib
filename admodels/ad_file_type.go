//
// @project GeniusRabbit rotator 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package admodels

// import (
// 	diskmodels "geniusrabbit.dev/disk/models"
// )

// AdFileType type
type AdFileType uint

// AdFileType values
const (
	AdFileUndefinedType AdFileType = 0
	AdFileImageType     AdFileType = 1
	AdFileVideoType     AdFileType = 2
	AdFileHTML5Type     AdFileType = 3
)

// AdFileTypeByName returns adfile value type
func AdFileTypeByName(name string) AdFileType {
	switch name {
	case "image", "img":
		return AdFileImageType
	case "video":
		return AdFileVideoType
	case "html5":
		return AdFileHTML5Type
	}
	return AdFileUndefinedType
}

// // AdFileTypeByObjectType value
// func AdFileTypeByObjectType(tp diskmodels.ObjectType) AdFileType {
// 	switch tp {
// 	case diskmodels.TypeImage:
// 		return AdFileImageType
// 	case diskmodels.TypeVideo:
// 		return AdFileVideoType
// 	case diskmodels.TypeHTMLArchType:
// 		return AdFileHTML5Type
// 	}
// 	return AdFileUndefinedType
// }

func (ft AdFileType) String() string {
	switch ft {
	case AdFileImageType:
		return "image" // diskmodels.TypeImage.String()
	case AdFileVideoType:
		return "video" //diskmodels.TypeVideo.String()
	case AdFileHTML5Type:
		return "html5" //diskmodels.TypeHTMLArchType.String()
	}
	return "undefined" //diskmodels.TypeUndefined.String()
}

// IsImage file type
func (ft AdFileType) IsImage() bool {
	return ft == AdFileImageType
}

// IsVideo file type
func (ft AdFileType) IsVideo() bool {
	return ft == AdFileVideoType
}
