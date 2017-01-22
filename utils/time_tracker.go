package utils

import (
	"time"
	"fmt"
)

func TimeTrack(start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf(" in %.3f seconds\n", elapsed.Seconds())
}