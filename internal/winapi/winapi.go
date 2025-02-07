package winapi

import (
	"activity-track/pkg"
	"path/filepath"
	"syscall"
	"unsafe"
)

var (
	User32                       = syscall.NewLazyDLL("user32.dll")
	procGetWindowThreadProcessId = User32.NewProc("GetWindowThreadProcessId")
	procGetCursorPos             = User32.NewProc("GetCursorPos")
	procSetWindowsHookExW        = User32.NewProc("SetWindowsHookExW")
	procUnhookWindowsHook        = User32.NewProc("UnhookWindowsHook")
	procCallNextHookEx           = User32.NewProc("CallNextHookEx")
	procGetMessageW              = User32.NewProc("GetMessageW")
	procSetWinEventHook          = User32.NewProc("SetWinEventHook")
	procUnhookWinEvent           = User32.NewProc("UnhookWinEvent")
	procGetWindowTextW           = User32.NewProc("GetWindowTextW")
	procGetWindowInfo            = User32.NewProc("GetWindowInfo")
	procGetWindowModuleFileNameW = User32.NewProc("GetWindowModuleFileNameW")
	HookHandle                   pkg.HHOOK
)

func GetProcAddress(name string) uintptr {
	proc := User32.NewProc(name)
	return proc.Addr()
}

func SetWindowsHookExW(idHook int, lpfn pkg.HOOKPROC, hMod pkg.HINSTANCE, dwThreadId pkg.DWORD) (pkg.HHOOK, error) {
	ret, _, err := procSetWindowsHookExW.Call(
		uintptr(idHook),
		uintptr(syscall.NewCallback(lpfn)),
		uintptr(hMod),
		uintptr(dwThreadId))

	return pkg.HHOOK(ret), err
}

func UnhookWindowsHook(hhk pkg.HHOOK) bool {
	ret, _, _ := procUnhookWindowsHook.Call(uintptr(hhk))
	return ret != 0
}

func CallNextHookEx(hhk pkg.HHOOK, nCode int, wParam pkg.WPARAM, lParam pkg.LPARAM) pkg.LRESULT {
	ret, _, _ := procCallNextHookEx.Call(
		uintptr(hhk),
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam))
	return pkg.LRESULT(ret)
}

func GetMessageW(lpMsg *pkg.MSG, hWnd pkg.HWND, wMsgFilterMin, wMsgFilterMax pkg.UINT) bool {
	ret, _, _ := procGetMessageW.Call(
		uintptr(unsafe.Pointer(lpMsg)),
		uintptr(hWnd),
		uintptr(wMsgFilterMin),
		uintptr(wMsgFilterMax))
	return ret != 0
}

func SetWinEventHook(eventMin, eventMax pkg.DWORD, hmodWinEventProc pkg.HMODULE, pfnWinEventProc pkg.WINEVENTPROC, idProcess pkg.DWORD, idThread pkg.DWORD, dwFlags pkg.DWORD) pkg.HWINEVENTHOOK {
	ret, _, _ := procSetWinEventHook.Call(
		uintptr(eventMin),
		uintptr(eventMax),
		uintptr(hmodWinEventProc),
		syscall.NewCallback(pfnWinEventProc),
		uintptr(idProcess),
		uintptr(idThread),
		uintptr(dwFlags))
	return pkg.HWINEVENTHOOK(ret)
}

func UnhookWinEvent(hWinEventHook pkg.HWINEVENTHOOK) bool {
	ret, _, _ := procUnhookWinEvent.Call(uintptr(hWinEventHook))
	return ret != 0
}

func GetWindowTextW(hWnd pkg.HWND, lpString pkg.LPWSTR, nMaxCount int32) int32 {
	ret, _, _ := procGetWindowTextW.Call(
		uintptr(hWnd),
		uintptr(unsafe.Pointer(lpString)),
		uintptr(nMaxCount))
	return int32(ret)
}

func GetWindowInfo(hwnd pkg.HWND, pwi *pkg.WINDOWINFO) bool {
	ret, _, _ := procGetWindowInfo.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(pwi)))
	return ret != 0
}

func GetWindowModuleFileNameW(hwnd pkg.HWND, lpszFileName pkg.LPWSTR, cchFileNameMax pkg.UINT) pkg.UINT {
	ret, _, _ := procGetWindowModuleFileNameW.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(lpszFileName)),
		uintptr(cchFileNameMax))
	return pkg.UINT(ret)
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

	hProcess, err := syscall.OpenProcess(pkg.PROCESS_QUERY_LIMITED_INFORMATION, false, processID)
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

func GetCursorPos(lpPoint *pkg.POINT) bool {
	ret, _, _ := procGetCursorPos.Call(uintptr(unsafe.Pointer(lpPoint)))
	return ret != 0
}
