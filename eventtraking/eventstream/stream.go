//
// @project GeniusRabbit corelib 2018 - 2019, 2022, 2024
// @author Dmitry Ponomarev <demdxx@gmail.com>
//

package eventstream

import (
	"context"
	"errors"

	nc "github.com/geniusrabbit/notificationcenter/v2"

	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventgenerator"
	"github.com/geniusrabbit/adcorelib/eventtraking/events"
)

var errInvalidResponse = errors.New(`response object can't be nil`)

type (
	EventType    = eventgenerator.EventType
	UserInfoType = eventgenerator.UserInfoType
)

// Stream accessor interface for the event tracking
type Stream interface {
	// SendEvent native action
	SendEvent(ctx context.Context, event any) error

	// Send response
	Send(event events.Type, status uint8, response adtype.Responser, it adtype.ResponserItem) error

	// SendLeadEvent as lead code type
	SendLeadEvent(ctx context.Context, event any) error

	// SendSourceSkip event for the response
	SendSourceSkip(response adtype.Responser) error

	// SendSourceNoBid event for the response
	SendSourceNoBid(response adtype.Responser) error

	// SendSourceFail event for the response
	SendSourceFail(response adtype.Responser) error

	// SendAccessPointBid event for the response
	SendAccessPointBid(response adtype.Responser, it ...adtype.ResponserItem) error

	// SendAccessPointSkip event for the response
	SendAccessPointSkip(response adtype.Responser) error

	// SendAccessPointNoBid event for the response
	SendAccessPointNoBid(response adtype.Responser) error

	// SendAccessPointFail event for the response
	SendAccessPointFail(response adtype.Responser) error
}

type stream[EventT EventType, UserInfoT UserInfoType] struct {
	events    nc.Publisher
	userInfo  nc.Publisher
	generator eventgenerator.Generator[EventT, UserInfoT]
}

// New stream object
func New[EventT EventType, UserInfoT UserInfoType](events, userInfo nc.Publisher, generator eventgenerator.Generator[EventT, UserInfoT]) Stream {
	return &stream[EventT, UserInfoT]{
		events:    events,
		userInfo:  userInfo,
		generator: generator,
	}
}

// SendEvent native action
func (s *stream[EventT, UserInfoT]) SendEvent(ctx context.Context, event any) error {
	return s.events.Publish(ctx, event)
}

// Send response
func (s *stream[EventT, UserInfoT]) Send(event events.Type, status uint8, response adtype.Responser, it adtype.ResponserItem) (err error) {
	if response == nil {
		return errInvalidResponse
	}

	var (
		info UserInfoT
		ctx  = response.Context()
	)

	for _, event := range s.generator.Events(event, status, response, it) {
		if err = s.SendEvent(ctx, event); err != nil {
			return err
		}
	}

	if s.userInfo != nil { // Send user info if it's possible
		if info, err = s.generator.UserInfo(response, it); err == nil {
			err = s.userInfo.Publish(ctx, info)
		}
	}
	return err
}

// SendLeadEvent as lead code type
func (s *stream[EventT, UserInfoT]) SendLeadEvent(ctx context.Context, event any) error {
	return s.events.Publish(ctx, event)
}

// SendSourceSkip event for the response
func (s *stream[EventT, UserInfoT]) SendSourceSkip(response adtype.Responser) error {
	return s.Send(events.SourceSkip, events.StatusUndefined, response, (*adtype.ResponseItemEmpty)(nil))
}

// SendSourceNoBid event for the response
func (s *stream[EventT, UserInfoT]) SendSourceNoBid(response adtype.Responser) error {
	return s.Send(events.SourceNoBid, events.StatusUndefined, response, (*adtype.ResponseItemEmpty)(nil))
}

// SendSourceFail event for the response
func (s *stream[EventT, UserInfoT]) SendSourceFail(response adtype.Responser) error {
	return s.Send(events.SourceFail, events.StatusFailed, response, (*adtype.ResponseItemEmpty)(nil))
}

// SendAccessPointBid event for the response
func (s *stream[EventT, UserInfoT]) SendAccessPointBid(response adtype.Responser, it ...adtype.ResponserItem) error {
	for _, item := range it {
		err := s.Send(events.AccessPointBid, events.StatusSuccess, response, item)
		if err != nil {
			return err
		}
	}
	return nil
}

// SendAccessPointSkip event for the response
func (s *stream[EventT, UserInfoT]) SendAccessPointSkip(response adtype.Responser) error {
	return s.Send(events.AccessPointSkip, events.StatusUndefined, response, (*adtype.ResponseItemEmpty)(nil))
}

// SendAccessPointNoBid event for the response
func (s *stream[EventT, UserInfoT]) SendAccessPointNoBid(response adtype.Responser) error {
	return s.Send(events.AccessPointNoBid, events.StatusUndefined, response, (*adtype.ResponseItemEmpty)(nil))
}

// SendAccessPointFail event for the response
func (s *stream[EventT, UserInfoT]) SendAccessPointFail(response adtype.Responser) error {
	return s.Send(events.AccessPointFail, events.StatusFailed, response, (*adtype.ResponseItemEmpty)(nil))
}
