package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

var ErrInvalidParseVersion = fmt.Errorf("invalid parse version")

// Version is a major.minor.patch triple. Each component fits in a uint16.
// Components after the third are dropped.
type Version struct {
	Major uint16 `json:"major"`
	Minor uint16 `json:"minor"`
	Patch uint16 `json:"patch"`
}

// ParseVersion parses str into a Version.
func ParseVersion(str string) (Version, error) {
	var v Version
	if err := v.SetFromStr(str); err != nil {
		return v, err
	}
	return v, nil
}

// MustParseVersion parses str into a Version. It panics on error.
func MustParseVersion(str string) Version {
	v, err := ParseVersion(str)
	if err != nil {
		panic(err)
	}
	return v
}

// IgnoreParseVersion parses str into a Version. It returns the empty version on error.
func IgnoreParseVersion(str string) Version {
	var v Version
	_ = v.SetFromStr(str)
	return v
}

// String returns the version as a string.
func (v *Version) String() string {
	if v.Patch > 0 {
		return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	}
	if v.Minor > 0 {
		return fmt.Sprintf("%d.%d", v.Major, v.Minor)
	}
	if v.Major == 0 {
		return ""
	}
	return fmt.Sprintf("%d", v.Major)
}

// Value implements the driver.Valuer interface, json field interface
func (v Version) Value() (driver.Value, error) {
	return v.String(), nil
}

// Scan implements the driver.Valuer interface, json field interface
func (v *Version) Scan(value any) error {
	switch t := value.(type) {
	case nil:
		*v = Version{}
	case string:
		return v.SetFromStr(t)
	case []byte:
		return v.SetFromStr(string(t))
	case Version:
		*v = t
	case *Version:
		*v = *t
	default:
		return fmt.Errorf("cannot convert %T to Version", t)
	}
	return nil
}

// IsEmpty returns true if the version is empty.
func (v *Version) IsEmpty() bool {
	return v.Major == 0 && v.Minor == 0 && v.Patch == 0
}

// Less returns true if v is less than other.
func (v *Version) Less(other Version) bool {
	if v.Major < other.Major {
		return true
	} else if v.Major > other.Major {
		return false
	} else if v.Minor < other.Minor {
		return true
	} else if v.Minor > other.Minor {
		return false
	}
	return v.Patch < other.Patch
}

// MarshalJSON returns the version as a JSON string.
func (v *Version) MarshalJSON() ([]byte, error) {
	return []byte(`"` + v.String() + `"`), nil
}

// UnmarshalJSON parses the JSON-encoded data and stores the result in the version.
func (v *Version) UnmarshalJSON(data []byte) error {
	s := string(data)
	if len(s) > 1 && s[0] == '"' && s[len(s)-1] == '"' {
		return v.SetFromStr(s[1 : len(s)-1])
	}
	return ErrInvalidParseVersion
}

// SetFromStr parses str into Major, Minor, and Patch.
//
// Leading and trailing space is ignored. An optional "v" or "V" prefix is
// accepted. Numbers are decimal and separated by '.' or '_'. Only the first
// three are stored: "145.0.0.0" is 145.0.0, "10.0" is 10.0, "10_15_7" is
// 10.15.7. A '-' or '+' after that prefix (semver pre-release or build), and
// any other trailer, is ignored: "v1.2.3-rc.1" is 1.2.3. "", "0", and
// "undefined" are the empty version.
//
// A string with no number, an empty component ("10..0"), or a component
// above 65535 returns ErrInvalidParseVersion.
func (v *Version) SetFromStr(str string) error {
	s := strings.TrimSpace(str)
	if s == "" || s == "0" || s == "undefined" {
		*v = Version{}
		return nil
	}
	if s[0] == 'v' || s[0] == 'V' {
		s = s[1:]
	}

	var parts [3]uint16
	n := 0
	i := 0
	for {
		if i >= len(s) || s[i] < '0' || s[i] > '9' {
			if n == 0 {
				return ErrInvalidParseVersion
			}
			break
		}
		num, next, err := scanVersionPart(s, i)
		if err != nil {
			return err
		}
		if n < 3 {
			parts[n] = num
		}
		n++
		i = next
		if i >= len(s) || (s[i] != '.' && s[i] != '_') {
			break
		}
		if n >= 3 {
			break
		}
		i++
		if i < len(s) && (s[i] < '0' || s[i] > '9') {
			return ErrInvalidParseVersion
		}
	}

	*v = Version{Major: parts[0], Minor: parts[1], Patch: parts[2]}
	return nil
}

// scanVersionPart reads a decimal integer at s[i]. i is on a digit.
// A value above 65535 is an error.
func scanVersionPart(s string, i int) (num uint16, next int, err error) {
	const maxPart = 65535
	var acc int
	for i < len(s) {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		d := int(c - '0')
		if acc > (maxPart-d)/10 {
			return 0, i, ErrInvalidParseVersion
		}
		acc = acc*10 + d
		i++
	}
	return uint16(acc), i, nil
}
