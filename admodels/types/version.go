package types

import (
	"database/sql/driver"
	"fmt"
)

var ErrInvalidParseVersion = fmt.Errorf("invalid parse version")

// Version model description of standard version type like: 1.2.3
type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

func ParseVersion(str string) (Version, error) {
	var v Version
	if err := v.SetFromStr(str); err != nil {
		return v, err
	}
	return v, nil
}

func MustParseVersion(str string) Version {
	v, err := ParseVersion(str)
	if err != nil {
		panic(err)
	}
	return v
}

func IgnoreParseVersion(str string) Version {
	var v Version
	_ = v.SetFromStr(str)
	return v
}

func (v *Version) String() string {
	if v.Patch > 0 {
		return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	}
	if v.Minor > 0 {
		return fmt.Sprintf("%d.%d", v.Major, v.Minor)
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

func (v *Version) IsEmpty() bool {
	return v.Major == 0 && v.Minor == 0 && v.Patch == 0
}

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

func (v *Version) MarshalJSON() ([]byte, error) {
	return []byte(`"` + v.String() + `"`), nil
}

func (v *Version) UnmarshalJSON(data []byte) error {
	s := string(data)
	if len(s) > 1 && s[0] == '"' && s[len(s)-1] == '"' {
		return v.SetFromStr(s[1 : len(s)-1])
	}
	return ErrInvalidParseVersion
}

func (v *Version) SetFromStr(str string) error {
	var major, minor, patch int = 0, 0, 0
	if str != "" && str != "0" && str != "undefined" {
		if _, err := fmt.Sscanf(str, "%d.%d.%d", &major, &minor, &patch); err != nil {
			if _, err = fmt.Sscanf(str, "%d.%d", &major, &minor); err != nil {
				if _, err = fmt.Sscanf(str, "%d", &major); err != nil {
					return err
				}
			}
		}
	}
	v.Major = major
	v.Minor = minor
	v.Patch = patch
	return nil
}
