package app

import (
	"github.com/kbinani/screenshot"
	"github.com/vilsol/oshabi/config"
)

func TranslateCoordinates(x int, y int) (int, int) {
	bounds := screenshot.GetDisplayBounds(config.Get().Display)
	return bounds.Min.X + x, bounds.Min.Y + y
}
