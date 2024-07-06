package storage

import (
	"time"

	"github.com/geniusrabbit/adcorelib/models"
)

// Reader interface collects information about advertisement objects
type Reader interface {
	CompanyList(updatedSince *time.Time) ([]*models.Company, error)
	ZoneList(updatedSince *time.Time) ([]*models.Zone, error)
	CampaignList(updatedSince *time.Time) ([]*models.Campaign, error)
	SourceList(updatedSince *time.Time) ([]*models.RTBSource, error)
	AccessPointList(updatedSince *time.Time) ([]*models.RTBAccessPoint, error)
	FormatList(updatedSince *time.Time) ([]*models.Format, error)
}

// ObjectReader interface of access to the single object
type ObjectReader interface {
	// ProjectByID(id uint64) (*models.Project, error)
	ZoneByID(id uint64) (*models.Zone, error)
	CampaignByID(id uint64) (*models.Campaign, error)
	SourceByID(id uint64) (*models.RTBSource, error)
	AccessPointByID(id uint64) (*models.RTBAccessPoint, error)
}
