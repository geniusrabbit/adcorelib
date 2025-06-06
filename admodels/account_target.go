//
// @project GeniusRabbit corelib 2017, 2019, 2024
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019, 2024
//

package admodels

// AccountTarget wrapper for replace of epsent target object
type AccountTarget struct {
	Acc *Account
}

// ID of object (Zone OR SmartLink only)
func (c AccountTarget) ID() uint64 {
	return 0
}

// Size default of target item
func (c AccountTarget) Size() (w, h int) {
	return w, h
}

// CommissionShareFactor of current target
func (c AccountTarget) CommissionShareFactor() float64 {
	return c.Acc.CommissionShareFactor()
}

// Account object
func (c AccountTarget) Account() *Account {
	return c.Acc
}

// ProjectID number
func (c AccountTarget) ProjectID() uint64 {
	return 0
}

// AccountID of current target
func (c AccountTarget) AccountID() uint64 {
	return c.Acc.ID()
}
