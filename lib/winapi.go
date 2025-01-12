package lib

import (
	"syscall"
	"unsafe"
)

var (
	WinApi           = syscall.NewLazyDLL("user32.dll")
	procGetCursorPos = WinApi.NewProc("GetCursorPos")
)

func GetProcAddress(name string) uintptr {
	proc := WinApi.NewProc(name)
	return proc.Addr()
}

func GetCursorPos(lpPoint *POINT) bool {
	ret, _, _ := procGetCursorPos.Call(uintptr(unsafe.Pointer(lpPoint)))
	return ret != 0
}
