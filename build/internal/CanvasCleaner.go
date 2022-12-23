package internal

import (
	"time"
)

func CanvasCleaner(canvas *Canvas) {
	// Clear canvas every day at 00:00
	now := time.Now()
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

	time.Sleep(nextMidnight.Sub(now))
	for {
		canvas.clear()
		time.Sleep(24 * time.Hour)
	}
}
