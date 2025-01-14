package lib

import (
	"syscall"
	"unsafe"
)

var (
	WinApi                 = syscall.NewLazyDLL("user32.dll")
	kernel32               = syscall.NewLazyDLL("kernel32.dll")
	procGetCursorPos       = WinApi.NewProc("GetCursorPos")
	procSetWindowsHookExW  = WinApi.NewProc("SetWindowsHookExW")
	procUnhookWindowsHook  = WinApi.NewProc("UnhookWindowsHook")
	procCallNextHookEx     = WinApi.NewProc("CallNextHookEx")
	procGetMessageW        = WinApi.NewProc("GetMessageW")
	procGetModuleHandleExW = kernel32.NewProc("GetModuleHandleExW")
	procGetModuleHandleW   = kernel32.NewProc("GetModuleHandleW")
	procSetWinEventHook    = WinApi.NewProc("SetWinEventHook")
	procUnhookWinEvent     = WinApi.NewProc("UnhookWinEvent")
	procGetWindowTextW     = WinApi.NewProc("GetWindowTextW")
	HookHandle             HHOOK
)

func GetProcAddress(name string) uintptr {
	proc := WinApi.NewProc(name)
	return proc.Addr()
}

func SetWindowsHookExW(idHook int, lpfn HOOKPROC, hMod HINSTANCE, dwThreadId DWORD) (HHOOK, error) {
	ret, _, err := procSetWindowsHookExW.Call(
		uintptr(idHook),
		uintptr(syscall.NewCallback(lpfn)),
		uintptr(hMod),
		uintptr(dwThreadId))

	return HHOOK(ret), err
}

func UnhookWindowsHook(hhk HHOOK) bool {
	ret, _, _ := procUnhookWindowsHook.Call(uintptr(hhk))
	return ret != 0
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

func GetModuleHandleExW(dwFlags DWORD, lpModuleName LPCWSTR, phModule *HMODULE) (bool, error) {
	ret, _, err := procGetModuleHandleExW.Call(
		uintptr(dwFlags),
		uintptr(unsafe.Pointer(lpModuleName)),
		uintptr(unsafe.Pointer(phModule)))
	return ret != 0, err
}

func GetModuleHandleW(lpModuleName LPCWSTR) HMODULE {
	ret, _, _ := procGetModuleHandleW.Call(uintptr(unsafe.Pointer(lpModuleName)))
	return HMODULE(ret)
}

func SetWinEventHook(eventMin, eventMax DWORD, hmodWinEventProc HMODULE, pfnWinEventProc WINEVENTPROC, idProcess DWORD, idThread DWORD, dwFlags DWORD) HWINEVENTHOOK {
	ret, _, _ := procSetWinEventHook.Call(
		uintptr(eventMin),
		uintptr(eventMax),
		uintptr(hmodWinEventProc),
		syscall.NewCallback(pfnWinEventProc),
		uintptr(idProcess),
		uintptr(idThread),
		uintptr(dwFlags))
	return HWINEVENTHOOK(ret)
}

func UnhookWinEvent(hWinEventHook HWINEVENTHOOK) bool {
	ret, _, _ := procUnhookWinEvent.Call(uintptr(hWinEventHook))
	return ret != 0
}

func GetWindowTextW(hWnd HWND, lpString LPWSTR, nMaxCount int32) int32 {
	ret, _, _ := procGetWindowTextW.Call(
		uintptr(hWnd),
		uintptr(unsafe.Pointer(lpString)),
		uintptr(nMaxCount))
	return int32(ret)
}

func GetCursorPos(lpPoint *POINT) bool {
	ret, _, _ := procGetCursorPos.Call(uintptr(unsafe.Pointer(lpPoint)))
	return ret != 0
}
