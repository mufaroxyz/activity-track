package main

import (
	"activity-track/lib"
	"time"
)

func main() {
	if lib.WinApi == nil {
		panic("Failed to initialize WinApi")
	}

	var activityPayload = lib.ActivityPayload{}

	var lastMousePos = lib.POINT{}
	var lastKeyboardEvent = lib.KBDLLHOOKSTRUCT{}
	mousePosChannel := make(chan lib.CursorPosData)
	mouseClickChannel := make(chan lib.MSLLHOOKSTRUCTExtended, 10)
	keyboardEventChannel := make(chan lib.KBDLLHOOKSTRUCT, 10)
	go lib.MousePosTrack(mousePosChannel)
	go lib.MouseClickTrack(mouseClickChannel)
	go lib.KeyboardEventTrack(keyboardEventChannel)

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
		case mouseClick := <-mouseClickChannel:
			activityPayload.MouseClicks = append(activityPayload.MouseClicks, mouseClick)
		case keyboardEvent := <-keyboardEventChannel:
			if keyboardEvent.VkCode == lastKeyboardEvent.VkCode {
				continue
			}

			lastKeyboardEvent = keyboardEvent
		case <-ticker.C:
			lib.SaveDataInDb(activityPayload)
			println("Freeing up payload memory")
			activityPayload = lib.ActivityPayload{}
		}
	}
}
