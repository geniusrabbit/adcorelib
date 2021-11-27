package eventstream

import (
	"context"
	"time"

	"geniusrabbit.dev/corelib/adtype"
	"github.com/geniusrabbit/notificationcenter"
)

// WinNotifier redeclared type
type WinNotifier struct {
	p notificationcenter.Publisher
}

// WinNotifications returns win notifier wrapper
func WinNotifications(p notificationcenter.Publisher) *WinNotifier {
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
