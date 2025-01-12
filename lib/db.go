package lib

func SaveDataInDb(payload ActivityPayload) {
	// TODO: implement the logic to process data to more optimized format & save it in db
	finalData := ActivityPayloadFinal{}

	if len(payload.CursorPositions) > 2 {
		for i := 1; i < len(payload.CursorPositions); i++ {
			var pixelDist = pixelDistance(payload.CursorPositions[i].X, payload.CursorPositions[i].Y, payload.CursorPositions[i-1].X, payload.CursorPositions[i-1].Y)
			finalData.TotalMouseDistance += pixelsToMeters(pixelDist)
		}
	}

	println("Total mouse distance: ", finalData.TotalMouseDistance)
}
