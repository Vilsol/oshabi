package cv

import (
	"fmt"
	"github.com/vilsol/oshabi/data"
	"image/png"
	"os"
	"testing"
)

func init() {
	if err := data.InitData(); err != nil {
		panic(err)
	}
}

func TestCalculateScaling(t *testing.T) {
	f, err := os.Open("../testdata/diff_scaling.png")
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

	scaling, err := CalculateScaling(img)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(scaling)
}
