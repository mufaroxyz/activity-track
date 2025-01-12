package main

import (
	"activity-track/lib"
	"time"
)

func main() {
	if lib.WinApi == nil {
		panic("Failed to initialize WinApi")
	}

	var getCursorPosAddr = lib.GetProcAddress("GetCursorPos")
	if getCursorPosAddr == 0 {
		panic("Failed to get GetCursorPos address")
	}

	println("GetCursorPos addr: ", getCursorPosAddr)

	var activityPayload = lib.ActivityPayload{}

	var lastMousePos = lib.POINT{}
	mousePosChannel := make(chan lib.CursorPosData)
	go lib.MousePosTrack(mousePosChannel)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case mousePos := <-mousePosChannel:
			if !lib.IsMouseMoved(lib.CursorPosData{POINT: lastMousePos}, mousePos) {
				continue
			}

			println(mousePos.X, mousePos.Y, mousePos.TimeStamp)

			if len(activityPayload.CursorPositions) > 0 {
				lastMousePos = activityPayload.CursorPositions[len(activityPayload.CursorPositions)-1].POINT
			}

			activityPayload.CursorPositions = append(activityPayload.CursorPositions, mousePos)
		case <-ticker.C:
			lib.SaveDataInDb(activityPayload)
			println("Freeing up payload memory")
			activityPayload = lib.ActivityPayload{}
		}
	}
}
