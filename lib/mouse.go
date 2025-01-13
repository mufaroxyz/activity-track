package lib

import (
	"time"
	"unsafe"
)

var (
	mouseClickChannel chan<- MSLLHOOKSTRUCT
)

func MousePosTrack(ch chan<- CursorPosData) {
	for {
		POINT := &POINT{}
		CursorPosData := &CursorPosData{}
		GetCursorPos(POINT)
		CursorPosData.POINT = *POINT
		CursorPosData.TimeStamp = time.Now().UnixNano()
		ch <- *CursorPosData
		time.Sleep(100 * time.Millisecond)
	}
}

func LowLevelMouseProc(nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	if nCode < 0 {
		return CallNextHookEx(0, nCode, wParam, lParam)
	}

	if nCode >= 0 && wParam == WM_LBUTTONDOWN || wParam == WM_RBUTTONDOWN {
		mouseStruct := (*MSLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		mouseClickChannel <- *mouseStruct
	}

	return CallNextHookEx(0, nCode, wParam, lParam)
}

func MouseClickTrack(ch chan<- MSLLHOOKSTRUCT) {
	println("Hooking mouse events")
	mouseClickChannel = ch
	hook := SetWindowsHookExW(WH_MOUSE_LL, LowLevelMouseProc, 0, 0)
	if hook == 0 {
		panic("Failed to set hook")
	}
	defer UnhookWindowsHook(hook)

	var msg MSG
	for {
		if !GetMessageW(&msg, 0, 0, 0) {
			break
		}
	}
}

func pixelDistance(x1, y1, x2, y2 int32) float64 {
	return float64((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2))
}

func IsMouseMoved(prevPos, newPos CursorPosData) bool {
	return pixelDistance(prevPos.X, prevPos.Y, newPos.X, newPos.Y) > 100
}

func pixelsToMeters(pixels float64) float64 {
	return pixels * 0.0002645833
}
