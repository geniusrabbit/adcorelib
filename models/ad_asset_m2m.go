package models

import "time"

type AdAssetAdM2M struct {
	Ad        *Ad       `json:"ad,omitempty"`
	AdID      uint64    `gorm:"primaryKey" json:"ad_id,omitempty"`
	Asset     *AdAsset  `json:"asset,omitempty"`
	AssetID   uint64    `gorm:"primaryKey" json:"asset_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
