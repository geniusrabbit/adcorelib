package admodels

import "testing"

func TestClosestThumbBy(t *testing.T) {
	asset := &AdFileAsset{
		Thumbs: []AdFileAssetThumb{
			{URL: "small", Width: 100, Height: 80},
			{URL: "slot", Width: 300, Height: 250},
			{URL: "wide", Width: 728, Height: 90},
			{URL: "large", Width: 600, Height: 500},
		},
	}

	t.Run("prefers fitting closest to target", func(t *testing.T) {
		got := asset.ClosestThumbBy(300, 250, 0, 0)
		if got == nil || got.URL != "slot" {
			t.Fatalf("got %+v, want slot", got)
		}
	})

	t.Run("falls back to closest when none fit", func(t *testing.T) {
		got := asset.ClosestThumbBy(80, 60, 0, 0)
		if got == nil || got.URL != "small" {
			t.Fatalf("got %+v, want small", got)
		}
	})

	t.Run("empty thumbs", func(t *testing.T) {
		if (&AdFileAsset{}).ClosestThumbBy(300, 250, 0, 0) != nil {
			t.Fatal("expected nil")
		}
	})
}

func TestThumbByPicksSmallestFitting(t *testing.T) {
	asset := &AdFileAsset{
		Thumbs: []AdFileAssetThumb{
			{URL: "small", Width: 100, Height: 80},
			{URL: "slot", Width: 300, Height: 250},
		},
	}
	got := asset.ThumbBy(300, 250, 0, 0)
	if got == nil || got.URL != "small" {
		t.Fatalf("got %+v, want small", got)
	}
}
