package lib

type POINT struct {
	X, Y int32
}

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
