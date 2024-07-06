//
// @project GeniusRabbit corelib 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package auction

import (
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

type equipmentSlice []adtype.ResponserItemCommon

func (l equipmentSlice) Len() int      { return len(l) }
func (l equipmentSlice) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l equipmentSlice) Less(i, j int) bool {
	e1, _ := l[i].(adtype.ResponserMultipleItem)
	e2, _ := l[j].(adtype.ResponserMultipleItem)
	if e1 != nil && (e2 == nil || e1.Count() > e2.Count()) {
		return true
	}
	if e2 != nil && (e1 == nil || e1.Count() < e2.Count()) {
		return false
	}
	return avgBid(e1, l[i]) > avgBid(e2, l[j])
}

func avgBid(mit adtype.ResponserMultipleItem, it adtype.ResponserItemCommon) billing.Money {
	if mit != nil {
		return mit.AuctionCPMBid() / billing.Money(mit.Count())
	}
	return it.AuctionCPMBid()
}
