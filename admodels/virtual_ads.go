//
// @project GeniusRabbit corelib 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package admodels

// VirtualAds extract for targeting
type VirtualAds struct {
	Campaign *Campaign
	Bids     []TargetBid
}
