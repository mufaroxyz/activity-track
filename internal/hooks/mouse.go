package hooks

import (
	"activity-track/internal/winapi"
	"activity-track/pkg"
	"time"
	"unsafe"
)

var (
	innerMouseEventChannel chan<- pkg.MSLLHOOKSTRUCTExtended
)

func MousePosTrack(ch chan<- pkg.CursorPosData) {
	for {
		POINT := &pkg.POINT{}
		CursorPosData := &pkg.CursorPosData{}
		winapi.GetCursorPos(POINT)
		CursorPosData.POINT = *POINT
		CursorPosData.TimeStamp = time.Now().UnixNano()
		ch <- *CursorPosData
		time.Sleep(100 * time.Millisecond)
	}
}

func LowLevelMouseProc(nCode int, wParam pkg.WPARAM, lParam pkg.LPARAM) pkg.LRESULT {
	if nCode < 0 {
		return winapi.CallNextHookEx(0, nCode, wParam, lParam)
	}

	if nCode >= 0 {
		mouseStruct := (*pkg.MSLLHOOKSTRUCTExtended)(unsafe.Pointer(lParam))

		switch wParam {
		case pkg.WM_LBUTTONDOWN:
			mouseStruct.ButtonType = pkg.WM_LBUTTONDOWN
		case pkg.WM_RBUTTONDOWN:
			mouseStruct.ButtonType = pkg.WM_RBUTTONDOWN
		}

		if wParam == pkg.WM_LBUTTONDOWN || wParam == pkg.WM_RBUTTONDOWN {
			innerMouseEventChannel <- *mouseStruct
		}
	}

	return winapi.CallNextHookEx(0, nCode, wParam, lParam)
}

func MouseClickTrack(ch chan<- pkg.MSLLHOOKSTRUCTExtended) {
	println("Hooking mouse events")
	innerMouseEventChannel = ch
	hook, _ := winapi.SetWindowsHookExW(pkg.WH_MOUSE_LL, LowLevelMouseProc, 0, 0)
	if hook == 0 {
		panic("Failed to set hook")
	}
	defer winapi.UnhookWindowsHook(hook)

	var msg pkg.MSG
	for {
		if !winapi.GetMessageW(&msg, 0, 0, 0) {
			break
		}
	}
}

func PixelDistance(x1, y1, x2, y2 int32) float64 {
	return float64((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2))
}

func IsMouseMoved(prevPos, newPos pkg.CursorPosData) bool {
	return PixelDistance(prevPos.X, prevPos.Y, newPos.X, newPos.Y) > 100
}

func PixelsToMeters(pixels float64) float64 {
	return pixels * 0.0002645833
}
