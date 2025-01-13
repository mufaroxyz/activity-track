package lib

import "time"

func SaveDataInDb(payload ActivityPayload) {
	// TODO: implement the logic to process data to more optimized format & save it in db
	finalData := ActivityPayloadFinal{}

	finalData.SnapshotTime = time.Now().UnixNano()

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
}
