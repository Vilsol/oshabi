package app

import (
	"context"
	"image"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
	"github.com/pkg/errors"
	"github.com/vilsol/oshabi/config"
	"github.com/vilsol/oshabi/cv"
	"github.com/vilsol/oshabi/types"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func ReadFull(ctx context.Context) ([]types.ParsedListing, error) {
	runtime.EventsEmit(ctx, "reading_listings")
	defer runtime.EventsEmit(ctx, "listings_read")

	img, err := CaptureScreen()
	if err != nil {
		return nil, err
	}

	infoButtonLocation, err := cv.FindInfoButton(img)
	if err != nil {
		return nil, err
	}

	canScroll, err := ScrollTop(infoButtonLocation)
	if err != nil {
		return nil, err
	}

	if !canScroll {
		return nil, errors.New("cannot scroll up")
	}

	inGrove, err := cv.IsInGrove(img)
	if err != nil {
		return nil, err
	}

	allListings := make([]types.ParsedListing, 0)

	totalScrollCount := 0
	offset := 0
	limit := 5
	for limit > 0 && totalScrollCount < 20 {
		img, err = CaptureScreen()
		if err != nil {
			return nil, err
		}

		listingChan := make(chan struct {
			listings []types.ParsedListing
			err      error
		})
		go func(img image.Image, offset int, limit int) {
			listings, err := cv.ReadImage(img, offset, limit)
			listingChan <- struct {
				listings []types.ParsedListing
				err      error
			}{listings: listings, err: err}
		}(img, offset, limit)

		scrollCount := 0
		for i := 0; i < 5; i++ {
			canScrollDown, err := cv.CanScrollDown(infoButtonLocation, inGrove, nil)
			if err != nil {
				return nil, err
			}

			if canScrollDown {
				scrollCount++
				if err := ScrollDown(infoButtonLocation); err != nil {
					return nil, err
				}
			} else {
				break
			}
		}

		offset = 5 - scrollCount
		limit = scrollCount
		totalScrollCount += scrollCount

		response := <-listingChan
		close(listingChan)

		if response.err != nil {
			return nil, err
		}

		allListings = append(allListings, response.listings...)
	}

	return allListings, nil
}

func CaptureScreen() (image.Image, error) {
	bounds := screenshot.GetDisplayBounds(config.Get().Display)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, errors.Wrap(err, "failed capturing screen")
	}
	return img, nil
}

func ScrollTop(infoButtonLocation image.Point) (bool, error) {
	realX, realY := TranslateCoordinates(infoButtonLocation.X, infoButtonLocation.Y)

	robotgo.Move(realX, realY)
	time.Sleep(time.Millisecond * 100)
	robotgo.Click()
	time.Sleep(time.Millisecond * 200)
	robotgo.Move(realX+cv.ScaleN(100), realY+cv.ScaleN(250))
	time.Sleep(time.Millisecond * 100)
	robotgo.Scroll(0, 20)
	time.Sleep(time.Millisecond * 100)
	robotgo.Move(realX+cv.ScaleN(100), realY)
	time.Sleep(time.Millisecond * 100)

	return true, nil
}

func ScrollDown(infoButtonLocation image.Point) error {
	realX, realY := TranslateCoordinates(infoButtonLocation.X, infoButtonLocation.Y)

	robotgo.Move(realX, realY)
	time.Sleep(time.Millisecond * 100)
	robotgo.Click()
	time.Sleep(time.Millisecond * 100)
	robotgo.Move(realX+cv.ScaleN(100), realY+cv.ScaleN(250))
	time.Sleep(time.Millisecond * 100)
	robotgo.Scroll(0, -1)
	time.Sleep(time.Millisecond * 100)
	robotgo.Move(realX+cv.ScaleN(100), realY)
	time.Sleep(time.Millisecond * 100)

	return nil
}
