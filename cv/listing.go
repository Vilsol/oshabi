package cv

import (
	"bytes"
	"github.com/disintegration/imaging"
	"github.com/otiai10/gosseract/v2"
	"github.com/vilsol/oshabi/data"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

var client *gosseract.Client

func InitOCR() {
	client = gosseract.NewClient()
	if err := client.SetLanguage("eng"); err != nil {
		panic(err)
	}
	println("Tesseract Version", client.Version())
}

func OCRListing(img image.Image) string {
	buff := new(bytes.Buffer)
	if err := png.Encode(buff, PrepareForOCR(img)); err != nil {
		panic(err)
	}

	if err := client.SetWhitelist("1234567890abcdefghijklmnopqrstuvwxyz,.%' "); err != nil {
		panic(err)
	}

	if err := client.SetImageFromBytes(buff.Bytes()); err != nil {
		panic(err)
	}

	text, err := client.Text()
	if err != nil {
		panic(err)
	}

	return text
}

func OCRListingCount(img image.Image) string {
	buff := new(bytes.Buffer)
	if err := png.Encode(buff, PrepareForOCR(img)); err != nil {
		panic(err)
	}

	if err := client.SetWhitelist("1234567890"); err != nil {
		panic(err)
	}

	if err := client.SetImageFromBytes(buff.Bytes()); err != nil {
		panic(err)
	}

	text, err := client.Text()
	if err != nil {
		panic(err)
	}

	return text
}

func OCRListingLevel(img image.Image) string {
	buff := new(bytes.Buffer)
	if err := png.Encode(buff, PrepareForOCR(img)); err != nil {
		panic(err)
	}

	if err := client.SetWhitelist("levl1234567890 "); err != nil {
		panic(err)
	}

	if err := client.SetImageFromBytes(buff.Bytes()); err != nil {
		panic(err)
	}

	text, err := client.Text()
	if err != nil {
		panic(err)
	}

	return text
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

func listingTrackers(img image.Image) (image.Image, image.Point, image.Image, image.Point) {
	scrollUp, _, err := image.Decode(bytes.NewReader(data.ScrollUpPNG))
	if err != nil {
		panic(err)
	}

	scrollLocation, _, err := Find(img, scrollUp)
	if err != nil {
		panic(err)
	}

	infoButton, _, err := image.Decode(bytes.NewReader(data.InfoButtonPNG))
	if err != nil {
		panic(err)
	}

	infoButtonLocation, _, err := Find(img, infoButton)
	if err != nil {
		panic(err)
	}

	return scrollUp, scrollLocation, infoButton, infoButtonLocation
}

func ExtractToListings(img image.Image, offset int, limit int) []RawListing {
	_, scrollLocation, infoButton, infoButtonLocation := listingTrackers(img)

	listings := make([]RawListing, limit)

	for i := 0; i < limit; i++ {
		pxOffset := (i+offset)*173 - 10

		text := imaging.Crop(img, image.Rect(
			infoButtonLocation.X+infoButton.Bounds().Dx()+190,
			scrollLocation.Y+pxOffset,
			scrollLocation.X-35,
			scrollLocation.Y+pxOffset+155,
		))

		blackRect := image.Rect(text.Bounds().Dx()-145, text.Bounds().Dy()-35, text.Bounds().Dx(), text.Bounds().Dy())
		draw.Draw(text, blackRect, &image.Uniform{C: color.RGBA{A: 255}}, image.Point{}, draw.Src)

		count := imaging.Crop(img, image.Rect(
			infoButtonLocation.X+infoButton.Bounds().Dx()+20,
			scrollLocation.Y+pxOffset+5,
			infoButtonLocation.X+infoButton.Bounds().Dx()+50,
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

	return listings
}

func CanScrollDown(img image.Image) bool {
	_, scrollLocation, infoButton, infoButtonLocation := listingTrackers(img)

	nextOffset := 5*173 - 10

	nextCount := imaging.Crop(img, image.Rect(
		infoButtonLocation.X+infoButton.Bounds().Dx(),
		scrollLocation.Y+nextOffset-20,
		infoButtonLocation.X+infoButton.Bounds().Dx()+40,
		scrollLocation.Y+nextOffset+15,
	))

	countCorner, _, err := image.Decode(bytes.NewReader(data.CountCornerPNG))
	if err != nil {
		panic(err)
	}

	_, cornerVal, err := Find(nextCount, countCorner)
	if err != nil {
		panic(err)
	}

	return cornerVal >= 0.95
}
