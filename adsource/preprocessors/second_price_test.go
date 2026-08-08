//
// @project GeniusRabbit corelib 2026
// @author Dmitry Ponomarev <demdxx@gmail.com>
//

package preprocessors

import (
	"context"
	"iter"
	"testing"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

type stubAd struct {
	adtype.ResponseItemEmpty

	auctionCPM  billing.Money
	price       billing.Money
	bidFloorCPM billing.Money // Impression.BidFloorCPM
	purchaseImp billing.Money // Impression.PurchaseImpPrice → MinimumPrice when > 0
	second      *adtype.SecondAd

	setCalled bool
	setPrice  billing.Money
	withComm  bool
}

func (s *stubAd) InternalAuctionCPMBid() billing.Money { return s.auctionCPM }

func (s *stubAd) Price(action adtype.Action) billing.Money {
	if action.IsImpression() {
		return s.price
	}
	return 0
}

func (s *stubAd) SetBidPrice(action adtype.Action, price billing.Money, withCommission bool) error {
	if !action.IsImpression() {
		return adtype.ErrUnsupportedAction
	}
	s.setCalled = true
	s.setPrice = price
	s.withComm = withCommission
	s.price = price
	return nil
}

func (s *stubAd) Second() *adtype.SecondAd { return s.second }

func (s *stubAd) Impression() *adtype.Impression {
	return &adtype.Impression{
		BidFloorCPM:      s.bidFloorCPM,
		PurchaseImpPrice: s.purchaseImp,
	}
}

type stubResponse struct {
	items []adtype.ResponseItemCommon
}

func (r *stubResponse) AuctionID() string                     { return "" }
func (r *stubResponse) AuctionType() types.AuctionType        { return types.UndefinedAuctionType }
func (r *stubResponse) Source() adtype.Source                 { return nil }
func (r *stubResponse) Request() adtype.BidRequester          { return nil }
func (r *stubResponse) Ads() []adtype.ResponseItemCommon      { return r.items }
func (r *stubResponse) Item(string) adtype.ResponseItemCommon { return nil }
func (r *stubResponse) Count() int                            { return len(r.items) }
func (r *stubResponse) Validate() error                       { return nil }
func (r *stubResponse) Error() error                          { return nil }
func (r *stubResponse) Context(...context.Context) context.Context {
	return context.Background()
}
func (r *stubResponse) Get(string) any { return nil }
func (r *stubResponse) Release()       {}

func (r *stubResponse) IterAds() iter.Seq[adtype.ResponseItem] {
	return func(yield func(adtype.ResponseItem) bool) {
		for _, it := range r.items {
			if ad, ok := it.(adtype.ResponseItem); ok {
				if !yield(ad) {
					return
				}
			}
		}
	}
}

func newResp(ads ...*stubAd) adtype.Response {
	items := make([]adtype.ResponseItemCommon, len(ads))
	for i, ad := range ads {
		items[i] = ad
	}
	return &stubResponse{items: items}
}

func TestSecondPrice_NilOrEmpty(t *testing.T) {
	p := SecondPrice{}

	got, err := p.PreprocessResponse(nil)
	if err != nil || got != nil {
		t.Fatalf("nil: got=%v err=%v", got, err)
	}

	empty := &stubResponse{}
	got, err = p.PreprocessResponse(empty)
	if err != nil || got != empty {
		t.Fatalf("empty: got=%v err=%v", got, err)
	}
}

func TestSecondPrice_Single_BidFloorHalf(t *testing.T) {
	// BidFloorCPM=200_000 → per-imp floor 200; midpoint with 1000 = 600
	ad := &stubAd{
		auctionCPM:  5000,
		price:       1000,
		bidFloorCPM: 200_000,
	}
	_, err := SecondPrice{}.PreprocessResponse(newResp(ad))
	if err != nil {
		t.Fatal(err)
	}
	want := billing.Money(200 + (1000-200)/2) // 600
	if !ad.setCalled || ad.setPrice != want || !ad.withComm {
		t.Fatalf("setCalled=%v price=%d want=%d withComm=%v", ad.setCalled, ad.setPrice, want, ad.withComm)
	}
}

func TestSecondPrice_Single_NoBidFloor_UsesMinimumPrice(t *testing.T) {
	ad := &stubAd{
		auctionCPM:  5000,
		price:       1000,
		purchaseImp: 200,
	}
	_, err := SecondPrice{}.PreprocessResponse(newResp(ad))
	if err != nil {
		t.Fatal(err)
	}
	if !ad.setCalled || ad.setPrice != 200 || !ad.withComm {
		t.Fatalf("setCalled=%v price=%d want=200 withComm=%v", ad.setCalled, ad.setPrice, ad.withComm)
	}
}

func TestSecondPrice_Single_NoFloors_TenPercent(t *testing.T) {
	ad := &stubAd{auctionCPM: 5000, price: 1000}
	_, err := SecondPrice{}.PreprocessResponse(newResp(ad))
	if err != nil {
		t.Fatal(err)
	}
	if !ad.setCalled || ad.setPrice != 100 || !ad.withComm {
		t.Fatalf("setCalled=%v price=%d want=100 withComm=%v", ad.setCalled, ad.setPrice, ad.withComm)
	}
}

func TestSecondPrice_Single_UsesSecond(t *testing.T) {
	ad := &stubAd{
		auctionCPM:  5000,
		price:       1000,
		purchaseImp: 200,
		second:      &adtype.SecondAd{Price: 350},
	}
	_, err := SecondPrice{}.PreprocessResponse(newResp(ad))
	if err != nil {
		t.Fatal(err)
	}
	if !ad.setCalled || ad.setPrice != 350 || !ad.withComm {
		t.Fatalf("setCalled=%v price=%d want=350 withComm=%v", ad.setCalled, ad.setPrice, ad.withComm)
	}
}

func TestSecondPrice_Single_ZeroPrice_Skip(t *testing.T) {
	ad := &stubAd{auctionCPM: 1, price: 0}
	_, err := SecondPrice{}.PreprocessResponse(newResp(ad))
	if err != nil {
		t.Fatal(err)
	}
	if ad.setCalled {
		t.Fatal("expected SetBidPrice not called for zero cleared price")
	}
}

func TestSecondPrice_GSP_ThreeAds(t *testing.T) {
	a := &stubAd{auctionCPM: 9000, price: 900}
	b := &stubAd{auctionCPM: 5000, price: 500}
	c := &stubAd{auctionCPM: 3000, price: 300, purchaseImp: 100}

	_, err := SecondPrice{}.PreprocessResponse(newResp(c, a, b))
	if err != nil {
		t.Fatal(err)
	}

	if !a.setCalled || a.setPrice != 500 || !a.withComm {
		t.Fatalf("a: setCalled=%v price=%d want=500 withComm=%v", a.setCalled, a.setPrice, a.withComm)
	}
	if !b.setCalled || b.setPrice != 300 || !b.withComm {
		t.Fatalf("b: setCalled=%v price=%d want=300 withComm=%v", b.setCalled, b.setPrice, b.withComm)
	}
	// no BidFloor → MinimumPrice (purchaseImp=100)
	if !c.setCalled || c.setPrice != 100 || !c.withComm {
		t.Fatalf("c: setCalled=%v price=%d want=100 withComm=%v", c.setCalled, c.setPrice, c.withComm)
	}
}

func TestSecondPrice_Tail_TenPercentCappedByPrevious(t *testing.T) {
	a := &stubAd{auctionCPM: 9000, price: 1000, second: &adtype.SecondAd{Price: 50}}
	b := &stubAd{auctionCPM: 1000, price: 800}
	// b: 10% of 800 = 80, capped by previous cleared (a=50) → 50

	_, err := SecondPrice{}.PreprocessResponse(newResp(a, b))
	if err != nil {
		t.Fatal(err)
	}
	if a.setPrice != 50 {
		t.Fatalf("a price=%d want=50", a.setPrice)
	}
	if b.setPrice != 50 {
		t.Fatalf("b price=%d want=50 (capped by previous)", b.setPrice)
	}
}

func TestSecondPrice_Tail_BidFloorHalfCapped(t *testing.T) {
	a := &stubAd{auctionCPM: 9000, price: 1000, second: &adtype.SecondAd{Price: 50}}
	b := &stubAd{auctionCPM: 1000, price: 800, bidFloorCPM: 100_000} // floor=100; half=450 → cap 50

	_, err := SecondPrice{}.PreprocessResponse(newResp(a, b))
	if err != nil {
		t.Fatal(err)
	}
	if b.setPrice != 50 {
		t.Fatalf("b price=%d want=50 (capped)", b.setPrice)
	}
}

func TestSecondPrice_EmptySecondIgnored(t *testing.T) {
	a := &stubAd{auctionCPM: 9000, price: 900, second: &adtype.SecondAd{}}
	b := &stubAd{auctionCPM: 5000, price: 400}

	_, err := SecondPrice{}.PreprocessResponse(newResp(a, b))
	if err != nil {
		t.Fatal(err)
	}
	if a.setPrice != 400 {
		t.Fatalf("a should use next price, got %d", a.setPrice)
	}
}

var _ adtype.ResponseItem = (*stubAd)(nil)
