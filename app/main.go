package app

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
	hook "github.com/robotn/gohook"
	"github.com/vilsol/oshabi/config"
	"github.com/vilsol/oshabi/cv"
	"github.com/vilsol/oshabi/data"
	"github.com/vilsol/oshabi/types"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"image"
	"math"
	"strconv"
	"strings"
	"time"
)

const DISPLAY = 0 // TODO Replace with config

func InitializeApp(ctx context.Context) {
	hook.Register(hook.KeyDown, []string{"ctrl", "j"}, func(e hook.Event) {
		config.AddListings(ctx, ReadFull(ctx))
	})

	go func() {
		s := hook.Start()
		<-hook.Process(s)

		hook.Start()
	}()
}

func ReadFull(ctx context.Context) []types.ParsedListing {
	runtime.EventsEmit(ctx, "reading_listings")
	defer runtime.EventsEmit(ctx, "listings_read")

	if !ScrollTop() {
		// TODO Return error
		return nil
	}

	allListings := make([]types.ParsedListing, 0)

	offset := 0
	limit := 5
	for limit > 0 {
		allListings = append(allListings, ReadScreen(offset, limit)...)

		scrollCount := 0
		for i := 0; i < 5; i++ {
			img := CaptureScreen()
			if cv.CanScrollDown(img) {
				scrollCount++
				ScrollDown(img)
			} else {
				break
			}
		}

		offset = 5 - scrollCount
		limit = scrollCount
	}

	return allListings
}

func CaptureScreen() image.Image {
	bounds := screenshot.GetDisplayBounds(DISPLAY)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}
	return img
}

func ReadScreen(offset int, limit int) []types.ParsedListing {
	img := CaptureScreen()

	listings := cv.ExtractToListings(img, offset, limit)

	realListings := make([]types.ParsedListing, 0)
	for _, listing := range listings {
		count := cv.OCRListingCount(listing.Count)
		countInt, err := strconv.ParseInt(count, 10, 32)
		if err != nil || countInt <= 0 {
			continue
		}

		level := cv.OCRListingLevel(listing.Level)
		splitLevel := strings.Split(level, " ")
		levelInt, err := strconv.ParseInt(splitLevel[len(splitLevel)-1], 10, 32)
		if err != nil {
			continue
		}

		craft := data.FindCraft(cv.OCRListing(listing.Text))

		realListings = append(realListings, types.ParsedListing{
			Type:  data.GetCraftByText(craft).Type,
			Count: int(countInt),
			Level: int(levelInt),
		})
	}

	return realListings
}

func ScrollTop() bool {
	img := CaptureScreen()

	scrollUp, _, err := image.Decode(bytes.NewReader(data.ScrollUpPNG))
	if err != nil {
		panic(err)
	}

	point, pointValue, err := cv.Find(img, scrollUp)
	if err != nil {
		panic(err)
	}

	fmt.Println(pointValue)
	if pointValue < 0.9 {
		return false
	}

	realX, realY := TranslateCoordinates(point.X, point.Y)

	realX += scrollUp.Bounds().Dx() / 2
	realY += scrollUp.Bounds().Dy() * 2

	robotgo.Move(realX, realY)
	time.Sleep(time.Millisecond * 100)
	robotgo.Click()
	time.Sleep(time.Millisecond * 100)
	robotgo.Click()
	time.Sleep(time.Millisecond * 100)
	robotgo.Move(realX+100, realY)
	time.Sleep(time.Millisecond * 100)

	return true
}

func ScrollDown(img image.Image) {
	scrollDown, _, err := image.Decode(bytes.NewReader(data.ScrollDownPNG))
	if err != nil {
		panic(err)
	}

	point, _, err := cv.Find(img, scrollDown)
	if err != nil {
		panic(err)
	}

	realX, realY := TranslateCoordinates(point.X, point.Y)

	realX += scrollDown.Bounds().Dx() / 2
	realY += scrollDown.Bounds().Dy() / 2

	robotgo.Move(realX, realY)
	time.Sleep(time.Millisecond * 25)
	robotgo.Click()
	time.Sleep(time.Millisecond * 25)
	robotgo.Move(realX+100, realY)
	time.Sleep(time.Millisecond * 25)
}

func TranslateCoordinates(x int, y int) (int, int) {
	bounds := screenshot.GetDisplayBounds(DISPLAY)

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

//func ParseListing(img image.Image, position int) *ParsedListing {
//	count := cv.OCRListingCount(listing.Count)
//	countInt, err := strconv.ParseInt(count, 10, 32)
//	if err != nil || countInt <= 0 {
//		return nil
//	}
//
//	level := cv.OCRListingLevel(listing.Level)
//	splitLevel := strings.Split(level, " ")
//	levelInt, err := strconv.ParseInt(splitLevel[len(splitLevel)-1], 10, 32)
//	if err != nil {
//		return nil
//	}
//
//	craft := data.FindCraft(cv.OCRListing(listing.Text))
//
//	return &ParsedListing{
//		Type:  data.GetCraftByText(craft).Type,
//		Count: int(countInt),
//		Level: int(levelInt),
//	}
//}
