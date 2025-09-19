package price

import (
	"testing"

	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

func TestPriceScope_PricePerAction(t *testing.T) {
	tests := []struct {
		name       string
		priceScope PriceScope
		actionType adtype.Action
		expected   billing.Money
	}{
		{
			name: "ActionView returns ViewPrice",
			priceScope: PriceScope{
				ViewPrice:  billing.MoneyFloat(1.5),
				ClickPrice: billing.MoneyFloat(2.0),
				LeadPrice:  billing.MoneyFloat(3.0),
			},
			actionType: adtype.ActionView,
			expected:   billing.MoneyFloat(1.5),
		},
		{
			name: "ActionClick returns ClickPrice",
			priceScope: PriceScope{
				ViewPrice:  billing.MoneyFloat(1.5),
				ClickPrice: billing.MoneyFloat(2.0),
				LeadPrice:  billing.MoneyFloat(3.0),
			},
			actionType: adtype.ActionClick,
			expected:   billing.MoneyFloat(2.0),
		},
		{
			name: "ActionLead returns LeadPrice",
			priceScope: PriceScope{
				ViewPrice:  billing.MoneyFloat(1.5),
				ClickPrice: billing.MoneyFloat(2.0),
				LeadPrice:  billing.MoneyFloat(3.0),
			},
			actionType: adtype.ActionLead,
			expected:   billing.MoneyFloat(3.0),
		},
		{
			name: "Unknown action returns 0",
			priceScope: PriceScope{
				ViewPrice:  billing.MoneyFloat(1.5),
				ClickPrice: billing.MoneyFloat(2.0),
				LeadPrice:  billing.MoneyFloat(3.0),
			},
			actionType: adtype.Action(99), // unknown action
			expected:   0,
		},
		{
			name: "Zero prices return 0",
			priceScope: PriceScope{
				ViewPrice:  0,
				ClickPrice: 0,
				LeadPrice:  0,
			},
			actionType: adtype.ActionView,
			expected:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.priceScope.PricePerAction(tt.actionType)
			if result != tt.expected {
				t.Errorf("PricePerAction() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPriceScope_SetBidViewPrice(t *testing.T) {
	tests := []struct {
		name        string
		priceScope  PriceScope
		price       billing.Money
		maxIfZero   bool
		expectedOk  bool
		expectedBid billing.Money
	}{
		{
			name: "Set valid bid price within max",
			priceScope: PriceScope{
				MaxBidImpPrice: billing.MoneyFloat(10.0),
				BidImpPrice:    billing.MoneyFloat(5.0),
			},
			price:       billing.MoneyFloat(3.0),
			maxIfZero:   false,
			expectedOk:  true,
			expectedBid: billing.MoneyFloat(3.0),
		},
		{
			name: "Set bid price higher than current bid",
			priceScope: PriceScope{
				MaxBidImpPrice: billing.MoneyFloat(10.0),
				BidImpPrice:    billing.MoneyFloat(5.0),
			},
			price:       billing.MoneyFloat(7.0),
			maxIfZero:   false,
			expectedOk:  false,
			expectedBid: billing.MoneyFloat(5.0), // unchanged
		},
		{
			name: "Set zero price with maxIfZero=true",
			priceScope: PriceScope{
				MaxBidImpPrice: billing.MoneyFloat(10.0),
				BidImpPrice:    billing.MoneyFloat(5.0),
			},
			price:       0,
			maxIfZero:   true,
			expectedOk:  true,
			expectedBid: billing.MoneyFloat(10.0),
		},
		{
			name: "Set zero price with maxIfZero=false",
			priceScope: PriceScope{
				MaxBidImpPrice: billing.MoneyFloat(10.0),
				BidImpPrice:    billing.MoneyFloat(5.0),
			},
			price:       0,
			maxIfZero:   false,
			expectedOk:  true,
			expectedBid: 0,
		},
		{
			name: "Set negative price with maxIfZero=true",
			priceScope: PriceScope{
				MaxBidImpPrice: billing.MoneyFloat(10.0),
				BidImpPrice:    billing.MoneyFloat(5.0),
			},
			price:       billing.MoneyFloat(-1.0),
			maxIfZero:   true,
			expectedOk:  true,
			expectedBid: billing.MoneyFloat(10.0),
		},
		{
			name: "Set negative price with maxIfZero=false",
			priceScope: PriceScope{
				MaxBidImpPrice: billing.MoneyFloat(10.0),
				BidImpPrice:    billing.MoneyFloat(5.0),
			},
			price:       billing.MoneyFloat(-1.0),
			maxIfZero:   false,
			expectedOk:  true,
			expectedBid: 0, // max(price, 0)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := tt.priceScope
			result := ps.SetBidImpressionPrice(tt.price, tt.maxIfZero)

			if result != tt.expectedOk {
				t.Errorf("SetBidImpressionPrice() returned %v, want %v", result, tt.expectedOk)
			}

			if ps.BidImpPrice != tt.expectedBid {
				t.Errorf("BidImpPrice = %v, want %v", ps.BidImpPrice, tt.expectedBid)
			}
		})
	}
}

func TestPriceScope_SetViewPrice(t *testing.T) {
	tests := []struct {
		name          string
		priceScope    PriceScope
		price         billing.Money
		maxIfZero     bool
		expectedOk    bool
		expectedPrice billing.Money
	}{
		{
			name: "Set valid view price within max",
			priceScope: PriceScope{
				MaxBidImpPrice: billing.MoneyFloat(10.0),
				ImpPrice:       billing.MoneyFloat(3.0),
			},
			price:         billing.MoneyFloat(5.0),
			maxIfZero:     false,
			expectedOk:    true,
			expectedPrice: billing.MoneyFloat(5.0),
		},
		{
			name: "Set view price higher than max",
			priceScope: PriceScope{
				MaxBidImpPrice: billing.MoneyFloat(10.0),
				ImpPrice:       billing.MoneyFloat(3.0),
			},
			price:         billing.MoneyFloat(15.0),
			maxIfZero:     false,
			expectedOk:    true,
			expectedPrice: billing.MoneyFloat(10.0), // capped to max
		},
		{
			name: "Set zero price with maxIfZero=true",
			priceScope: PriceScope{
				MaxBidImpPrice: billing.MoneyFloat(10.0),
				ImpPrice:       billing.MoneyFloat(3.0),
			},
			price:         0,
			maxIfZero:     true,
			expectedOk:    true,
			expectedPrice: billing.MoneyFloat(10.0),
		},
		{
			name: "Set zero price with maxIfZero=false",
			priceScope: PriceScope{
				MaxBidImpPrice: billing.MoneyFloat(10.0),
				ImpPrice:       billing.MoneyFloat(3.0),
			},
			price:         0,
			maxIfZero:     false,
			expectedOk:    true,
			expectedPrice: 0,
		},
		{
			name: "Set negative price with maxIfZero=true",
			priceScope: PriceScope{
				MaxBidImpPrice: billing.MoneyFloat(10.0),
				ImpPrice:       billing.MoneyFloat(3.0),
			},
			price:         billing.MoneyFloat(-2.0),
			maxIfZero:     true,
			expectedOk:    true,
			expectedPrice: billing.MoneyFloat(10.0),
		},
		{
			name: "Set negative price with maxIfZero=false",
			priceScope: PriceScope{
				MaxBidImpPrice: billing.MoneyFloat(10.0),
				ImpPrice:       billing.MoneyFloat(3.0),
			},
			price:         billing.MoneyFloat(-2.0),
			maxIfZero:     false,
			expectedOk:    true,
			expectedPrice: 0, // max(price, 0)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := tt.priceScope
			result := ps.SetImpressionPrice(tt.price, tt.maxIfZero)

			if result != tt.expectedOk {
				t.Errorf("SetImpressionPrice() returned %v, want %v", result, tt.expectedOk)
			}

			if ps.ImpPrice != tt.expectedPrice {
				t.Errorf("ImpPrice = %v, want %v", ps.ImpPrice, tt.expectedPrice)
			}
		})
	}
}

func TestPriceScope_PrepareBidViewPrice(t *testing.T) {
	tests := []struct {
		name       string
		priceScope PriceScope
		price      billing.Money
		expected   billing.Money
	}{
		{
			name: "Price within max limit",
			priceScope: PriceScope{
				MaxBidImpPrice: billing.MoneyFloat(10.0),
			},
			price:    billing.MoneyFloat(5.0),
			expected: billing.MoneyFloat(5.0),
		},
		{
			name: "Price exceeds max limit",
			priceScope: PriceScope{
				MaxBidImpPrice: billing.MoneyFloat(10.0),
			},
			price:    billing.MoneyFloat(15.0),
			expected: billing.MoneyFloat(10.0),
		},
		{
			name: "Price equals max limit",
			priceScope: PriceScope{
				MaxBidImpPrice: billing.MoneyFloat(10.0),
			},
			price:    billing.MoneyFloat(10.0),
			expected: billing.MoneyFloat(10.0),
		},
		{
			name: "Zero price",
			priceScope: PriceScope{
				MaxBidImpPrice: billing.MoneyFloat(10.0),
			},
			price:    0,
			expected: 0,
		},
		{
			name: "Negative price",
			priceScope: PriceScope{
				MaxBidImpPrice: billing.MoneyFloat(10.0),
			},
			price:    billing.MoneyFloat(-5.0),
			expected: billing.MoneyFloat(-5.0),
		},
		{
			name: "Zero max bid price",
			priceScope: PriceScope{
				MaxBidImpPrice: 0,
			},
			price:    billing.MoneyFloat(5.0),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.priceScope.PrepareBidImpressionPrice(tt.price)
			if result != tt.expected {
				t.Errorf("PrepareBidImpressionPrice() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPriceScope_CompleteWorkflow(t *testing.T) {
	// Test a complete workflow scenario
	ps := PriceScope{
		TestMode:       true,
		MaxBidImpPrice: billing.MoneyFloat(10.0),
		BidImpPrice:    billing.MoneyFloat(8.0), // Set initial bid price
		ImpPrice:       billing.MoneyFloat(2.0),
		ClickPrice:     billing.MoneyFloat(4.0),
		LeadPrice:      billing.MoneyFloat(8.0),
		ECPM:           billing.MoneyFloat(5.0),
	}

	// Test all action types
	if got := ps.PricePerAction(adtype.ActionImpression); got != billing.MoneyFloat(2.0) {
		t.Errorf("ImpressionPrice = %v, want %v", got, billing.MoneyFloat(2.0))
	}

	if got := ps.PricePerAction(adtype.ActionClick); got != billing.MoneyFloat(4.0) {
		t.Errorf("ClickPrice = %v, want %v", got, billing.MoneyFloat(4.0))
	}

	if got := ps.PricePerAction(adtype.ActionLead); got != billing.MoneyFloat(8.0) {
		t.Errorf("LeadPrice = %v, want %v", got, billing.MoneyFloat(8.0))
	}

	// Test setting bid impression price (should succeed because 6.0 < 8.0 current bid)
	if !ps.SetBidImpressionPrice(billing.MoneyFloat(6.0), false) {
		t.Error("SetBidImpressionPrice should succeed for valid price")
	}
	if ps.BidImpPrice != billing.MoneyFloat(6.0) {
		t.Errorf("BidImpPrice = %v, want %v", ps.BidImpPrice, billing.MoneyFloat(6.0))
	}

	// Test setting impression price
	if !ps.SetImpressionPrice(billing.MoneyFloat(3.0), false) {
		t.Error("SetImpressionPrice should succeed for valid price")
	}
	if ps.ImpPrice != billing.MoneyFloat(3.0) {
		t.Errorf("ImpPrice = %v, want %v", ps.ImpPrice, billing.MoneyFloat(3.0))
	}

	// Test prepare bid impression price
	prepared := ps.PrepareBidImpressionPrice(billing.MoneyFloat(12.0))
	if prepared != billing.MoneyFloat(10.0) {
		t.Errorf("PrepareBidImpressionPrice = %v, want %v", prepared, billing.MoneyFloat(10.0))
	}
}

func TestPriceScope_EdgeCases(t *testing.T) {
	t.Run("All zero values", func(t *testing.T) {
		ps := PriceScope{}

		// All prices should be 0
		if got := ps.PricePerAction(adtype.ActionView); got != 0 {
			t.Errorf("ViewPrice = %v, want 0", got)
		}

		// Setting prices with zero max should work
		if !ps.SetBidImpressionPrice(0, true) {
			t.Error("SetBidImpressionPrice should succeed with zero max")
		}

		if ps.BidImpPrice != 0 {
			t.Errorf("BidImpPrice = %v, want 0", ps.BidImpPrice)
		}
	})

	t.Run("Large values", func(t *testing.T) {
		largeValue := billing.MoneyFloat(999999.99)
		ps := PriceScope{
			MaxBidImpPrice: largeValue,
			ImpPrice:       largeValue,
			ClickPrice:     largeValue,
			LeadPrice:      largeValue,
		}

		if got := ps.PricePerAction(adtype.ActionImpression); got != largeValue {
			t.Errorf("ImpressionPrice = %v, want %v", got, largeValue)
		}

		prepared := ps.PrepareBidImpressionPrice(largeValue * 2)
		if prepared != largeValue {
			t.Errorf("PrepareBidImpressionPrice = %v, want %v", prepared, largeValue)
		}
	})
}
