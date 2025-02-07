package main

import (
	"activity-track/internal/cloudflare"
	"activity-track/internal/config"
	"activity-track/internal/db"
	"activity-track/internal/hooks"
	"activity-track/internal/winapi"
	"activity-track/pkg"
	"fmt"
	"strings"
	"syscall"
	"time"
)

func main() {
	if winapi.User32 == nil {
		panic("Failed to initialize WinApi")
	}

	config.InitConfig()
	cloudflare.SetupCloudflareClient()

	var activityPayload = pkg.ActivityPayload{}

	var lastMousePos = pkg.POINT{}
	var lastKeyboardEvent = pkg.KBDLLHOOKSTRUCT{}
	var lastWindowActivity = pkg.WindowActivity{}
	mousePosChannel := make(chan pkg.CursorPosData)
	mouseEventChannel := make(chan pkg.MSLLHOOKSTRUCTExtended, 10)
	keyboardEventChannel := make(chan pkg.KBDLLHOOKSTRUCT, 10)
	activeWindowEventChannel := make(chan pkg.ActiveWindowEvent)
	go hooks.MousePosTrack(mousePosChannel)
	go hooks.MouseClickTrack(mouseEventChannel)
	go hooks.KeyboardEventTrack(keyboardEventChannel)
	go hooks.TrackWindowReplaced(activeWindowEventChannel)

	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case mousePos := <-mousePosChannel:
			if !hooks.IsMouseMoved(pkg.CursorPosData{POINT: lastMousePos}, mousePos) {
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
			winapi.GetWindowTextW(activeWindowEvent.WindowHandle, &buffer[0], 256)
			windowTitle := syscall.UTF16ToString(buffer)

			if pkg.IsTitleIgnored(windowTitle) {
				continue
			}

			processName := winapi.GetProcessExeName(syscall.Handle(activeWindowEvent.WindowHandle))
			association := pkg.GetAssociation(strings.ToLower(processName), windowTitle)

			println(fmt.Sprintf("<-activeWindowEventChannel ts: %v, group: %v, process_name: %s", activeWindowEvent.TimeStamp, association, processName))

			activityPayload.WindowActivities = append(activityPayload.WindowActivities, pkg.WindowActivity{
				Activity:  association,
				TimeStamp: activeWindowEvent.TimeStamp,
			})
		case <-ticker.C:
			db.SaveDataInDb(activityPayload)
			// println("Freeing up payload memory")
			lastWindowActivity = activityPayload.WindowActivities[len(activityPayload.WindowActivities)-1]
			activityPayload = pkg.ActivityPayload{}
			activityPayload.WindowActivities = append(activityPayload.WindowActivities, lastWindowActivity)
		}
	}
}
