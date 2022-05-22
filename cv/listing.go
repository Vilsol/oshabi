package cv

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"github.com/vilsol/oshabi/data"
)

func OCRListing(img image.Image) (string, error) {
	return ocr(img, "1234567890abcdefghijklmnopqrstuvwxyz,.%' ")
}

func OCRListingCount(img image.Image) (string, error) {
	return ocr(img, "1234567890")
}

func OCRListingLevel(img image.Image) (string, error) {
	return ocr(img, "levl1234567890 ")
}

func PrepareForOCR(img image.Image) image.Image {
	src := imaging.Grayscale(img)
	src = imaging.Invert(src)
	return imaging.Sharpen(src, 4)
}

type RawListing struct {
	Count image.Image
	Text  image.Image
	Level image.Image
}

const (
	listingHeight = float64(170.5)
	textHeight    = listingHeight - 12

	scrollHorizontalOffset = 35
	scrollVerticalOffset   = -10

	levelWidth  = 180
	levelHeight = 40

	infoHorizontalOffset = 185
)

func listingTrackers(img image.Image) (image.Point, image.Point, error) {
	scrollLocation, _, err := Find(img, Scale(data.ScrollUp))
	if err != nil {
		return image.Point{}, image.Point{}, errors.Wrap(err, "failed to find scroll up button")
	}

	infoButtonLocation, _, err := Find(img, Scale(data.InfoButton))
	if err != nil {
		return image.Point{}, image.Point{}, errors.Wrap(err, "failed to find info button")
	}

	return scrollLocation, infoButtonLocation, nil
}

func ExtractToListings(img image.Image, offset int, limit int) ([]RawListing, error) {
	scrollLocation, infoButtonLocation, err := listingTrackers(img)
	if err != nil {
		return nil, err
	}

	listings := make([]RawListing, limit)

	for i := 0; i < limit; i++ {
		pxOffset := int(float64(i+offset)*ScaleNf(listingHeight) + float64(ScaleN(scrollVerticalOffset)))

		infoButtonBounds := ScaleBounds(data.InfoButton)
		text := imaging.Crop(img, image.Rect(
			infoButtonLocation.X+infoButtonBounds.Dx()+ScaleN(infoHorizontalOffset),
			scrollLocation.Y+pxOffset,
			scrollLocation.X-ScaleN(scrollHorizontalOffset),
			scrollLocation.Y+pxOffset+int(ScaleNf(textHeight)),
		))

		blackRect := image.Rect(text.Bounds().Dx()-ScaleN(levelWidth), text.Bounds().Dy()-ScaleN(levelHeight), text.Bounds().Dx(), text.Bounds().Dy())
		draw.Draw(text, blackRect, &image.Uniform{C: color.RGBA{A: 255}}, image.Point{}, draw.Src)

		count := imaging.Crop(img, image.Rect(
			infoButtonLocation.X+infoButtonBounds.Dx()+ScaleN(20),
			scrollLocation.Y+pxOffset+ScaleN(5),
			infoButtonLocation.X+infoButtonBounds.Dx()+ScaleN(50),
			scrollLocation.Y+pxOffset+ScaleN(45),
		))

		level := imaging.Crop(img, image.Rect(
			scrollLocation.X-ScaleN(levelWidth),
			scrollLocation.Y+pxOffset+(int(ScaleNf(textHeight))-ScaleN(levelHeight)),
			scrollLocation.X-ScaleN(scrollHorizontalOffset),
			scrollLocation.Y+pxOffset+int(ScaleNf(textHeight)),
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
	scrollLocation, infoButtonLocation, err := listingTrackers(img)
	if err != nil {
		return false, err
	}

	nextOffset := int(5*listingHeight - ScaleNf(scrollVerticalOffset))

	infoButtonBounds := ScaleBounds(data.InfoButton)
	nextCount := imaging.Crop(img, image.Rect(
		infoButtonLocation.X+infoButtonBounds.Dx(),
		scrollLocation.Y+nextOffset-ScaleN(20),
		infoButtonLocation.X+infoButtonBounds.Dx()+ScaleN(40),
		scrollLocation.Y+nextOffset+ScaleN(15),
	))

	_, cornerVal, err := Find(nextCount, Scale(data.CountCorner))
	if err != nil {
		return false, errors.Wrap(err, "failed to find count corner")
	}

	return cornerVal >= 0.95, nil
}
