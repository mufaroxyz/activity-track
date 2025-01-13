package lib

type (
	HOOKPROC  func(int, WPARAM, LPARAM) LRESULT
	HINSTANCE uintptr
	HWND      uintptr
	WPARAM    uintptr
	LPARAM    uintptr
	LRESULT   uintptr
	DWORD     uint32
	UINT      uint32
	MSG       struct {
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
	WH_MOUSE_LL    = 14
	WM_LBUTTONDOWN = 0x0201
	WM_RBUTTONDOWN = 0x0204
)

type MSLLHOOKSTRUCT struct {
	Point     POINT
	MouseData uint32
	Flags     uint32
	Time      uint32
	ExtraInfo uintptr
}

type HHOOK uintptr

type CursorPosData struct {
	POINT
	TimeStamp int64
}

type ActivityPayload struct {
	CursorPositions []CursorPosData
}

type ActivityPayloadFinal struct {
	TotalMouseDistance float64
}
