package tests

import (
	"image/png"
	"os"
	"path"
	"testing"

	"github.com/pkg/errors"

	"github.com/MarvinJWendt/testza"
	"github.com/vilsol/oshabi/types"

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

func TestScalingHorticrafting(t *testing.T) {
	runScalingTest(t, "../testdata/horticrafting", true, []types.ParsedListing{
		{Type: "FiveSockets", Count: 1, Level: 77},
		{Type: "ChangeGemToGem", Count: 1, Level: 77},
		{Type: "EnchantArmourLife", Count: 1, Level: 77},
		{Type: "QualityFlask", Count: 1, Level: 77},
		{Type: "EnchantWeaponCritical", Count: 1, Level: 77},
	})
}

func TestScalingGrove(t *testing.T) {
	runScalingTest(t, "../testdata/grove", false, []types.ParsedListing{
		{Type: "SacrificeMap1Anarchy", Count: 1, Level: 72},
		{Type: "ReforgeCaster", Count: 1, Level: 72},
		{Type: "ReforgeCasterMoreCommon", Count: 1, Level: 72},
		{Type: "ReforgeFire", Count: 5, Level: 72},
		{Type: "ReforgeFireMoreCommon", Count: 1, Level: 72},
	})
}

func runScalingTest(t *testing.T, dirPath string, shouldScrollDown bool, expected []types.ParsedListing) {
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		t.Error(err)
		return
	}

	for _, entry := range dir {
		if !entry.IsDir() {
			t.Run(entry.Name(), func(t *testing.T) {
				imgPath := path.Join(dirPath, entry.Name())

				f, err := os.Open(imgPath)
				if err != nil {
					t.Error(err)
					return
				}

				img, err := png.Decode(f)
				if err != nil {
					t.Error(err)
					return
				}

				if err := f.Close(); err != nil {
					t.Error(err)
					return
				}

				scaling, err := cv.CalculateScaling(img)
				if err != nil {
					t.Error(err)
					return
				}

				config.Cfg.Scaling = scaling

				listings, err := cv.ReadImage(img, 0, 5)
				if err != nil {
					t.Error(err)
					return
				}

				testza.AssertEqual(t, expected, listings)

				infoButtonLocation, infoButtonValue, err := cv.Find(img, cv.Scale(data.InfoButton))
				if err != nil {
					t.Error(err)
					return
				}

				if infoButtonValue < 0.7 {
					t.Error(errors.New("info button not found"))
					return
				}

				canScrollDown, err := cv.CanScrollDown(infoButtonLocation, false, img)
				if err != nil {
					t.Error(err)
					return
				}

				testza.AssertEqual(t, shouldScrollDown, canScrollDown)
			})
		}
	}
}
