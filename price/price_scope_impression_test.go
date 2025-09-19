package price

import "testing"

func TestPriceScopeImpression(t *testing.T) {
	ps := &PriceScopeImpression{
		MaxBidImpPrice: 100,
	}
	if !ps.SetBidImpressionPrice(50, false) {
		t.Error("SetBidImpressionPrice failed for valid price")
	}
	if ps.BidImpPrice != 50 {
		t.Errorf("Expected BidImpPrice to be 50, got %v", ps.BidImpPrice)
	}

	if ps.SetBidImpressionPrice(150, false) {
		t.Error("SetBidImpressionPrice should fail for price exceeding MaxBidImpPrice")
	}

	if !ps.SetBidImpressionPrice(0, true) {
		t.Error("SetBidImpressionPrice failed for zero price with maxIfZero true")
	}
}
