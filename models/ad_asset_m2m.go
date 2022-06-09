package models

import "time"

type AdAssetAdM2M struct {
	Ad        *Ad       `json:"ad,omitempty"`
	AdID      uint64    `gorm:"column:ad_id" json:"ad_id,omitempty"`
	Asset     *AdFile   `json:"asset,omitempty"`
	AssetID   uint64    `gorm:"column:file_id" json:"asset_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
