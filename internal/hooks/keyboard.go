package hooks

import (
	"activity-track/internal/winapi"
	"activity-track/pkg"
	"unsafe"
)

var (
	innerKeyboardChannel chan<- pkg.KBDLLHOOKSTRUCT
)

func LowLevelKeyboardProc(nCode int, wParam pkg.WPARAM, lParam pkg.LPARAM) pkg.LRESULT {
	if nCode < 0 {
		return winapi.CallNextHookEx(0, nCode, wParam, lParam)
	}

	if nCode >= 0 {
		if wParam != pkg.WM_KEYDOWN && wParam != pkg.WM_SYSKEYDOWN {
			kbdStruct := (*pkg.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
			innerKeyboardChannel <- *kbdStruct
		}
	}

	return winapi.CallNextHookEx(0, nCode, wParam, lParam)
}

func KeyboardEventTrack(ch chan<- pkg.KBDLLHOOKSTRUCT) {
	println("Hooking keyboard events")
	innerKeyboardChannel = ch
	hook, _ := winapi.SetWindowsHookExW(pkg.WH_KEYBOARD_LL, LowLevelKeyboardProc, 0, 0)
	if hook == 0 {
		panic("Failed to set hook WH_KEYBOARD_LL")
	}
	defer winapi.UnhookWindowsHook(hook)

	var msg pkg.MSG
	for {
		if !winapi.GetMessageW(&msg, 0, 0, 0) {
			break
		}
	}
}
