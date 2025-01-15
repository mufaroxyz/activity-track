package main

import (
	"activity-track/lib"
	"fmt"
	"strings"
	"syscall"
	"time"
)

func main() {
	if lib.User32 == nil {
		panic("Failed to initialize WinApi")
	}

	if lib.DEBUG == 0 {
		lib.InitConfig()
	}

	lib.SetupCloudflareClient()

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

	ticker := time.NewTicker(20 * time.Second)
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
			activityPayload.KeyboardPresses = append(activityPayload.KeyboardPresses, keyboardEvent)
		case activeWindowEvent := <-activeWindowEventChannel:
			buffer := make([]uint16, 256)
			lib.GetWindowTextW(activeWindowEvent.WindowHandle, &buffer[0], 256)
			windowTitle := syscall.UTF16ToString(buffer)

			if lib.IsTitleIgnored(windowTitle) {
				continue
			}

			processName := lib.GetProcessExeName(syscall.Handle(activeWindowEvent.WindowHandle))
			association := lib.GetAssociation(strings.ToLower(processName), windowTitle)

			println(fmt.Sprintf("<-activeWindowEventChannel ts: %v, group: %v, process_name: %s", activeWindowEvent.TimeStamp, association, processName))

			activityPayload.WindowActivities = append(activityPayload.WindowActivities, lib.WindowActivity{
				Activity:  association,
				TimeStamp: activeWindowEvent.TimeStamp,
			})
		case <-ticker.C:
			lib.SaveDataInDb(activityPayload)
			// println("Freeing up payload memory")
			activityPayload = lib.ActivityPayload{}
		}
	}
}
