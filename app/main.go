package app

import (
	"context"
	"image"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
	"github.com/pkg/errors"
	hook "github.com/robotn/gohook"
	"github.com/vilsol/oshabi/config"
	"github.com/vilsol/oshabi/cv"
	"github.com/vilsol/oshabi/data"
	"github.com/vilsol/oshabi/types"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func InitializeApp(ctx context.Context) {
	hook.Register(hook.KeyDown, []string{"ctrl", "j"}, func(e hook.Event) {
		listings, err := ReadFull(ctx)
		if err != nil {
			runtime.EventsEmit(ctx, "error", err.Error())
			return
		}

		if err := config.AddListings(ctx, listings); err != nil {
			runtime.EventsEmit(ctx, "error", err.Error())
			return
		}
	})

	go func() {
		s := hook.Start()
		<-hook.Process(s)

		hook.Start()
	}()
}

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

	offset := 0
	limit := 5
	for limit > 0 {
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
			Type:  data.GetCraftByText(craft).Type,
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

	point, pointValue, err := cv.Find(img, data.ScrollUp)
	if err != nil {
		return false, errors.Wrap(err, "failed to find scroll up button")
	}

	if pointValue < 0.9 {
		return false, nil
	}

	realX, realY := TranslateCoordinates(point.X, point.Y)

	realX += data.ScrollUp.Bounds().Dx() / 2
	realY += data.ScrollUp.Bounds().Dy() * 2

	robotgo.Move(realX, realY)
	time.Sleep(time.Millisecond * 100)
	robotgo.Click()
	time.Sleep(time.Millisecond * 100)
	robotgo.Click()
	time.Sleep(time.Millisecond * 100)
	robotgo.Move(realX+100, realY)
	time.Sleep(time.Millisecond * 100)

	return true, nil
}

func ScrollDown(img image.Image) error {
	point, _, err := cv.Find(img, data.ScrollDown)
	if err != nil {
		return errors.Wrap(err, "failed to find scroll down button")
	}

	realX, realY := TranslateCoordinates(point.X, point.Y)

	realX += data.ScrollDown.Bounds().Dx() / 2
	realY += data.ScrollDown.Bounds().Dy() / 2

	robotgo.Move(realX, realY)
	time.Sleep(time.Millisecond * 25)
	robotgo.Click()
	time.Sleep(time.Millisecond * 25)
	robotgo.Move(realX+100, realY)
	time.Sleep(time.Millisecond * 25)

	return nil
}

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

	realX := int(math.Abs(float64(leftMost))) + x
	realY := int(math.Abs(float64(topMost))) + y

	return realX, realY
}
