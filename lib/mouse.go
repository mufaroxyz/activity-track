package lib

import "time"

func MousePosTrack(ch chan<- CursorPosData) {
	for {
		POINT := &POINT{}
		CursorPosData := &CursorPosData{}
		GetCursorPos(POINT)
		CursorPosData.POINT = *POINT
		CursorPosData.TimeStamp = time.Now().UnixNano()
		ch <- *CursorPosData
		time.Sleep(100 * time.Millisecond)
	}
}

func pixelDistance(x1, y1, x2, y2 int32) float64 {
	return float64((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2))
}

func IsMouseMoved(prevPos, newPos CursorPosData) bool {
	return pixelDistance(prevPos.X, prevPos.Y, newPos.X, newPos.Y) > 100
}

func pixelsToMeters(pixels float64) float64 {
	return pixels * 0.0002645833
}
