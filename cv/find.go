package cv

import (
	"gocv.io/x/gocv"
	"image"
)

func Find(static image.Image, dynamic image.Image) (image.Point, float32, error) {
	matStatic, err := gocv.ImageToMatRGB(static)
	if err != nil {
		return image.Point{}, 0, err
	}

	grayStatic := gocv.NewMat()
	gocv.CvtColor(matStatic, &grayStatic, gocv.ColorRGBToGray)

	matDynamic, err := gocv.ImageToMatRGB(dynamic)
	if err != nil {
		return image.Point{}, 0, err
	}

	grayDynamic := gocv.NewMat()
	gocv.CvtColor(matDynamic, &grayDynamic, gocv.ColorRGBToGray)

	m := gocv.NewMat()
	result := gocv.NewMat()
	gocv.MatchTemplate(grayStatic, grayDynamic, &result, gocv.TmCcoeffNormed, m)

	_, maxVal, _, location := gocv.MinMaxLoc(result)
	_ = result.Close()
	_ = m.Close()
	_ = matDynamic.Close()
	_ = grayDynamic.Close()

	return location, maxVal, nil
}
