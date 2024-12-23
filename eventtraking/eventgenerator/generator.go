//
// @project GeniusRabbit corelib 2018 - 2019, 2024
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019, 2024
//

package eventgenerator

import (
	"errors"

	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/eventtraking/events"
)

// Errors set
var (
	ErrInvalidMultipleItemAsSingle = errors.New("can`t convert multipleitem to single action")
)

// Generator object
type Generator[EventT EventType, UserInfoT UserInfoType] interface {
	// Event object by response
	Event(event events.Type, status uint8, response adtype.Responser, it adtype.ResponserItem) (EventT, error)

	// Events object list
	Events(event events.Type, status uint8, response adtype.Responser, it adtype.ResponserItemCommon) []EventT

	// UserInfo event object by response
	UserInfo(response adtype.Responser, it adtype.ResponserItem) (UserInfoT, error)
}

type generator[EventT EventType, UserInfoT UserInfoType] struct {
	service           string
	eventAllocator    Allocator[EventT]
	userInfoAllocator Allocator[UserInfoT]
}

// New generator object
func New[EventT EventType, UserInfoT UserInfoType](
	service string,
	eventAllocator Allocator[EventT],
	userInfoAllocator Allocator[UserInfoT],
) Generator[EventT, UserInfoT] {
	return generator[EventT, UserInfoT]{
		service:           service,
		eventAllocator:    eventAllocator,
		userInfoAllocator: userInfoAllocator,
	}
}

// Event object by response
func (g generator[EventT, UserInfoT]) Event(event events.Type, status uint8, response adtype.Responser, it adtype.ResponserItem) (EventT, error) {
	eventObj := g.eventAllocator()
	if err := eventObj.Fill(g.service, event, status, response, it); err != nil {
		return eventObj, err
	}
	return eventObj, nil
}

// Events object list
func (g generator[EventT, UserInfoT]) Events(event events.Type, status uint8, response adtype.Responser, it adtype.ResponserItemCommon) (events []EventT) {
	if mit, _ := it.(adtype.ResponserMultipleItem); mit != nil {
		ads := mit.Ads()
		events = make([]EventT, 0, len(ads))
		for _, it := range ads {
			if event, err := g.Event(event, status, response, it); err == nil {
				events = append(events, event)
			}
		}
	} else if event, err := g.Event(event, status, response, it.(adtype.ResponserItem)); err == nil {
		events = append(events, event)
	}
	return events
}

// UserInfo event object by response
func (g generator[EventT, UserInfoT]) UserInfo(response adtype.Responser, it adtype.ResponserItem) (UserInfoT, error) {
	userInfoObj := g.userInfoAllocator()
	if err := userInfoObj.Fill(response, it); err != nil {
		return userInfoObj, err
	}
	return userInfoObj, nil
}
