package notify

import (
	"time"

	"github.com/godbus/dbus"
	"golang.org/x/xerrors"
)

func Send(summary string, body string, expire time.Duration) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return xerrors.Errorf("notify send: %w", err)
	}
	obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	call := obj.Call(
		"org.freedesktop.Notifications.Notify",
		0,
		"",
		uint32(0),
		"",
		summary,
		body,
		[]string{},
		map[string]dbus.Variant{},
		int32(expire)/1e6)
	if call.Err != nil {
		return xerrors.Errorf("notify send: %w", call.Err)
	}
	if err := conn.Close(); err != nil {
		return xerrors.Errorf("notify send: %w", err)
	}
	return nil
	//var ret uint32
	//err := call.Store(&ret)
	//if err != nil {
	//log.Printf("error getting uint32 ret value: %v", err)
	//return ret, err
	//}
	//return ret, nil
}
