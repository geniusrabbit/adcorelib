//
// @project GeniusRabbit AdNet
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package models

import (
	"time"

	"github.com/guregu/null"

	"geniusrabbit.dev/corelib/admodels/types"
	"geniusrabbit.dev/corelib/billing"
)

// ```pg
// CREATE TABLE company_m2m_member
// ( company_id          BIGINT                      NOT NULL
// , user_id             BIGINT                      NOT NULL
//
// , is_admin            BOOL                        NOT NULL      DEFAULT FALSE -- for current company
// , acl                 JSONB                                     DEFAULT NULL  -- {model:flags,@custom:value}
// , roles               BIGINT[]                                  DEFAULT NULL
//
// , created_at          TIMESTAMPTZ                 NOT NULL      DEFAULT NOW()
// , updated_at          TIMESTAMPTZ                 NOT NULL      DEFAULT NOW()
// , deleted_at          TIMESTAMPTZ
//
// , PRIMARY KEY (user_id, company_id)
// );
// ```

// Company model
type Company struct {
	ID          uint64              `json:"id"`                                           //
	Name        string              `json:"name"`                                         // Unique project name. Like login
	Title       string              `json:"title"`                                        //
	Description string              `json:"description"`                                  //
	Status      types.ApproveStatus `json:"status"`                                       //
	Members     []*User             `gorm:"many2many:company_m2m_member;" json:"members"` // Members of project
	CompanyName string              `json:"company_name"`                                 // Company info
	Country     string              `json:"country"`                                      // - // -
	City        string              `json:"city"`                                         // - // -
	Address     string              `json:"address"`                                      // - // -
	Phone       string              `json:"phone"`                                        // Contacts
	Email       string              `json:"email"`                                        // - // -
	Messanger   string              `json:"messanger"`                                    // - // -

	MaxDaily billing.Money `json:"max_daily,omitempty"`
	// RevenueShare it's amount of percent of the raw incode which will be shared with the publisher company
	// For example:
	//   Displayed ads for 100$
	//   Company revenue share 60%
	//   In such case the ad network have 40$
	//   The publisher have 60$
	RevenueShare float64 `json:"revenue_share,omitempty"` // % 100_00, 10000 -> 100%, 6550 -> 65.5%

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
}

// TableName in database
func (c *Company) TableName() string {
	return "company"
}
