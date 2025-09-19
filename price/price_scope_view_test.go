package price

import (
	"testing"

	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

func TestPriceScopeView_PricePerAction(t *testing.T) {
	tests := []struct {
		name       string
		priceScope PriceScopeView
		actionType adtype.Action
		expected   billing.Money
	}{
		{
			name: "ActionView returns ViewPrice",
			priceScope: PriceScopeView{
				ViewPrice: billing.MoneyFloat(1.5),
			},
			actionType: adtype.ActionView,
			expected:   billing.MoneyFloat(1.5),
		},
		{
			name: "ActionClick returns 0",
			priceScope: PriceScopeView{
				ViewPrice: billing.MoneyFloat(1.5),
			},
			actionType: adtype.ActionClick,
			expected:   0,
		},
		{
			name: "ActionLead returns 0",
			priceScope: PriceScopeView{
				ViewPrice: billing.MoneyFloat(1.5),
			},
			actionType: adtype.ActionLead,
			expected:   0,
		},
		{
			name: "Unknown action returns 0",
			priceScope: PriceScopeView{
				ViewPrice: billing.MoneyFloat(1.5),
			},
			actionType: adtype.Action(99), // unknown action
			expected:   0,
		},
		{
			name: "Zero ViewPrice returns 0",
			priceScope: PriceScopeView{
				ViewPrice: 0,
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

func TestPriceScopeView_SetBidViewPrice(t *testing.T) {
	tests := []struct {
		name        string
		priceScope  PriceScopeView
		price       billing.Money
		maxIfZero   bool
		expectedOk  bool
		expectedBid billing.Money
	}{
		{
			name: "Set valid bid price within max",
			priceScope: PriceScopeView{
				MaxBidViewPrice: billing.MoneyFloat(10.0),
				BidViewPrice:    billing.MoneyFloat(5.0),
			},
			price:       billing.MoneyFloat(3.0),
			maxIfZero:   false,
			expectedOk:  true,
			expectedBid: billing.MoneyFloat(3.0),
		},
		{
			name: "Set bid price higher than current bid",
			priceScope: PriceScopeView{
				MaxBidViewPrice: billing.MoneyFloat(10.0),
				BidViewPrice:    billing.MoneyFloat(5.0),
			},
			price:       billing.MoneyFloat(7.0),
			maxIfZero:   false,
			expectedOk:  true, // PriceScopeView checks against MaxBidViewPrice, not current BidViewPrice
			expectedBid: billing.MoneyFloat(7.0),
		},
		{
			name: "Set zero price with maxIfZero=true",
			priceScope: PriceScopeView{
				MaxBidViewPrice: billing.MoneyFloat(10.0),
				BidViewPrice:    billing.MoneyFloat(5.0),
			},
			price:       0,
			maxIfZero:   true,
			expectedOk:  true,
			expectedBid: billing.MoneyFloat(10.0),
		},
		{
			name: "Set zero price with maxIfZero=false",
			priceScope: PriceScopeView{
				MaxBidViewPrice: billing.MoneyFloat(10.0),
				BidViewPrice:    billing.MoneyFloat(5.0),
			},
			price:       0,
			maxIfZero:   false,
			expectedOk:  true,
			expectedBid: 0,
		},
		{
			name: "Set negative price with maxIfZero=true",
			priceScope: PriceScopeView{
				MaxBidViewPrice: billing.MoneyFloat(10.0),
				BidViewPrice:    billing.MoneyFloat(5.0),
			},
			price:       billing.MoneyFloat(-1.0),
			maxIfZero:   true,
			expectedOk:  true,
			expectedBid: billing.MoneyFloat(10.0),
		},
		{
			name: "Set negative price with maxIfZero=false",
			priceScope: PriceScopeView{
				MaxBidViewPrice: billing.MoneyFloat(10.0),
				BidViewPrice:    billing.MoneyFloat(5.0),
			},
			price:       billing.MoneyFloat(-1.0),
			maxIfZero:   false,
			expectedOk:  true,
			expectedBid: 0, // max(price, 0)
		},
		{
			name: "Set price higher than max bid",
			priceScope: PriceScopeView{
				MaxBidViewPrice: billing.MoneyFloat(10.0),
				BidViewPrice:    billing.MoneyFloat(5.0),
			},
			price:       billing.MoneyFloat(15.0),
			maxIfZero:   false,
			expectedOk:  false,
			expectedBid: billing.MoneyFloat(5.0), // unchanged
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := tt.priceScope
			result := ps.SetBidViewPrice(tt.price, tt.maxIfZero)

			if result != tt.expectedOk {
				t.Errorf("SetBidViewPrice() returned %v, want %v", result, tt.expectedOk)
			}

			if ps.BidViewPrice != tt.expectedBid {
				t.Errorf("BidViewPrice = %v, want %v", ps.BidViewPrice, tt.expectedBid)
			}
		})
	}
}

func TestPriceScopeView_SetViewPrice(t *testing.T) {
	tests := []struct {
		name          string
		priceScope    PriceScopeView
		price         billing.Money
		maxIfZero     bool
		expectedOk    bool
		expectedPrice billing.Money
	}{
		{
			name: "Set valid view price within max",
			priceScope: PriceScopeView{
				MaxBidViewPrice: billing.MoneyFloat(10.0),
				ViewPrice:       billing.MoneyFloat(3.0),
			},
			price:         billing.MoneyFloat(5.0),
			maxIfZero:     false,
			expectedOk:    true,
			expectedPrice: billing.MoneyFloat(5.0),
		},
		{
			name: "Set view price higher than max",
			priceScope: PriceScopeView{
				MaxBidViewPrice: billing.MoneyFloat(10.0),
				ViewPrice:       billing.MoneyFloat(3.0),
			},
			price:         billing.MoneyFloat(15.0),
			maxIfZero:     false,
			expectedOk:    true,
			expectedPrice: billing.MoneyFloat(10.0), // capped to max
		},
		{
			name: "Set zero price with maxIfZero=true",
			priceScope: PriceScopeView{
				MaxBidViewPrice: billing.MoneyFloat(10.0),
				ViewPrice:       billing.MoneyFloat(3.0),
			},
			price:         0,
			maxIfZero:     true,
			expectedOk:    true,
			expectedPrice: billing.MoneyFloat(10.0),
		},
		{
			name: "Set zero price with maxIfZero=false",
			priceScope: PriceScopeView{
				MaxBidViewPrice: billing.MoneyFloat(10.0),
				ViewPrice:       billing.MoneyFloat(3.0),
			},
			price:         0,
			maxIfZero:     false,
			expectedOk:    true,
			expectedPrice: 0,
		},
		{
			name: "Set negative price with maxIfZero=true",
			priceScope: PriceScopeView{
				MaxBidViewPrice: billing.MoneyFloat(10.0),
				ViewPrice:       billing.MoneyFloat(3.0),
			},
			price:         billing.MoneyFloat(-2.0),
			maxIfZero:     true,
			expectedOk:    true,
			expectedPrice: billing.MoneyFloat(10.0),
		},
		{
			name: "Set negative price with maxIfZero=false",
			priceScope: PriceScopeView{
				MaxBidViewPrice: billing.MoneyFloat(10.0),
				ViewPrice:       billing.MoneyFloat(3.0),
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
			result := ps.SetViewPrice(tt.price, tt.maxIfZero)

			if result != tt.expectedOk {
				t.Errorf("SetViewPrice() returned %v, want %v", result, tt.expectedOk)
			}

			if ps.ViewPrice != tt.expectedPrice {
				t.Errorf("ViewPrice = %v, want %v", ps.ViewPrice, tt.expectedPrice)
			}
		})
	}
}

func TestPriceScopeView_PrepareBidViewPrice(t *testing.T) {
	tests := []struct {
		name       string
		priceScope PriceScopeView
		price      billing.Money
		expected   billing.Money
	}{
		{
			name: "Price within max limit",
			priceScope: PriceScopeView{
				MaxBidViewPrice: billing.MoneyFloat(10.0),
			},
			price:    billing.MoneyFloat(5.0),
			expected: billing.MoneyFloat(5.0),
		},
		{
			name: "Price exceeds max limit",
			priceScope: PriceScopeView{
				MaxBidViewPrice: billing.MoneyFloat(10.0),
			},
			price:    billing.MoneyFloat(15.0),
			expected: billing.MoneyFloat(10.0),
		},
		{
			name: "Price equals max limit",
			priceScope: PriceScopeView{
				MaxBidViewPrice: billing.MoneyFloat(10.0),
			},
			price:    billing.MoneyFloat(10.0),
			expected: billing.MoneyFloat(10.0),
		},
		{
			name: "Zero price",
			priceScope: PriceScopeView{
				MaxBidViewPrice: billing.MoneyFloat(10.0),
			},
			price:    0,
			expected: 0,
		},
		{
			name: "Negative price",
			priceScope: PriceScopeView{
				MaxBidViewPrice: billing.MoneyFloat(10.0),
			},
			price:    billing.MoneyFloat(-5.0),
			expected: billing.MoneyFloat(-5.0),
		},
		{
			name: "Zero max bid price",
			priceScope: PriceScopeView{
				MaxBidViewPrice: 0,
			},
			price:    billing.MoneyFloat(5.0),
			expected: billing.MoneyFloat(5.0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.priceScope.PrepareBidViewPrice(tt.price)
			if result != tt.expected {
				t.Errorf("PrepareBidViewPrice() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPriceScopeView_CompleteWorkflow(t *testing.T) {
	// Test a complete workflow scenario
	ps := PriceScopeView{
		MaxBidViewPrice: billing.MoneyFloat(10.0),
		ViewPrice:       billing.MoneyFloat(2.0),
		ECPM:            billing.MoneyFloat(5.0),
	}

	// Test action types - only ActionView should return non-zero
	if got := ps.PricePerAction(adtype.ActionView); got != billing.MoneyFloat(2.0) {
		t.Errorf("ViewPrice = %v, want %v", got, billing.MoneyFloat(2.0))
	}

	if got := ps.PricePerAction(adtype.ActionClick); got != 0 {
		t.Errorf("ClickPrice = %v, want 0", got)
	}

	if got := ps.PricePerAction(adtype.ActionLead); got != 0 {
		t.Errorf("LeadPrice = %v, want 0", got)
	}

	// Test setting bid view price
	if !ps.SetBidViewPrice(billing.MoneyFloat(6.0), false) {
		t.Error("SetBidViewPrice should succeed for valid price")
	}
	if ps.BidViewPrice != billing.MoneyFloat(6.0) {
		t.Errorf("BidViewPrice = %v, want %v", ps.BidViewPrice, billing.MoneyFloat(6.0))
	}

	// Test setting view price
	if !ps.SetViewPrice(billing.MoneyFloat(3.0), false) {
		t.Error("SetViewPrice should succeed for valid price")
	}
	if ps.ViewPrice != billing.MoneyFloat(3.0) {
		t.Errorf("ViewPrice = %v, want %v", ps.ViewPrice, billing.MoneyFloat(3.0))
	}

	// Test prepare bid view price
	prepared := ps.PrepareBidViewPrice(billing.MoneyFloat(12.0))
	if prepared != billing.MoneyFloat(10.0) {
		t.Errorf("PrepareBidViewPrice = %v, want %v", prepared, billing.MoneyFloat(10.0))
	}
}

func TestPriceScopeView_EdgeCases(t *testing.T) {
	t.Run("All zero values", func(t *testing.T) {
		ps := PriceScopeView{}

		// All prices should be 0
		if got := ps.PricePerAction(adtype.ActionView); got != 0 {
			t.Errorf("ViewPrice = %v, want 0", got)
		}

		// Setting prices with zero max should work
		if !ps.SetBidViewPrice(0, true) {
			t.Error("SetBidViewPrice should succeed with zero max")
		}

		if ps.BidViewPrice != 0 {
			t.Errorf("BidViewPrice = %v, want 0", ps.BidViewPrice)
		}
	})

	t.Run("Large values", func(t *testing.T) {
		largeValue := billing.MoneyFloat(999999.99)
		ps := PriceScopeView{
			MaxBidViewPrice: largeValue,
			ViewPrice:       largeValue,
		}

		if got := ps.PricePerAction(adtype.ActionView); got != largeValue {
			t.Errorf("ViewPrice = %v, want %v", got, largeValue)
		}

		prepared := ps.PrepareBidViewPrice(largeValue * 2)
		if prepared != largeValue {
			t.Errorf("PrepareBidViewPrice = %v, want %v", prepared, largeValue)
		}
	})
}

func TestPriceScopeView_CompareWithPriceScope(t *testing.T) {
	// Test to ensure PriceScopeView behaves differently than PriceScope for non-view actions
	psView := PriceScopeView{
		MaxBidViewPrice: billing.MoneyFloat(10.0),
		ViewPrice:       billing.MoneyFloat(2.0),
	}

	ps := PriceScope{
		MaxBidImpPrice: billing.MoneyFloat(10.0),
		ViewPrice:      billing.MoneyFloat(2.0),
		ClickPrice:     billing.MoneyFloat(4.0),
		LeadPrice:      billing.MoneyFloat(8.0),
	}

	// Both should return the same for ActionView
	if psView.PricePerAction(adtype.ActionView) != ps.PricePerAction(adtype.ActionView) {
		t.Error("PriceScopeView and PriceScope should return same price for ActionView")
	}

	// PriceScopeView should return 0 for Click and Lead, while PriceScope returns actual prices
	if psView.PricePerAction(adtype.ActionClick) != 0 {
		t.Error("PriceScopeView should return 0 for ActionClick")
	}

	if psView.PricePerAction(adtype.ActionLead) != 0 {
		t.Error("PriceScopeView should return 0 for ActionLead")
	}

	if ps.PricePerAction(adtype.ActionClick) == 0 {
		t.Error("PriceScope should return non-zero for ActionClick")
	}

	if ps.PricePerAction(adtype.ActionLead) == 0 {
		t.Error("PriceScope should return non-zero for ActionLead")
	}
}

func TestPriceScopeView_StructFields(t *testing.T) {
	// Test struct field access
	ps := PriceScopeView{
		MaxBidViewPrice: billing.MoneyFloat(10.0),
		BidViewPrice:    billing.MoneyFloat(6.0),
		ViewPrice:       billing.MoneyFloat(2.0),
		ECPM:            billing.MoneyFloat(5.0),
	}

	if ps.MaxBidViewPrice != billing.MoneyFloat(10.0) {
		t.Errorf("MaxBidViewPrice = %v, want %v", ps.MaxBidViewPrice, billing.MoneyFloat(10.0))
	}

	if ps.BidViewPrice != billing.MoneyFloat(6.0) {
		t.Errorf("BidViewPrice = %v, want %v", ps.BidViewPrice, billing.MoneyFloat(6.0))
	}

	if ps.ViewPrice != billing.MoneyFloat(2.0) {
		t.Errorf("ViewPrice = %v, want %v", ps.ViewPrice, billing.MoneyFloat(2.0))
	}

	if ps.ECPM != billing.MoneyFloat(5.0) {
		t.Errorf("ECPM = %v, want %v", ps.ECPM, billing.MoneyFloat(5.0))
	}
}
