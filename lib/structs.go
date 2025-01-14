package lib

type (
	HOOKPROC      func(int, WPARAM, LPARAM) LRESULT
	WINEVENTPROC  func(HWINEVENTHOOK, DWORD, HWND, LONG, LONG, DWORD, DWORD) uintptr
	HWINEVENTHOOK uintptr
	LONG          int32
	HINSTANCE     uintptr
	HMODULE       uintptr
	LPCWSTR       *uint16
	HWND          uintptr
	WPARAM        uintptr
	LPARAM        uintptr
	LRESULT       uintptr
	DWORD         uint32
	UINT          uint32
	ULONG_PTR     uintptr
	LPWSTR        *uint16
	MSG           struct {
		HWND    HWND
		Message UINT
		WParam  WPARAM
		LParam  LPARAM
		Time    DWORD
		Pt      POINT
	}
)

type POINT struct {
	X, Y int32
}

const (
	WH_KEYBOARD_LL                               = 13
	WM_KEYDOWN                                   = 0x0100
	WM_SYSKEYDOWN                                = 0x0104
	WH_MOUSE_LL                                  = 14
	WM_LBUTTONDOWN                               = 0x0201
	WM_RBUTTONDOWN                               = 0x0204
	WH_SHELL                                     = 10
	HSHELL_WINDOWREPLACED                        = 13
	GET_MODULE_HANDLE_EX_FLAG_FROM_ADDRESS       = 0x00000004
	GET_MODULE_HANDLE_EX_FLAG_UNCHANGED_REFCOUNT = 0x00000002
	EVENT_SYSTEM_FOREGROUND                      = 0x0003
	WINEVENT_OUTOFCONTEXT                        = 0x0000
)

type MSLLHOOKSTRUCT struct {
	Point     POINT
	MouseData uint32
	Flags     uint32
	Time      uint32
	ExtraInfo uintptr
}

type MSLLHOOKSTRUCTExtended struct {
	MSLLHOOKSTRUCT
	ButtonType int
}

type KBDLLHOOKSTRUCT struct {
	VkCode      DWORD
	ScanCode    DWORD
	Flags       DWORD
	Time        DWORD
	DwExtraInfo ULONG_PTR
}

type HHOOK uintptr

type CursorPosData struct {
	POINT
	TimeStamp int64
}

type ActiveWindowEvent struct {
	WindowHandle HWND
	TimeStamp    int64
}

type ActivityPayload struct {
	CursorPositions []CursorPosData
	MouseClicks     []MSLLHOOKSTRUCTExtended
	KeyboardPresses []KBDLLHOOKSTRUCT
}

type MouseActivity struct {
	TotalMouseDistance float64
	RightClicks        int
	LeftClicks         int
}

type ActivityPayloadFinal struct {
	MouseActivity   MouseActivity
	KeyboardPresses int
	SnapshotTime    int64
}
