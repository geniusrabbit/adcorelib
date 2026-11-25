package simple

import (
	"context"
	"testing"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/udetect"
)

func ver(s string) types.Version {
	return types.MustParseVersion(s)
}

func TestSimpleClientDetectIDs(t *testing.T) {
	client := New(
		[]*Item{
			{ID: 8, Name: "Chrome", Versions: []*VersionItem{
				{ID: 580, Min: ver("58.0"), Max: ver("59.0")},
				{ID: 590, Min: ver("59.0"), Max: ver("60.0")},
			}},
			{ID: 5, Name: "Firefox"},
			{ID: 15, Name: "Safari"},
			{ID: 41, Name: "Edge"},
		},
		[]*Item{
			{ID: 1, Name: "Windows", Versions: []*VersionItem{
				{ID: 61, Min: ver("6.1"), Max: ver("6.2")},
			}},
			{ID: 10, Name: "macOS"},
			{ID: 61, Name: "iOS"},
			{ID: 89, Name: "Android"},
		},
	)

	tests := []struct {
		name      string
		ua        string
		browserID uint64
		osID      uint
	}{
		{
			name:      "safari macos",
			ua:        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8",
			browserID: 15,
			osID:      10,
		},
		{
			name:      "chrome windows",
			ua:        "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
			browserID: 8,
			osID:      1,
		},
		{
			name:      "safari ios",
			ua:        "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1",
			browserID: 15,
			osID:      61,
		},
		{
			name:      "chrome android",
			ua:        "Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.125 Mobile Safari/537.36",
			browserID: 8,
			osID:      89,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.Detect(context.Background(), &udetect.Request{UA: tt.ua})
			if err != nil {
				t.Fatalf("Detect: %v", err)
			}
			if resp.Device == nil || resp.Device.OS == nil || resp.Device.Browser == nil {
				t.Fatal("expected device, os, and browser")
			}
			if resp.Device.OS.ID != tt.osID {
				t.Fatalf("os id: got %d want %d (parsed %q)", resp.Device.OS.ID, tt.osID, resp.Device.OS.Name)
			}
			if resp.Device.Browser.ID != tt.browserID {
				t.Fatalf("browser id: got %d want %d (parsed %q)", resp.Device.Browser.ID, tt.browserID, resp.Device.Browser.Name)
			}
		})
	}
}

func TestSimpleClientVersionRanges(t *testing.T) {
	client := New(
		[]*Item{{
			ID: 8, Name: "Chrome",
			Versions: []*VersionItem{
				{ID: 580, Min: ver("58.0"), Max: ver("59.0")},
				{ID: 590, Min: ver("59.0"), Max: ver("60.0")},
			},
		}},
		[]*Item{{
			ID: 1, Name: "Windows",
			Versions: []*VersionItem{
				{ID: 100, Min: ver("6.1"), Max: ver("6.2")},
			},
		}},
	)

	chrome := "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36"
	resp, err := client.Detect(context.Background(), &udetect.Request{UA: chrome})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Device.Browser.ID != 8 || resp.Device.Browser.VersionID != 590 {
		t.Fatalf("chrome: root %d version %d", resp.Device.Browser.ID, resp.Device.Browser.VersionID)
	}
	if resp.Device.OS.ID != 1 || resp.Device.OS.VersionID != 100 {
		t.Fatalf("windows: root %d version %d", resp.Device.OS.ID, resp.Device.OS.VersionID)
	}

	boundary := "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.0.0 Safari/537.36"
	resp, err = client.Detect(context.Background(), &udetect.Request{UA: boundary})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Device.Browser.ID != 8 || resp.Device.Browser.VersionID != 0 {
		t.Fatalf("chrome 60 boundary: root %d version %d", resp.Device.Browser.ID, resp.Device.Browser.VersionID)
	}

	unknown := "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.0.0 Safari/537.36"
	resp, err = client.Detect(context.Background(), &udetect.Request{UA: unknown})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Device.Browser.ID != 8 || resp.Device.Browser.VersionID != 0 {
		t.Fatalf("chrome 57: root %d version %d", resp.Device.Browser.ID, resp.Device.Browser.VersionID)
	}
}
