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

func listingTrackers(img image.Image) (image.Point, image.Point, error) {
	scrollLocation, _, err := Find(img, data.ScrollUp)
	if err != nil {
		return image.Point{}, image.Point{}, errors.Wrap(err, "failed to find scroll up button")
	}

	infoButtonLocation, _, err := Find(img, data.InfoButton)
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
		pxOffset := (i+offset)*173 - 10

		text := imaging.Crop(img, image.Rect(
			infoButtonLocation.X+data.InfoButton.Bounds().Dx()+190,
			scrollLocation.Y+pxOffset,
			scrollLocation.X-35,
			scrollLocation.Y+pxOffset+155,
		))

		blackRect := image.Rect(text.Bounds().Dx()-145, text.Bounds().Dy()-35, text.Bounds().Dx(), text.Bounds().Dy())
		draw.Draw(text, blackRect, &image.Uniform{C: color.RGBA{A: 255}}, image.Point{}, draw.Src)

		count := imaging.Crop(img, image.Rect(
			infoButtonLocation.X+data.InfoButton.Bounds().Dx()+20,
			scrollLocation.Y+pxOffset+5,
			infoButtonLocation.X+data.InfoButton.Bounds().Dx()+50,
			scrollLocation.Y+pxOffset+40,
		))

		level := imaging.Crop(img, image.Rect(
			scrollLocation.X-180,
			scrollLocation.Y+pxOffset+120,
			scrollLocation.X-35,
			scrollLocation.Y+pxOffset+155,
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

	nextOffset := 5*173 - 10

	nextCount := imaging.Crop(img, image.Rect(
		infoButtonLocation.X+data.InfoButton.Bounds().Dx(),
		scrollLocation.Y+nextOffset-20,
		infoButtonLocation.X+data.InfoButton.Bounds().Dx()+40,
		scrollLocation.Y+nextOffset+15,
	))

	_, cornerVal, err := Find(nextCount, data.CountCorner)
	if err != nil {
		return false, errors.Wrap(err, "failed to find count corner")
	}

	return cornerVal >= 0.95, nil
}
