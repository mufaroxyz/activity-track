package lib

import (
	"time"
)

var (
	innerActiveWindowEventChannel chan<- ActiveWindowEvent
)

func WinEventProc(hWinEventHook HWINEVENTHOOK, event DWORD, hwnd HWND,
	idObject LONG, idChild LONG, idEventThread DWORD, dwmsEventTime DWORD) uintptr {

	activeWindowEvent := ActiveWindowEvent{
		WindowHandle: hwnd,
		TimeStamp:    time.Now().Unix(),
	}
	innerActiveWindowEventChannel <- activeWindowEvent
	return 0
}

func TrackWindowReplaced(ch chan<- ActiveWindowEvent) {
	innerActiveWindowEventChannel = ch

	hook := SetWinEventHook(EVENT_SYSTEM_FOREGROUND, EVENT_SYSTEM_FOREGROUND,
		0, WinEventProc, 0, 0, WINEVENT_OUTOFCONTEXT)

	defer UnhookWinEvent(hook)

	var msg MSG
	for {
		if !GetMessageW(&msg, 0, 0, 0) {
			break
		}
	}
}
