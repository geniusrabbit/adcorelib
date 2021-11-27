package loaders

import (
	"encoding/json"
	"io/ioutil"

	"gorm.io/gorm"

	"geniusrabbit.dev/corelib/admodels"
	"geniusrabbit.dev/corelib/models"
)

// TargetReloader function
func TargetReloader(source interface{}) func() ([]admodels.Target, error) {
	switch src := source.(type) {
	case *gorm.DB:
		return func() ([]admodels.Target, error) {
			return dbTargetReloader(src)
		}
	case string:
		return func() ([]admodels.Target, error) {
			return fileTargetReloader(src)
		}
	}
	return nil
}

func dbTargetReloader(database *gorm.DB) (list []admodels.Target, err error) {
	var zones []*models.Zone
	if err = database.Find(&zones).Error; err != nil {
		return nil, err
	}
	for _, zone := range zones {
		if zone == nil {
			continue
		}
		list = append(list, admodels.TargetFromModel(*zone))
	}
	return list, err
}

type targetData struct {
	Zones []*models.Zone `json:"zones"`
}

func fileTargetReloader(filename string) (list []admodels.Target, _ error) {
	var (
		targets   targetData
		data, err = ioutil.ReadFile(filename)
	)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &targets); err != nil {
		return nil, err
	}
	for _, zone := range targets.Zones {
		if zone == nil {
			continue
		}
		list = append(list, admodels.TargetFromModel(*zone))
	}
	return list, err
}
