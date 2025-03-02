package lib

import "unsafe"

var (
	innerKeyboardChannel chan<- KBDLLHOOKSTRUCT
)

func LowLevelKeyboardProc(nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	if nCode < 0 {
		return CallNextHookEx(0, nCode, wParam, lParam)
	}

	if nCode >= 0 {
		if wParam != WM_KEYDOWN && wParam != WM_SYSKEYDOWN {
			kbdStruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
			innerKeyboardChannel <- *kbdStruct
		}
	}

	return CallNextHookEx(0, nCode, wParam, lParam)
}

func KeyboardEventTrack(ch chan<- KBDLLHOOKSTRUCT) {
	println("Hooking keyboard events")
	innerKeyboardChannel = ch
	hook, _ := SetWindowsHookExW(WH_KEYBOARD_LL, LowLevelKeyboardProc, 0, 0)
	if hook == 0 {
		panic("Failed to set hook WH_KEYBOARD_LL")
	}
	defer UnhookWindowsHook(hook)

	var msg MSG
	for {
		if !GetMessageW(&msg, 0, 0, 0) {
			break
		}
	}
}
