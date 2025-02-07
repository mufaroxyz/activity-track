package hooks

import (
	"activity-track/internal/winapi"
	"activity-track/pkg"
	"time"
)

var (
	innerActiveWindowEventChannel chan<- pkg.ActiveWindowEvent
)

func WinEventProc(hWinEventHook pkg.HWINEVENTHOOK, event pkg.DWORD, hwnd pkg.HWND,
	idObject pkg.LONG, idChild pkg.LONG, idEventThread pkg.DWORD, dwmsEventTime pkg.DWORD) uintptr {

	activeWindowEvent := pkg.ActiveWindowEvent{
		WindowHandle: hwnd,
		TimeStamp:    time.Now().Unix(),
	}
	innerActiveWindowEventChannel <- activeWindowEvent
	return 0
}

func TrackWindowReplaced(ch chan<- pkg.ActiveWindowEvent) {
	innerActiveWindowEventChannel = ch

	hook := winapi.SetWinEventHook(pkg.EVENT_SYSTEM_FOREGROUND, pkg.EVENT_SYSTEM_FOREGROUND,
		0, WinEventProc, 0, 0, pkg.WINEVENT_OUTOFCONTEXT)

	defer winapi.UnhookWinEvent(hook)

	var msg pkg.MSG
	for {
		if !winapi.GetMessageW(&msg, 0, 0, 0) {
			break
		}
	}
}
