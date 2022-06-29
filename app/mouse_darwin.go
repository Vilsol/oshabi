package app

import (
	"math"

	"github.com/kbinani/screenshot"

	"github.com/vilsol/oshabi/config"
)

func TranslateCoordinates(x int, y int) (int, int) {
	bounds := screenshot.GetDisplayBounds(config.Get().Display)

	leftMost := bounds.Min.X
	topMost := bounds.Min.Y
	for i := 0; i < screenshot.NumActiveDisplays(); i++ {
		b := screenshot.GetDisplayBounds(i)
		if b.Min.X < leftMost {
			leftMost = b.Min.X
		}

		if b.Min.Y < topMost {
			topMost = b.Min.Y
		}
	}

	realX := int(math.Abs(float64(leftMost))) + x + bounds.Min.X
	realY := int(math.Abs(float64(topMost))) + y + bounds.Min.Y

	return realX, realY
}
