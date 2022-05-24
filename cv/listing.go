package cv

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/otiai10/gosseract/v2"

	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"github.com/vilsol/oshabi/data"
)

func OCRListing(img image.Image) (string, error) {
	return ocr(img, "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ,.%' ", gosseract.PSM_AUTO)
}

func OCRListingCount(img image.Image) (string, error) {
	return ocr(img, "1234567890", gosseract.PSM_SINGLE_CHAR)
}

func OCRListingLevel(img image.Image) (string, error) {
	return ocr(img, "Levl1234567890 ", gosseract.PSM_SINGLE_LINE)
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

	listingTextOffset = 170
	listingTextWidth  = 870
	listingTextHeight = 165

	levelWidth  = 180
	levelHeight = 48

	countWidth            = 38
	countHeight           = 45
	countHorizontalOffset = 8
	countVerticalOffset   = 5

	groveOffset = -32
)

func listingTrackers(img image.Image) (bool, image.Point, error) {
	_, hortValue, err := Find(img, Scale(data.Horticrafting))
	if err != nil {
		return false, image.Point{}, errors.Wrap(err, "failed to find horticrafting label")
	}

	infoButtonLocation, _, err := Find(img, Scale(data.InfoButton))
	if err != nil {
		return false, image.Point{}, errors.Wrap(err, "failed to find info button")
	}

	return hortValue < 0.8, infoButtonLocation, nil
}

func ExtractToListings(img image.Image, offset int, limit int) ([]RawListing, error) {
	inGrove, infoButtonLocation, err := listingTrackers(img)
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
		listingTop := infoButtonLocation.Y + pxOffset

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
			listingTextRight-ScaleN(levelWidth),
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

func CanScrollDown(img image.Image) (bool, error) {
	_, infoButtonLocation, err := listingTrackers(img)
	if err != nil {
		return false, err
	}

	pxOffset := int(5*ScaleNf(listingHeight) + float64(ScaleN(infoListingVerticalOffset)))

	listingLeft := infoButtonLocation.X + ScaleN(infoListingHorizontalOffset)
	listingTop := infoButtonLocation.Y + pxOffset

	nextCount := imaging.Crop(img, image.Rect(
		listingLeft-ScaleN(20),
		listingTop-ScaleN(20),
		listingLeft+ScaleN(30),
		listingTop+ScaleN(30),
	))

	_, cornerVal, err := Find(nextCount, Scale(data.CountCorner))
	if err != nil {
		return false, errors.Wrap(err, "failed to find count corner")
	}

	return cornerVal >= 0.9, nil
}
