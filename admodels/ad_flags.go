package admodels

// AdFlag option state type
type AdFlag uint8

// Ad flag option types
const (
	AdFlagIsPremium AdFlag = 0x01
	AdFlagActive    AdFlag = 0x02
	AdFlagInsecure  AdFlag = 0x04
	AdFlagAsPopover AdFlag = 0x08
)

//go:inline
func (af AdFlag) IsPremium() bool {
	return af&AdFlagIsPremium != 0
}

//go:inline
func (af AdFlag) IsActive() bool {
	return af&AdFlagActive != 0
}

//go:inline
func (af AdFlag) IsInsecure() bool {
	return af&AdFlagInsecure != 0
}

//go:inline
func (af AdFlag) IsPopover() bool {
	return af&AdFlagAsPopover != 0
}
