package main

import (
	"fmt"
	"image/png"
	"os"
	"testing"

	"github.com/vilsol/oshabi/app"
	"github.com/vilsol/oshabi/config"
	"github.com/vilsol/oshabi/cv"
	"github.com/vilsol/oshabi/data"
)

func init() {
	if err := config.InitConfig(); err != nil {
		panic(err)
	}

	if err := data.InitData(); err != nil {
		panic(err)
	}

	if err := cv.InitOCR(); err != nil {
		panic(err)
	}
}

func TestCalculateScaling(t *testing.T) {
	f, err := os.Open("testdata/diff_scaling.png")
	if err != nil {
		t.Error(err)
		return
	}

	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		t.Error(err)
		return
	}

	scaling, err := cv.CalculateScaling(img)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("Calculated scaling:", scaling)
	config.Cfg.Scaling = scaling

	listings, err := app.ReadImage(img, 0, 5)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("%#v\n", listings)
}
