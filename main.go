package main

import (
	"activity-track/lib"
	"fmt"
	"syscall"
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
	mouseEventChannel := make(chan lib.MSLLHOOKSTRUCTExtended, 10)
	keyboardEventChannel := make(chan lib.KBDLLHOOKSTRUCT, 10)
	activeWindowEventChannel := make(chan lib.ActiveWindowEvent)
	go lib.MousePosTrack(mousePosChannel)
	go lib.MouseClickTrack(mouseEventChannel)
	go lib.KeyboardEventTrack(keyboardEventChannel)
	go lib.TrackWindowReplaced(activeWindowEventChannel)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case mousePos := <-mousePosChannel:
			if !lib.IsMouseMoved(lib.CursorPosData{POINT: lastMousePos}, mousePos) {
				continue
			}

			if len(activityPayload.CursorPositions) > 0 {
				lastMousePos = activityPayload.CursorPositions[len(activityPayload.CursorPositions)-1].POINT
			}

			activityPayload.CursorPositions = append(activityPayload.CursorPositions, mousePos)
		case mouseClick := <-mouseEventChannel:
			activityPayload.MouseClicks = append(activityPayload.MouseClicks, mouseClick)
		case keyboardEvent := <-keyboardEventChannel:
			if keyboardEvent.VkCode == lastKeyboardEvent.VkCode {
				continue
			}

			lastKeyboardEvent = keyboardEvent
		case activeWindowEvent := <-activeWindowEventChannel:
			buffer := make([]uint16, 256)
			lib.GetWindowTextW(activeWindowEvent.WindowHandle, &buffer[0], 256)
			windowTitle := syscall.UTF16ToString(buffer)

			println(fmt.Sprintf("<-activeWindowEventChannel ts: %v, handle: %v, title: %s", activeWindowEvent.TimeStamp, activeWindowEvent.WindowHandle, windowTitle))
		case <-ticker.C:
			lib.SaveDataInDb(activityPayload)
			// println("Freeing up payload memory")
			activityPayload = lib.ActivityPayload{}
		}
	}
}
