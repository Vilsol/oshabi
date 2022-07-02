package cv

import (
	"image"
	"image/color"
	"image/draw"
	"regexp"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/kbinani/screenshot"
	"github.com/otiai10/gosseract/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/vilsol/oshabi/config"
	"github.com/vilsol/oshabi/data"
	"github.com/vilsol/oshabi/types"
)

func OCRListing(img image.Image) (string, error) {
	return ocr(img, GetWhitelist(), gosseract.PSM_SINGLE_BLOCK, true, false)
}

func OCRListingCount(img image.Image) (string, error) {
	return ocr(img, "1234567890", gosseract.PSM_SINGLE_CHAR, false, true)
}

func OCRListingLevel(img image.Image) (string, error) {
	return ocr(img, "1234567890", gosseract.PSM_SINGLE_LINE, false, true)
}

type RawListing struct {
	Count image.Image
	Text  image.Image
	Level image.Image
}

const (
	listingHeight = float64(174)

	infoListingHorizontalOffset = 55
	infoListingVerticalOffset   = 170

	listingTextOffset         = 170
	listingTextWidth          = 870
	listingTextHeight         = 175
	listingTextVerticalOffset = -10

	levelWidth       = 180
	levelHeight      = 53
	levelNumberWidth = 65

	countWidth            = 38
	countHeight           = 60
	countHorizontalOffset = 8
	countVerticalOffset   = -5

	groveOffset = -32
)

func FindInfoButton(img image.Image) (image.Point, error) {
	infoButtonLocation, _, err := Find(img, Scale(data.InfoButton))
	if err != nil {
		return image.Point{}, errors.Wrap(err, "failed to find info button")
	}
	log.Debug().Int("x", infoButtonLocation.X).Int("y", infoButtonLocation.Y).Msg("info button found")
	return infoButtonLocation, err
}

func IsInGrove(img image.Image) (bool, error) {
	_, hortValue, err := Find(img, Scale(data.Horticrafting))
	if err != nil {
		return false, errors.Wrap(err, "failed to find horticrafting label")
	}
	log.Debug().Bool("grove", hortValue < 0.8).Msg("in grove or horticrafting")
	return hortValue < 0.8, nil
}

func ExtractToListings(img image.Image, offset int, limit int) ([]RawListing, error) {
	infoButtonLocation, err := FindInfoButton(img)
	if err != nil {
		return nil, err
	}

	inGrove, err := IsInGrove(img)
	if err != nil {
		return nil, err
	}

	listings := make([]RawListing, limit)

	for i := 0; i < limit; i++ {
		pxOffset := int(float64(i+offset)*ScaleNf(listingHeight) + float64(ScaleN(infoListingVerticalOffset)))

		if inGrove {
			pxOffset += ScaleN(groveOffset)
		}

		listingLeft := infoButtonLocation.X + ScaleN(infoListingHorizontalOffset)
		listingTop := infoButtonLocation.Y + pxOffset + ScaleN(listingTextVerticalOffset)

		listingTextLeft := listingLeft + ScaleN(listingTextOffset)
		listingTextRight := listingTextLeft + ScaleN(listingTextWidth)
		listingTextBottom := listingTop + ScaleN(listingTextHeight)

		text := imaging.Crop(img, image.Rect(
			listingTextLeft,
			listingTop,
			listingTextRight,
			listingTextBottom,
		))

		blackRect := image.Rect(text.Bounds().Dx()-ScaleN(levelWidth), text.Bounds().Dy()-ScaleN(levelHeight), text.Bounds().Dx(), text.Bounds().Dy())
		draw.Draw(text, blackRect, &image.Uniform{C: color.RGBA{A: 255}}, image.Point{}, draw.Src)

		count := imaging.Crop(img, image.Rect(
			listingLeft+ScaleN(countHorizontalOffset),
			listingTop+ScaleN(countVerticalOffset),
			listingLeft+ScaleN(countWidth)+ScaleN(countHorizontalOffset),
			listingTop+ScaleN(countHeight)+ScaleN(countVerticalOffset),
		))

		level := imaging.Crop(img, image.Rect(
			listingTextRight-ScaleN(levelNumberWidth),
			listingTextBottom-ScaleN(levelHeight),
			listingTextRight,
			listingTextBottom,
		))

		listings[i] = RawListing{
			Count: count,
			Text:  text,
			Level: level,
		}
	}

	return listings, nil
}

func CanScrollDown(infoButtonLocation image.Point, inGrove bool, img image.Image) (bool, error) {
	bounds := screenshot.GetDisplayBounds(config.Get().Display)

	pxOffset := int(5*ScaleNf(listingHeight) + float64(ScaleN(infoListingVerticalOffset)))

	if inGrove {
		pxOffset += ScaleN(groveOffset)
	}

	listingLeft := infoButtonLocation.X + ScaleN(infoListingHorizontalOffset)
	listingTop := infoButtonLocation.Y + pxOffset

	newRect := image.Rect(
		bounds.Min.X+listingLeft-ScaleN(20),
		bounds.Min.Y+listingTop-ScaleN(20),
		bounds.Min.X+listingLeft+ScaleN(30),
		bounds.Min.Y+listingTop+ScaleN(30),
	)

	if img == nil {
		nextCount, err := screenshot.CaptureRect(newRect)

		if err != nil {
			return false, errors.Wrap(err, "failed capturing screen")
		}

		img = nextCount
	} else {
		img = imaging.Crop(img, newRect)
	}

	_, cornerVal, err := Find(img, Scale(data.CountCorner))
	if err != nil {
		return false, errors.Wrap(err, "failed to find count corner")
	}

	return cornerVal >= 0.7, nil
}

var levelRegex = regexp.MustCompile(`\w*(\d\s*\d)`)

func ReadImage(img image.Image, offset int, limit int) ([]types.ParsedListing, error) {
	listings, err := ExtractToListings(img, offset, limit)
	if err != nil {
		return nil, err
	}

	realListings := make([]types.ParsedListing, 0)
	for _, listing := range listings {
		count, err := OCRListingCount(listing.Count)
		if err != nil {
			return nil, err
		}

		countInt, err := strconv.ParseInt(count, 10, 32)
		if err != nil || countInt <= 0 {
			continue
		}

		level, err := OCRListingLevel(listing.Level)
		if err != nil {
			return nil, err
		}

		matches := levelRegex.FindAllStringSubmatch(level, -1)
		if len(matches) == 0 {
			continue
		}

		levelClean := strings.Replace(matches[0][1], " ", "", -1)

		if levelClean[0] == '1' {
			levelClean = "7" + levelClean[1:]
		}

		levelInt, err := strconv.ParseInt(levelClean, 10, 32)
		if err != nil {
			continue
		}

		listingText, err := OCRListing(listing.Text)
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
