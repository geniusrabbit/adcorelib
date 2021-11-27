package loaders

import (
	"encoding/json"
	"io/ioutil"

	"gorm.io/gorm"

	"geniusrabbit.dev/corelib/admodels"
	"geniusrabbit.dev/corelib/admodels/types"
	"geniusrabbit.dev/corelib/models"
)

// FormatLoader from source like database of filesystem
func FormatLoader(source interface{}) func() ([]*types.Format, error) {
	switch src := source.(type) {
	case *gorm.DB:
		return func() ([]*types.Format, error) {
			return dbFormatLoader(src)
		}
	case string:
		return func() ([]*types.Format, error) {
			return fileFormatLoader(src)
		}
	}
	return nil
}

func dbFormatLoader(database *gorm.DB) (list []*types.Format, err error) {
	var formats []*models.Format
	if err = database.Find(&formats).Error; err != nil {
		return nil, err
	}
	for _, format := range formats {
		if format == nil {
			continue
		}
		list = append(list, admodels.FormatFromModel(format))
	}
	return list, err
}

type formatData struct {
	Formats []*models.Format `json:"formats"`
}

func fileFormatLoader(filename string) (list []*types.Format, _ error) {
	var (
		formats   formatData
		data, err = ioutil.ReadFile(filename)
	)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &formats); err != nil {
		return nil, err
	}
	for _, format := range formats.Formats {
		if format == nil {
			continue
		}
		list = append(list, admodels.FormatFromModel(format))
	}
	return list, err
}
