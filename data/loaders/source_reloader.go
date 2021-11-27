package loaders

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/hashicorp/hcl"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"

	"geniusrabbit.dev/corelib/admodels"
	"geniusrabbit.dev/corelib/adtype"
	"geniusrabbit.dev/corelib/context/ctxlogger"
	"geniusrabbit.dev/corelib/eventtraking/eventstream"
	"geniusrabbit.dev/corelib/models"
)

var errInvalidFileFormat = errors.New(`[SourceAccessLoader] invalid file format`)

type (
	companyGetter func(id uint64) *admodels.Company
	subReloader   func(sourceList []*models.RTBSource, logger *zap.Logger) (sources []adtype.Source, err error)
	Platformer    interface {
		New(src *admodels.RTBSource, options ...interface{}) (adtype.Source, error)
	}
)

// SourceAccessLoader provides object of elements reloading
type SourceAccessLoader struct {
	CompanyGetter companyGetter
	EventStream   eventstream.Stream
	WinNotify     nc.Publisher
	Platforms     map[string]Platformer
}

// SourceReloader accessor
func (sl *SourceAccessLoader) Reloader(ctx context.Context, dataSource interface{}) func() ([]adtype.Source, error) {
	logger := ctxlogger.Get(ctx)
	switch d := dataSource.(type) {
	case *gorm.DB:
		return _DBSourceReloader(logger, d, sl.reload)
	case string:
		return _FSSourceReloader(logger, d, sl.reload)
	}
	return nil
}

// DBSourceReloader accessor
func _DBSourceReloader(logger *zap.Logger, database *gorm.DB, reloader subReloader) func() ([]adtype.Source, error) {
	logger = logger.With(
		zap.String("module", "SourceAccessLoader"),
		zap.String("datasrc", "database"),
	)
	return func() ([]adtype.Source, error) {
		var sourceList []*models.RTBSource
		if err := database.Find(&sourceList).Error; err != nil {
			return nil, err
		}
		return reloader(sourceList, logger)
	}
}

// FSSourceReloader accessor
func _FSSourceReloader(logger *zap.Logger, filename string, reloader subReloader) func() ([]adtype.Source, error) {
	logger = logger.With(
		zap.String("module", "SourceAccessLoader"),
		zap.String("datasrc", "fs"),
	)
	return func() (sources []adtype.Source, err error) {
		sourceList, err := readtypes(filename)
		if err != nil {
			return
		}
		newSourceList := make([]*models.RTBSource, 0, len(sourceList))
		for _, src := range sourceList {
			if src.Active.IsActive() && src.Status.IsApproved() {
				newSourceList = append(newSourceList, src)
			}
		}
		return reloader(newSourceList, logger)
	}
}

func (sl *SourceAccessLoader) reload(sourceList []*models.RTBSource, logger *zap.Logger) (sources []adtype.Source, err error) {
	for _, baseSource := range sourceList {
		var (
			err    error
			src    = admodels.RTBSourceFromModel(baseSource, sl.CompanyGetter(baseSource.CompanyID))
			source adtype.Source
		)

		if src == nil {
			logger.Error("invalid RTB source object",
				zap.Uint64("source_id", baseSource.ID),
				zap.String("protocol", baseSource.Protocol))
			continue
		}

		if fact := sl.Platforms[src.Protocol]; fact != nil {
			source, err = fact.New(src,
				eventstream.WinNotifications(sl.WinNotify),
				sl.EventStream,
				logger.With(
					zap.Uint64("platform_id", src.ID),
					zap.String("protocol", src.Protocol),
				))
		} else {
			logger.Error("invalid RTB client",
				zap.Uint64("source_id", src.ID),
				zap.String("protocol", src.Protocol))
			continue
		}

		if err != nil {
			logger.Error("invalid RTB source",
				zap.Uint64("source_id", src.ID),
				zap.String("protocol", src.Protocol),
				zap.Error(err))
		} else {
			sources = append(sources, source)
		}
	}
	return sources, err
}

func readtypes(filename string) ([]*models.RTBSource, error) {
	type sourcesDataInfo struct {
		Sources []*models.RTBSource `json:"sources" yaml:"sources"`
	}
	var (
		sourcesData sourcesDataInfo
		data, err   = ioutil.ReadFile(filename)
	)
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".json":
		err = json.Unmarshal(data, &sourcesData)
	case ".yml", ".yaml":
		err = yaml.Unmarshal(data, &sourcesData)
	case ".xml":
		err = xml.Unmarshal(data, &sourcesData)
	case ".hcl":
		err = hcl.Unmarshal(data, &sourcesData)
	default:
		err = errInvalidFileFormat
	}
	return sourcesData.Sources, err
}
