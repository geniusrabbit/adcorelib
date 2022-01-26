package models

import (
	"time"

	"github.com/geniusrabbit/gosql"
	"github.com/guregu/null"
)

// CompanyM2MMember model
type CompanyM2MMember struct {
	CompanyID int      `gorm:"primaryKey" json:"company_id"`
	Company   *Company `gorm:"foreignkey:Company" json:"company"`
	UserID    int      `gorm:"primaryKey" json:"user_id"`
	User      *User    `gorm:"foreignkey:User" json:"user"`

	IsAdmin bool                    `json:"is_admin"`
	ACL     gosql.NullableJSON      `gorm:"type:JSONB" json:"acl"`
	Roles   gosql.NullableUintArray `gorm:"type:BIGINT[]" json:"roles"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
}

// TableName in database
func (c *CompanyM2MMember) TableName() string {
	return "company_m2m_member"
}
