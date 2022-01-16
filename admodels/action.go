//
// @project GeniusRabbit rotator 2016 – 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019
//

package admodels

// Action type
type Action int

func (a Action) String() string {
	switch a {
	case ActionImpression:
		return "impression"
	case ActionClick:
		return "click"
	case ActionLead:
		return "lead"
	}
	return "undefined"
}

// Int value of action
func (a Action) Int() int {
	return int(a)
}

// IsImpression action type
func (a Action) IsImpression() bool {
	return a == ActionImpression
}

// IsClick action type
func (a Action) IsClick() bool {
	return a == ActionClick
}

// IsLead action type
func (a Action) IsLead() bool {
	return a == ActionLead
}

// Campaign actions
const (
	ActionImpression Action = 1
	ActionView       Action = 1
	ActionClick      Action = 2
	ActionLead       Action = 3
)

// LeadAcceptCoef delimiter magic value
const (
	LeadAcceptCoef = 100
)
