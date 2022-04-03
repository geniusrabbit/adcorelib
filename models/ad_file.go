//
// @project GeniusRabbit::corelib 2018 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019
//

package models

import (
	"time"

	"github.com/geniusrabbit/gosql"
	"github.com/guregu/null"
)

// AdFileThumb of the file
type AdFileThumb struct {
	Path        string // Path to image or video
	Type        int
	Width       int
	Height      int
	ContentType string
}

// Meta AdFile info
type Meta struct {
	Thumbs []AdFileThumb `json:"thumbs"`
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
// , created_at            TIMESTAMPTZ       NOT NULL  DEFAULT NOW()
//
// , PRIMARY KEY (file_id, ad_id)
// );

// AdFile structure which describes the paticular file attached to advertisement
// Image advertisement: Title=Image title, Description=My description
//         ID,             HashID,                     path,  size, name,  type, content_type,                          meta
//   File:  1, dhg321h3ndp43u2hfc, 'images/a/c/banner1.jpg', 64322, NULL, image,   image/jpeg, {"main": {...}, "items": [{...}]}
//   File:  2, xxg321h3xxx43u2hfc,  'images/a/c/video1.mp4', 44322, NULL, video,  video/x-mp4, {"main": {...}, "items": [{...}]}
type AdFile struct {
	ID        uint64   `json:"id"`
	HashID    string   `json:"hashid" gorm:"column:hashid"` // File hash
	Company   *Company `json:"company,omitempty"`           // Owner Project
	CompanyID uint64   `json:"company_id"`

	Path        string             `json:"path"`
	Name        null.String        `json:"name,omitempty"` // Internal file name
	ContentType string             `json:"content_type"`
	Type        ObjectType         `gorm:"type:INT" json:"type"`
	Meta        gosql.NullableJSON `gorm:"type:JSONB" json:"meta,omitempty"`
	Size        int64              `json:"size,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at,omitempty"`
}

// TableName in database
func (fl *AdFile) TableName() string {
	return "adv_ad_file"
}

// ObjectMeta information of file object
func (fl AdFile) ObjectMeta() (meta *Meta) {
	meta = new(Meta)
	fl.Meta.UnmarshalTo(meta)
	return meta
}
