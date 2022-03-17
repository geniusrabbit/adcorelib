package eventstream

import (
	"context"
	"time"

	"geniusrabbit.dev/corelib/adtype"
	nc "github.com/geniusrabbit/notificationcenter/v2"
)

// WinNotifier redeclared type
type WinNotifier struct {
	p nc.Publisher
}

// WinNotifications returns win notifier wrapper
func WinNotifications(p nc.Publisher) *WinNotifier {
	return &WinNotifier{p: p}
}

// Send URL win notify
func (w *WinNotifier) Send(ctx context.Context, url string) error {
	return w.p.Publish(ctx, &adtype.WinEvent{URL: url, Time: time.Now()})
}

// SendEvent win notify
func (w *WinNotifier) SendEvent(ctx context.Context, event *adtype.WinEvent) error {
	return w.p.Publish(ctx, event)
}
