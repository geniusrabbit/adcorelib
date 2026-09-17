package geo

// RegionCodeByPartial accepts a full ISO 3166-2 code ("SK-BA") or a
// subdivision suffix ("BA") when the parent country is known.
// Empty and unknown values return UndefinedRegionCode.
func RegionCodeByPartial(code string, country Code2) RegionCode {
	if code == "" {
		return UndefinedRegionCode
	}
	if rc := RegionCodeByString(code); rc != UndefinedRegionCode {
		return rc
	}
	if country.ISO2() != UndefinedCountryCodeISO2 {
		return RegionCodeByString(country.ISO2() + "-" + code)
	}
	return UndefinedRegionCode
}
