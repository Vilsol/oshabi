package app

import (
	"context"
	"fmt"
	"image"
	"strconv"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
	"github.com/pkg/errors"
	"github.com/vilsol/oshabi/config"
	"github.com/vilsol/oshabi/cv"
	"github.com/vilsol/oshabi/data"
	"github.com/vilsol/oshabi/types"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func ReadFull(ctx context.Context) ([]types.ParsedListing, error) {
	runtime.EventsEmit(ctx, "reading_listings")
	defer runtime.EventsEmit(ctx, "listings_read")

	canScroll, err := ScrollTop()
	if err != nil {
		return nil, err
	}

	if !canScroll {
		return nil, errors.New("cannot scroll up")
	}

	allListings := make([]types.ParsedListing, 0)

	totalScrollCount := 0
	offset := 0
	limit := 5
	for limit > 0 && totalScrollCount < 20 {
		listings, err := ReadScreen(offset, limit)
		if err != nil {
			return nil, err
		}

		allListings = append(allListings, listings...)

		scrollCount := 0
		for i := 0; i < 5; i++ {
			img, err := CaptureScreen()
			if err != nil {
				return nil, err
			}

			canScrollDown, err := cv.CanScrollDown(img)
			if err != nil {
				return nil, err
			}

			if canScrollDown {
				scrollCount++
				if err := ScrollDown(img); err != nil {
					return nil, err
				}
			} else {
				break
			}
		}

		offset = 5 - scrollCount
		limit = scrollCount
		totalScrollCount += scrollCount
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

func ReadScreen(offset int, limit int) ([]types.ParsedListing, error) {
	img, err := CaptureScreen()
	if err != nil {
		return nil, err
	}

	return ReadImage(img, offset, limit)
}

func ReadImage(img image.Image, offset int, limit int) ([]types.ParsedListing, error) {
	listings, err := cv.ExtractToListings(img, offset, limit)
	if err != nil {
		return nil, err
	}

	realListings := make([]types.ParsedListing, 0)
	for _, listing := range listings {
		count, err := cv.OCRListingCount(listing.Count)
		if err != nil {
			return nil, err
		}

		countInt, err := strconv.ParseInt(count, 10, 32)
		if err != nil || countInt <= 0 {
			continue
		}

		level, err := cv.OCRListingLevel(listing.Level)
		if err != nil {
			return nil, err
		}

		splitLevel := strings.Split(level, " ")
		levelInt, err := strconv.ParseInt(splitLevel[len(splitLevel)-1], 10, 32)
		if err != nil {
			continue
		}

		listingText, err := cv.OCRListing(listing.Text)
		if err != nil {
			return nil, err
		}

		craft := data.FindCraft(listingText)

		realListings = append(realListings, types.ParsedListing{
			Type:  data.GetCraft(craft).Type,
			Count: int(countInt),
			Level: int(levelInt),
		})
	}

	return realListings, nil
}

func ScrollTop() (bool, error) {
	img, err := CaptureScreen()
	if err != nil {
		return false, err
	}

	point, pointValue, err := cv.Find(img, cv.Scale(data.InfoButton))
	if err != nil {
		return false, errors.Wrap(err, "failed to find info button")
	}

	if pointValue < 0.8 {
		return false, fmt.Errorf("info button not found: %f", pointValue)
	}

	realX, realY := TranslateCoordinates(point.X, point.Y)

	robotgo.Move(realX, realY)
	time.Sleep(time.Millisecond * 100)
	robotgo.Click()
	time.Sleep(time.Millisecond * 100)
	robotgo.Move(realX+cv.ScaleN(100), realY+cv.ScaleN(250))
	time.Sleep(time.Millisecond * 100)
	robotgo.Scroll(0, 20)
	time.Sleep(time.Millisecond * 100)
	robotgo.Move(realX+cv.ScaleN(100), realY)
	time.Sleep(time.Millisecond * 100)

	return true, nil
}

func ScrollDown(img image.Image) error {
	point, _, err := cv.Find(img, cv.Scale(data.InfoButton))
	if err != nil {
		return errors.Wrap(err, "failed to find info button")
	}

	realX, realY := TranslateCoordinates(point.X, point.Y)

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
