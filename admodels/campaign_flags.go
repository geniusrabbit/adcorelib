package admodels

// CampaignFlagType type of the flags
type CampaignFlagType uint8

// Flags set
const (
	CampaignFlagActive  CampaignFlagType = 1 << iota // 0x01
	CampaignFlagDeleted                              //
	CampaignFlagPrivate                              // Private campaigns not avalable for publick usage
	CampaignFlagPremium
	CampaignFlagHasWeighedAds
)

//go:inline
func (f CampaignFlagType) IsActive() bool {
	return f&CampaignFlagActive != 0
}

//go:inline
func (f CampaignFlagType) IsDeleted() bool {
	return f&CampaignFlagDeleted != 0
}

//go:inline
func (f CampaignFlagType) IsPrivate() bool {
	return f&CampaignFlagPrivate != 0
}

//go:inline
func (f CampaignFlagType) IsPremium() bool {
	return f&CampaignFlagPremium != 0
}

//go:inline
func (f CampaignFlagType) HasWeighedAds() bool {
	return f&CampaignFlagHasWeighedAds != 0
}
