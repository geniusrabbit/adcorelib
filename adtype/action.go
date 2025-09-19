package adtype

// Action type
type Action int

// Campaign actions
const (
	ActionImpression Action = 1
	ActionImp               = ActionImpression // Alias
	ActionView       Action = 2
	ActionClick      Action = 3
	ActionLead       Action = 4
)

func (a Action) String() string {
	switch a {
	case ActionImpression:
		return "impression"
	case ActionView:
		return "view"
	case ActionClick:
		return "click"
	case ActionLead:
		return "lead"
	}
	return "undefined"
}

// Int value of action
//
//go:inline
func (a Action) Int() int { return int(a) }

// IsImpression action type
//
//go:inline
func (a Action) IsImpression() bool { return a == ActionImpression }

// IsView action type
//
//go:inline
func (a Action) IsView() bool { return a == ActionView }

// IsClick action type
//
//go:inline
func (a Action) IsClick() bool { return a == ActionClick }

// IsLead action type
//
//go:inline
func (a Action) IsLead() bool { return a == ActionLead }

// LeadAcceptCoef delimiter magic value
const (
	LeadAcceptCoef = 100
)
