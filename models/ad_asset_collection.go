package models

import (
	"time"

	"github.com/guregu/null"
)

// -- collection & file Many2many
//
// CREATE TABLE m2m_adv_ad_file_collection
// ( asset_id              BIGINT          NOT NULL    REFERENCES  adv_ad_file(id) MATCH SIMPLE
//                                                       ON UPDATE NO ACTION
//                                                       ON DELETE CASCADE
// , collection_id         BIGINT          NOT NULL    REFERENCES  adv_ad_asset_collection(id) MATCH SIMPLE
//                                                       ON UPDATE NO ACTION
//                                                       ON DELETE CASCADE
//
// , created_at            TIMESTAMPTZ     NOT NULL  DEFAULT NOW()
//
// , PRIMARY KEY (asset_id, collection_id)
// );

// AdAssetCollection represents the list of prepared assets
type AdAssetCollection struct {
	ID        uint64 `json:"id"`
	AccountID uint64 `json:"account_id"`

	// Assets related to advertisement
	Assets []*AdAsset `json:"assets,omitempty" gorm:"many2many:m2m_adv_ad_file_collection;foreignKey:collection_id;association_autoupdate:false"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at,omitempty"`
}

// TableName in database
func (aac *AdAssetCollection) TableName() string {
	return "adv_ad_asset_collection"
}

// RBACResourceName returns the name of the resource for the RBAC
func (aac *AdAssetCollection) RBACResourceName() string {
	return "adv_ad_asset_collection"
}
