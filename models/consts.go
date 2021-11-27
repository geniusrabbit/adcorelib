//
// @project geniusrabbit::corelib 2016 - 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 - 2018
//

package models

// ApproveStatus type
type ApproveStatus uint

// Status approve
const (
	StatusPending      ApproveStatus = 0
	StatusApproved     ApproveStatus = 1
	StatusRejected     ApproveStatus = 2
	StatusPendingName                = `pending`
	StatusApprovedName               = `approved`
	StatusRejectedName               = `rejected`
)

// Name of the status
func (st ApproveStatus) Name() string {
	switch st {
	case StatusApproved:
		return StatusApprovedName
	case StatusRejected:
		return StatusRejectedName
	}
	return StatusPendingName
}

// IsApproved status of the object
func (st ApproveStatus) IsApproved() bool {
	return st == StatusApproved
}

// ApproveNameToStatus name to const
func ApproveNameToStatus(name string) ApproveStatus {
	switch name {
	case StatusApprovedName:
		return StatusApproved
	case StatusRejectedName:
		return StatusRejected
	}
	return StatusPending
}

// ActiveStatus of the object
type ActiveStatus uint

// Status active
const (
	StatusPause  ActiveStatus = 0
	StatusActive ActiveStatus = 1
)

// Name of the status
func (st ActiveStatus) Name() string {
	if st == StatusActive {
		return `active`
	}
	return `pause`
}

// IsActive status of the object
func (st ActiveStatus) IsActive() bool {
	return st == StatusActive
}

// Status private
const (
	StatusPublic  = 0
	StatusPrivate = 1
)
