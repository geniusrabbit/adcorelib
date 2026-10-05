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

	t.Run("minimums pick the closest fitting thumb", func(t *testing.T) {
		mins := &AdFileAsset{
			Thumbs: []AdFileAssetThumb{
				{URL: "large", Width: 1200, Height: 800},
				{URL: "slot", Width: 320, Height: 250},
			},
		}
		got := mins.ClosestThumbBy(0, 0, 300, 200)
		if got == nil || got.URL != "slot" {
			t.Fatalf("got %+v, want slot", got)
		}
	})

	t.Run("fit beats a closer thumb outside the box", func(t *testing.T) {
		sized := &AdFileAsset{
			Thumbs: []AdFileAssetThumb{
				{URL: "close", Width: 310, Height: 240},
				{URL: "fit", Width: 200, Height: 200},
			},
		}
		got := sized.ClosestThumbBy(300, 250, 150, 150)
		if got == nil || got.URL != "fit" {
			t.Fatalf("got %+v, want fit", got)
		}
	})

	t.Run("equal distance keeps the smaller thumb", func(t *testing.T) {
		tied := &AdFileAsset{
			Thumbs: []AdFileAssetThumb{
				{URL: "wide", Width: 400, Height: 100},
				{URL: "tall", Width: 100, Height: 400},
			},
		}
		got := tied.ClosestThumbBy(250, 250, 0, 0)
		if got == nil || got.URL != "tall" {
			t.Fatalf("got %+v, want tall", got)
		}
	})

	t.Run("url fallback", func(t *testing.T) {
		got := (&AdFileAsset{URL: "main.jpg", Width: 1626, Height: 967}).ClosestThumbBy(300, 250, 0, 0)
		if got == nil || got.URL != "main.jpg" || got.Width != 1626 || got.Height != 967 {
			t.Fatalf("got %+v, want main.jpg 1626x967", got)
		}
	})

	t.Run("nil asset", func(t *testing.T) {
		if (*AdFileAsset)(nil).ClosestThumbBy(300, 250, 0, 0) != nil {
			t.Fatal("expected nil")
		}
	})

	t.Run("arbitrary order", func(t *testing.T) {
		thumbs := []AdFileAssetThumb{
			{URL: "huge", Width: 1600, Height: 900},
			{URL: "tiny", Width: 100, Height: 80},
			{URL: "wide", Width: 1200, Height: 200},
			{URL: "mid", Width: 640, Height: 360},
		}
		orders := [][]int{
			{0, 1, 2, 3},
			{3, 2, 1, 0},
			{1, 3, 0, 2},
			{2, 0, 3, 1},
		}
		for _, order := range orders {
			shuffled := make([]AdFileAssetThumb, len(order))
			for i, idx := range order {
				shuffled[i] = thumbs[idx]
			}
			ordered := &AdFileAsset{Thumbs: shuffled}
			if got := ordered.ClosestThumbBy(640, 360, 0, 0); got == nil || got.URL != "mid" {
				t.Fatalf("exact %v: got %+v, want mid", order, got)
			}
			if got := ordered.ClosestThumbBy(0, 0, 600, 300); got == nil || got.URL != "mid" {
				t.Fatalf("minimums %v: got %+v, want mid", order, got)
			}
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
