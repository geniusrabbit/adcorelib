package eventgenerator

import (
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/eventtraking/events"
)

// EventUnpacFunc object for event unpack function
type EventUnpacFunc = events.EventUnpacFunc

// Allocator object
type Allocator[T any] func() T

// EventType object for event basic type interface
type EventType interface {
	// SetDateTime set date time of event
	SetDateTime(t int64)

	// EventType returns type of event
	EventType() events.Type

	// EventURL returns url of event target
	EventURL() string

	// PrepareURL prepare url for event
	PrepareURL(url string) string

	// SetEventPurchaseViewPrice set price of event
	SetEventPurchaseViewPrice(price int64) error

	// Pack event object to byte array
	Pack() events.Code

	// Unpack event object from byte array
	Unpack(data []byte, unpuckFnc ...EventUnpacFunc) error

	// Fill event object by response
	Fill(service string, event events.Type, status uint8, response adtype.Response, it adtype.ResponseItem) error
}

// LeadType object for lead basic type interface
type LeadType interface {
	// String returns string representation of object
	String() string

	// EventAuctionID returns auction id of event
	EventAuctionID() string

	// SetDateTime set date time of event
	SetDateTime(t int64)

	// Pack event object to byte array
	Pack() events.Code

	// Unpack event object from byte array
	Unpack(data []byte) error

	// Fill event object by response
	Fill(item adtype.ResponseItem, response adtype.Response) error
}

// UserInfoType object for user info basic type interface
type UserInfoType interface {
	Fill(response adtype.Response, it adtype.ResponseItem) error
}

// TestEvent object for testing
type TestEvent struct{}

// SetDateTime set date time of event
func (e *TestEvent) SetDateTime(t int64) {}

// EventType returns type of event
func (e *TestEvent) EventType() events.Type { return events.View }

// EventURL returns url of event target
func (e *TestEvent) EventURL() string { return "" }

// PrepareURL prepare url for event
func (e *TestEvent) PrepareURL(url string) string { return "" }

// SetEventPurchaseViewPrice set price of event
func (e *TestEvent) SetEventPurchaseViewPrice(price int64) error { return nil }

// Pack event object to byte array
func (e *TestEvent) Pack() events.Code { return events.CodeObj([]byte{1, 2, 3}, nil) }

// Unpack event object from byte array
func (e *TestEvent) Unpack(data []byte, unpuckFnc ...EventUnpacFunc) error { return nil }

// Fill event object by response
func (e *TestEvent) Fill(service string, event events.Type, status uint8, response adtype.Response, it adtype.ResponseItem) error {
	return nil
}

var _ EventType = &TestEvent{}

// TestLead object for testing
type TestLead struct{}

// String returns string representation of object
func (l *TestLead) String() string { return "" }

// EventAuctionID returns auction id of event
func (l *TestLead) EventAuctionID() string { return "" }

// SetDateTime set date time of event
func (l *TestLead) SetDateTime(t int64) {}

// Pack event object to byte array
func (l *TestLead) Pack() events.Code { return events.CodeObj([]byte{1, 2, 3}, nil) }

// Unpack event object from byte array
func (l *TestLead) Unpack(data []byte) error { return nil }

// Fill event object by response
func (l *TestLead) Fill(item adtype.ResponseItem, response adtype.Response) error { return nil }

var _ LeadType = &TestLead{}

// TestUserInfo object for testing
type TestUserInfo struct{}

// Fill user info object by response
func (u *TestUserInfo) Fill(response adtype.Response, it adtype.ResponseItem) error { return nil }

var _ UserInfoType = &TestUserInfo{}
