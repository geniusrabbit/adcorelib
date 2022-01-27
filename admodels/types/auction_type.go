package types

import (
	"bytes"
	"database/sql/driver"

	"github.com/geniusrabbit/gosql"
	"github.com/pkg/errors"
)

// AuctionType data
// CREATE TYPE AuctionType AS ENUM ('undefined', 'first_price', 'second_price')
type AuctionType uint8

// Auction types
const (
	UndefinedAuctionType   AuctionType = 0
	FirstPriceAuctionType  AuctionType = 1
	SecondPriceAuctionType AuctionType = 2
)

// IsFirtsPrice auction type
func (at AuctionType) IsFirtsPrice() bool {
	return at == FirstPriceAuctionType
}

// IsSecondPrice auction type
func (at AuctionType) IsSecondPrice() bool {
	return at == SecondPriceAuctionType
}

// Name of the status
func (at AuctionType) Name() string {
	switch at {
	case FirstPriceAuctionType:
		return `first_price`
	case SecondPriceAuctionType:
		return `second_price`
	}
	return `undefined`
}

// AuctionTypeNameToType name to const
func AuctionTypeNameToType(name string) AuctionType {
	switch name {
	case `first_price`, `first`, `1`:
		return FirstPriceAuctionType
	case `second_proce`, `second`, `2`:
		return SecondPriceAuctionType
	}
	return UndefinedAuctionType
}

// Value implements the driver.Valuer interface, json field interface
func (at AuctionType) Value() (driver.Value, error) {
	return at.Name(), nil
}

// Scan implements the driver.Valuer interface, json field interface
func (at *AuctionType) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		*at = AuctionTypeNameToType(v[1 : len(v)-1])
	case []byte:
		*at = AuctionTypeNameToType(string(v[1 : len(v)-1]))
	case nil:
		*at = UndefinedAuctionType
	default:
		return gosql.ErrInvalidScan
	}
	return nil
}

// MarshalJSON implements the json.Marshaler
func (at AuctionType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + at.Name() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaller
func (at *AuctionType) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errors.Wrap(errInvalidUnmarshalValue, "`"+string(b)+"`")
	}
	if bytes.HasPrefix(b, []byte(`"`)) {
		*at = AuctionTypeNameToType(string(b[1 : len(b)-1]))
	} else {
		*at = AuctionTypeNameToType(string(b))
	}
	return nil
}

// Implements the Unmarshaler interface of the yaml pkg.
func (at *AuctionType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var yamlString string
	if err := unmarshal(&yamlString); err != nil {
		return err
	}
	*at = AuctionTypeNameToType(yamlString)
	return nil
}
