package cv

import (
	"image"
	"math"

	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"github.com/vilsol/oshabi/config"
	"github.com/vilsol/oshabi/data"
	"gocv.io/x/gocv"
)

func CalculateScaling(img image.Image) (float64, error) {
	bottomLeftCrop := imaging.Crop(img, image.Rect(0, img.Bounds().Dy()/2, img.Bounds().Dx()/2, img.Bounds().Dy()))
	topRightCrop := imaging.Crop(img, image.Rect(img.Bounds().Dx()/2, 0, img.Bounds().Dx(), img.Bounds().Dy()/2))

	menuButtonScaling, menuButtonValue, _, err := ScaleAndFind(bottomLeftCrop, data.MenuButton)
	if err != nil {
		return 0, errors.New("failed finding menu button")
	}

	if menuButtonValue < 0.8 {
		return 0, errors.New("could not find menu button on bottom left of screen")
	}

	topRightXScaling, topRightXValue, _, err := ScaleAndFind(topRightCrop, data.TopRightX)
	if err != nil {
		return 0, errors.New("failed finding top right x")
	}

	if topRightXValue < 0.8 {
		return 0, errors.New("could not find x button on top right of screen")
	}

	if math.Max(menuButtonScaling, topRightXScaling)-math.Min(menuButtonScaling, topRightXScaling) > 0.1 {
		return 0, errors.New("cv differentiates >0.1, please post a bug report with a screenshot of an opened horticrafting station")
	}

	avg := (menuButtonScaling + topRightXScaling) / 2

	// Round to the nearest 0.005
	return math.Round(avg/0.005) * 0.005, nil
}

type resultTuple struct {
	Result gocv.Mat
	Scale  float64
}

func ScaleAndFind(static image.Image, dynamic image.Image) (float64, float32, image.Point, error) {
	matLeftCrop, err := gocv.ImageToMatRGB(static)
	if err != nil {
		return 0, 0, image.Point{}, errors.Wrap(err, "failed converting image to mat")
	}

	grayLeftCrop := gocv.NewMat()
	gocv.CvtColor(matLeftCrop, &grayLeftCrop, gocv.ColorRGBToGray)

	size := dynamic.Bounds()

	results := make(chan resultTuple)
	count := 0

	for i := 0.25; i < 2.01; i += 0.005 {
		count++
		scale := math.Round(i*1000) / 1000
		width := int(float64(size.Dx()) * scale)
		height := int(float64(size.Dy()) * scale)

		if width > static.Bounds().Dx() || height > static.Bounds().Dy() {
			continue
		}

		resized := imaging.Resize(dynamic, width, height, imaging.NearestNeighbor)

		matResized, err := gocv.ImageToMatRGB(resized)
		if err != nil {
			return 0, 0, image.Point{}, errors.Wrap(err, "failed converting image to mat")
		}

		grayResized := gocv.NewMat()
		gocv.CvtColor(matResized, &grayResized, gocv.ColorRGBToGray)

		go func(grayLeftCrop gocv.Mat, grayResized gocv.Mat, scale float64) {
			m := gocv.NewMat()
			result := gocv.NewMat()
			gocv.MatchTemplate(grayLeftCrop, grayResized, &result, gocv.TmCcoeffNormed, m)
			results <- resultTuple{
				Result: result,
				Scale:  scale,
			}
			_ = m.Close()
			_ = matResized.Close()
			_ = grayResized.Close()
		}(grayLeftCrop, grayResized, scale)
	}

	bestScaling := 0.05
	bestValue := float32(0)
	var bestLocation image.Point

	for i := 0; i < count; i++ {
		result := <-results

		_, maxVal, _, maxLoc := gocv.MinMaxLoc(result.Result)

		if maxVal > bestValue {
			bestValue = maxVal
			bestScaling = result.Scale
			bestLocation = maxLoc
		}

		_ = result.Result.Close()
	}

	close(results)

	return bestScaling, bestValue, bestLocation, nil
}

func Scale(img image.Image) image.Image {
	size := img.Bounds()
	scale := math.Round(config.Get().Scaling*1000) / 1000
	width := int(float64(size.Dx()) * scale)
	height := int(float64(size.Dy()) * scale)
	return imaging.Resize(img, width, height, imaging.NearestNeighbor)
}

func ScaleBounds(img image.Image) image.Rectangle {
	size := img.Bounds()
	scale := math.Round(config.Get().Scaling*1000) / 1000
	width := int(float64(size.Dx()) * scale)
	height := int(float64(size.Dy()) * scale)
	return image.Rect(0, 0, width, height)
}

func ScaleN(n int) int {
	return int(ScaleNf(float64(n)))
}

func ScaleNf(n float64) float64 {
	scale := math.Round(config.Get().Scaling*1000) / 1000
	return math.Round(scale * n)
}
