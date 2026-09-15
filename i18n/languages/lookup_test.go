package languages

import "testing"

func Test_GetLanguageByID(t *testing.T) {
	if got := GetLanguageByID(40); got.Name != "English" {
		t.Fatalf("id 40: %+v", got)
	}
	if got := GetLanguageByID(0); !got.Code.IsUndefined() {
		t.Fatalf("id 0 should be undefined: %+v", got)
	}
	if got := GetLanguageByID(uint(len(languages))); got != &languages[0] {
		t.Fatal("out-of-range id should be Languages[0]")
	}
}

func Test_CatalogIndexEqualsID(t *testing.T) {
	if Languages[0].Code != UndefinedCode {
		t.Fatalf("Languages[0] code = %q", Languages[0].Code)
	}
	for i, lg := range Languages {
		if lg.ID != uint(i) {
			t.Fatalf("index %d has ID %d", i, lg.ID)
		}
		if GetLanguageByID(lg.ID) != &languages[i] {
			t.Fatalf("GetLanguageByID(%d) is not catalog slot", lg.ID)
		}
		if i == 0 {
			continue
		}
		if lg.Code[0] < 'a' || lg.Code[0] > 'z' || lg.Code[1] < 'a' || lg.Code[1] > 'z' {
			t.Fatalf("code %q is not lowercase", lg.Code)
		}
		if GetLanguageByCode(lg.Code) != &languages[i] {
			t.Fatalf("code %q does not round-trip", lg.Code)
		}
	}
}

func Test_LangCodes2IDs(t *testing.T) {
	if LangCodes2IDs(nil) != nil {
		t.Fatal("nil codes should return nil")
	}
	ids := LangCodes2IDs([]string{"en", "xx", "DE"})
	if len(ids) != 3 {
		t.Fatalf("len=%d", len(ids))
	}
	if ids[0] != 0 {
		t.Fatalf("unknown xx should sort first as 0, got %v", ids)
	}
	if ids[1] != uint64(GetLanguageIdByCodeString("EN")) {
		t.Fatalf("expected EN id, got %v", ids)
	}
	if ids[2] != uint64(GetLanguageIdByCodeString("DE")) {
		t.Fatalf("expected DE id, got %v", ids)
	}
}
