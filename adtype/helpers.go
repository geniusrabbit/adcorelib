//
// @project GeniusRabbit rotator 2016
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016
//

package adtype

import "github.com/geniusrabbit/gogeo"

func indexOfStringArray(v string, arr []string) int {
	if len(arr) < 1 {
		return -1
	}
	for i, s := range arr {
		if s == v {
			return i
		}
	}
	return -1
}

func defStr(v, def string) string {
	if len(v) < 1 {
		return def
	}
	return v
}

func defCountryCode(code1, code2 string) (cc string) {
	if code1 == "**" || code1 == "" {
		cc = code2
	} else {
		cc = code1
	}
	if len(cc) == 3 {
		return gogeo.CountryByCode3(cc).Code2
	}
	return cc
}

func defInt(v, def int) int {
	if v == 0 {
		return def
	}
	return v
}

func defFloat(v, def float64) float64 {
	if v == 0 {
		return def
	}
	return v
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
