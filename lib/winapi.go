package lib

import (
	"syscall"
	"unsafe"
)

var (
	WinApi                = syscall.NewLazyDLL("user32.dll")
	procGetCursorPos      = WinApi.NewProc("GetCursorPos")
	procSetWindowsHookExW = WinApi.NewProc("SetWindowsHookExW")
	procUnhookWindowsHook = WinApi.NewProc("UnhookWindowsHook")
	procLowLevelMouseProc = WinApi.NewProc("LowLevelMouseProc")
	procCallNextHookEx    = WinApi.NewProc("CallNextHookEx")
	procGetMessageW       = WinApi.NewProc("GetMessageW")
)

func GetProcAddress(name string) uintptr {
	proc := WinApi.NewProc(name)
	return proc.Addr()
}

func SetWindowsHookExW(idHook int, lpfn HOOKPROC, hMod HINSTANCE, dwThreadId DWORD) HHOOK {
	ret, _, _ := procSetWindowsHookExW.Call(
		uintptr(idHook),
		uintptr(syscall.NewCallback(lpfn)),
		uintptr(hMod),
		uintptr(dwThreadId))
	return HHOOK(ret)
}

func UnhookWindowsHook(hhk HHOOK) bool {
	ret, _, _ := procUnhookWindowsHook.Call(uintptr(hhk))
	return ret != 0
}

func LowLevelMouseProc(nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _, _ := procLowLevelMouseProc.Call(
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam))
	return LRESULT(ret)
}

func CallNextHookEx(hhk HHOOK, nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _, _ := procCallNextHookEx.Call(
		uintptr(hhk),
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam))
	return LRESULT(ret)
}

func GetMessageW(lpMsg *MSG, hWnd HWND, wMsgFilterMin, wMsgFilterMax UINT) bool {
	ret, _, _ := procGetMessageW.Call(
		uintptr(unsafe.Pointer(lpMsg)),
		uintptr(hWnd),
		uintptr(wMsgFilterMin),
		uintptr(wMsgFilterMax))
	return ret != 0
}

func GetCursorPos(lpPoint *POINT) bool {
	ret, _, _ := procGetCursorPos.Call(uintptr(unsafe.Pointer(lpPoint)))
	return ret != 0
}
