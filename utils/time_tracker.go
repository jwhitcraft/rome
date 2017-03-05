package utils

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func TimeTrack(start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf(" %v %v %v",
		color.GreenString("in"),
		color.CyanString("%.3f", elapsed.Seconds()),
		color.GreenString("seconds\n"))
}
