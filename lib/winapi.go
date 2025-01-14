package lib

import (
	"path/filepath"
	"syscall"
	"unsafe"
)

var (
	User32                         = syscall.NewLazyDLL("user32.dll")
	kernel32                       = syscall.NewLazyDLL("kernel32.dll")
	procQueryFullProcessImageNameW = kernel32.NewProc("QueryFullProcessImageNameW")
	procGetWindowThreadProcessId   = User32.NewProc("GetWindowThreadProcessId")
	procGetCursorPos               = User32.NewProc("GetCursorPos")
	procSetWindowsHookExW          = User32.NewProc("SetWindowsHookExW")
	procUnhookWindowsHook          = User32.NewProc("UnhookWindowsHook")
	procCallNextHookEx             = User32.NewProc("CallNextHookEx")
	procGetMessageW                = User32.NewProc("GetMessageW")
	procSetWinEventHook            = User32.NewProc("SetWinEventHook")
	procUnhookWinEvent             = User32.NewProc("UnhookWinEvent")
	procGetWindowTextW             = User32.NewProc("GetWindowTextW")
	procGetWindowInfo              = User32.NewProc("GetWindowInfo")
	procGetWindowModuleFileNameW   = User32.NewProc("GetWindowModuleFileNameW")
	HookHandle                     HHOOK
)

func GetProcAddress(name string) uintptr {
	proc := User32.NewProc(name)
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

func GetWindowInfo(hwnd HWND, pwi *WINDOWINFO) bool {
	ret, _, _ := procGetWindowInfo.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(pwi)))
	return ret != 0
}

func GetWindowModuleFileNameW(hwnd HWND, lpszFileName LPWSTR, cchFileNameMax UINT) UINT {
	ret, _, _ := procGetWindowModuleFileNameW.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(lpszFileName)),
		uintptr(cchFileNameMax))
	return UINT(ret)
}

func GetWindowThreadProcessId(hwnd syscall.Handle, processId *uint32) uint32 {
	ret, _, _ := procGetWindowThreadProcessId.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(processId)),
	)
	return uint32(ret)
}

func GetProcessExeName(hwnd syscall.Handle) string {
	var processID uint32
	GetWindowThreadProcessId(hwnd, &processID)

	hProcess, err := syscall.OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION, false, processID)
	if err != nil {
		return ""
	}
	defer syscall.CloseHandle(hProcess)

	var pathLen uint32 = 260
	var buffer = make([]uint16, pathLen)

	ret, _, _ := syscall.NewLazyDLL("psapi.dll").NewProc("GetProcessImageFileNameW").Call(
		uintptr(hProcess),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(pathLen),
	)

	if ret == 0 {
		return ""
	}

	path := syscall.UTF16ToString(buffer[:])
	return filepath.Base(path)
}

func GetCursorPos(lpPoint *POINT) bool {
	ret, _, _ := procGetCursorPos.Call(uintptr(unsafe.Pointer(lpPoint)))
	return ret != 0
}
