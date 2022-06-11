//
// @project GeniusRabbit::corelib 2018 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019
//

package models

import (
	"encoding/json"
	"time"

	"geniusrabbit.dev/corelib/admodels/types"
	"github.com/geniusrabbit/gosql/v2"
	"github.com/guregu/null"
)

type AdFileProcessingStatus int

const (
	AdFileProcessingUndefined AdFileProcessingStatus = 0
	AdFileProcessingProcessed                        = 1
	AdFileProcessingError                            = 2
	AdFileProcessingDeleted                          = 3
)

func (st AdFileProcessingStatus) Name() string {
	switch st {
	case AdFileProcessingProcessed:
		return "Processed"
	case AdFileProcessingError:
		return "Error"
	case AdFileProcessingDeleted:
		return "Deleted"
	}
	return "Undefined"
}

// AdFileThumb of the file
type AdFileThumb struct {
	Name        string            `json:"name"` // Path to image or video
	Type        types.AdAssetType `json:"type"`
	Width       int               `json:"width,omitempty"`
	Height      int               `json:"height,omitempty"`
	ContentType string            `json:"content_type,omitempty"`
	Ext         map[string]any    `json:"ext,omitempty"`
}

// Meta AdFile info
type AdAssetMeta struct {
	Main  AdFileThumb   `json:"main"`
	Items []AdFileThumb `json:"items,omitempty"`
}
type ObjectType int

// -- ad & file Many2many
//
// CREATE TABLE m2m_adv_ad_file_ad
// ( file_id               BIGINT          NOT NULL    REFERENCES  adv_ad_file(id) MATCH SIMPLE
//                                                       ON UPDATE NO ACTION
//                                                       ON DELETE CASCADE
// , ad_id                 BIGINT          NOT NULL    REFERENCES  adv_ad(id) MATCH SIMPLE
//                                                       ON UPDATE NO ACTION
//                                                       ON DELETE CASCADE
//
// , created_at            TIMESTAMPTZ     NOT NULL  DEFAULT NOW()
//
// , PRIMARY KEY (file_id, ad_id)
// );

// AdFile structure which describes the paticular file attached to advertisement
// Image advertisement: Title=Image title, Description=My description
//         ID,             HashID,                     path,  size, name,  type, content_type,                          meta
//   File:  1, dhg321h3ndp43u2hfc, 'images/a/c/banner1.jpg', 64322, NULL, image,   image/jpeg, {"main": {...}, "items": [{...}]}
//   File:  2, xxg321h3xxx43u2hfc,  'images/a/c/video1.mp4', 44322, NULL, video,  video/x-mp4, {"main": {...}, "items": [{...}]}
type AdAsset struct {
	ID        uint64   `json:"id"`
	HashID    string   `json:"hashid" gorm:"column:hashid"` // File hash
	Company   *Company `json:"company,omitempty"`           // Owner Project
	CompanyID uint64   `json:"company_id"`

	ProcessingStatus AdFileProcessingStatus `json:"processing_status"`

	ObjectID    string                              `json:"object_id"`
	FileInfo    gosql.NullableJSON[json.RawMessage] `json:"file_info"`
	Name        null.String                         `json:"name,omitempty"` // Internal file name
	ContentType string                              `json:"content_type"`
	Type        ObjectType                          `gorm:"type:INT" json:"type"`
	Meta        gosql.NullableJSON[AdAssetMeta]     `gorm:"type:JSONB" json:"meta,omitempty"`
	Size        int64                               `json:"size,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at,omitempty"`
}

// TableName in database
func (fl *AdAsset) TableName() string {
	return "adv_ad_file"
}
