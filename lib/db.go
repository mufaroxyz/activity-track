package lib

import (
	"encoding/json"
	"fmt"
	"time"
)

/*
	d1 instance:
	{
		"snapshot_time": timestamp,
		"mouse_activity": { ... }
		...
	}
	15 min interval with 15 min cache on the frontend reading the data
*/

func SaveDataInDb(payload ActivityPayload) {
	finalData := ActivityPayloadFinal{}

	finalData.SnapshotTime = time.Now().Unix()

	if len(payload.CursorPositions) > 2 {
		for i := 1; i < len(payload.CursorPositions); i++ {
			var pixelDist = pixelDistance(payload.CursorPositions[i].X, payload.CursorPositions[i].Y, payload.CursorPositions[i-1].X, payload.CursorPositions[i-1].Y)
			finalData.MouseActivity.TotalMouseDistance += pixelsToMeters(pixelDist)
		}
	}

	finalData.MouseActivity.LeftClicks = 0
	finalData.MouseActivity.RightClicks = 0

	for _, click := range payload.MouseClicks {
		if click.ButtonType == WM_LBUTTONDOWN {
			finalData.MouseActivity.LeftClicks++
		} else if click.ButtonType == WM_RBUTTONDOWN {
			finalData.MouseActivity.RightClicks++
		}
	}

	finalData.KeyboardPresses = len(payload.KeyboardPresses)

	activityMap := make(map[string]int64)

	for i := 1; i < len(payload.WindowActivities); i++ {
		activityType := payload.WindowActivities[i-1].Activity
		var duration int64
		if i == len(payload.WindowActivities)-1 {
			duration = time.Now().Unix() - payload.WindowActivities[i].TimeStamp
		} else {
			duration = payload.WindowActivities[i].TimeStamp - payload.WindowActivities[i-1].TimeStamp
		}
		activityMap[activityType] += duration
	}

	mergedWindowActivities := []WindowActivityFinal{}
	for activity, duration := range activityMap {
		mergedWindowActivities = append(mergedWindowActivities, WindowActivityFinal{
			Activity: activity,
			Time:     duration,
		})
	}

	finalData.WindowActivities = mergedWindowActivities

	mouseActivityJson, _ := json.Marshal(finalData.MouseActivity)
	windowActivityJson, _ := json.Marshal(finalData.WindowActivities)

	queryResult, err := Query(fmt.Sprintf(`
		INSERT INTO activity (snapshot_time, mouse_activity, keyboard_presses, window_activity)
		VALUES (%d, '%+v', %d, '%+v')
	`, finalData.SnapshotTime, string(mouseActivityJson), finalData.KeyboardPresses, string(windowActivityJson)))

	if err != nil {
		panic(err)
	}

	println(fmt.Sprintf("Query results: %+v \n", queryResult))
}
