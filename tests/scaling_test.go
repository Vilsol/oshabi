package tests

import (
	"context"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/MarvinJWendt/testza"
	"github.com/pkg/errors"

	"github.com/vilsol/oshabi/config"
	"github.com/vilsol/oshabi/cv"
	"github.com/vilsol/oshabi/data"
	"github.com/vilsol/oshabi/types"
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

	config.Cfg.Language = config.LanguageEnglish
}

func TestScalingHorticrafting(t *testing.T) {
	// Disabled until translations are fixed
	//runScalingTest(t, "../testdata/horticrafting/chi_tra", false, config.LanguageTaiwanese, []types.ParsedListing{
	//	{Type: data.HarvestReforgeCriticalMoreCommon, Count: 1, Level: 83},
	//	{Type: data.HarvestReforgeCriticalMoreCommon, Count: 1, Level: 83},
	//	{Type: data.HarvestReforgeCriticalMoreCommon, Count: 1, Level: 83},
	//	{Type: data.HarvestChangeResistColdToFire, Count: 1, Level: 83},
	//	{Type: data.HarvestReforgePhysical, Count: 1, Level: 83},
	//})

	runScalingTest(t, "../testdata/horticrafting/deu", true, config.LanguageGerman, []types.ParsedListing{
		{Type: data.HarvestRandomiseInfluenceArmour, Count: 1, Level: 83},
		{Type: data.HarvestChangeUniqueHarbinger, Count: 1, Level: 83},
		{Type: data.HarvestChangeGemToGem, Count: 1, Level: 77},
		{Type: data.HarvestEnchantArmourLife, Count: 1, Level: 77},
		{Type: data.HarvestEnchantArmourLife, Count: 1, Level: 83},
	})

	runScalingTest(t, "../testdata/horticrafting/eng", true, config.LanguageEnglish, []types.ParsedListing{
		{Type: data.HarvestFiveSockets, Count: 1, Level: 77},
		{Type: data.HarvestChangeGemToGem, Count: 1, Level: 77},
		{Type: data.HarvestEnchantArmourLife, Count: 1, Level: 77},
		{Type: data.HarvestQualityFlask, Count: 1, Level: 77},
		{Type: data.HarvestEnchantWeaponCritical, Count: 1, Level: 77},
	})

	runScalingTest(t, "../testdata/horticrafting/fra", true, config.LanguageFrench, []types.ParsedListing{
		{Type: data.HarvestRandomiseInfluenceArmour, Count: 1, Level: 83},
		{Type: data.HarvestChangeUniqueHarbinger, Count: 1, Level: 83},
		{Type: data.HarvestChangeGemToGem, Count: 1, Level: 77},
		{Type: data.HarvestEnchantArmourLife, Count: 1, Level: 77},
		{Type: data.HarvestEnchantArmourLife, Count: 1, Level: 83},
	})

	runScalingTest(t, "../testdata/horticrafting/kor", true, config.LanguageKorean, []types.ParsedListing{
		{Type: data.HarvestRandomiseInfluenceArmour, Count: 1, Level: 83},
		{Type: data.HarvestChangeUniqueHarbinger, Count: 1, Level: 83},
		{Type: data.HarvestChangeGemToGem, Count: 1, Level: 77},
		{Type: data.HarvestEnchantArmourLife, Count: 1, Level: 77},
		{Type: data.HarvestEnchantArmourLife, Count: 1, Level: 83},
	})

	runScalingTest(t, "../testdata/horticrafting/por", true, config.LanguagePortuguese, []types.ParsedListing{
		{Type: data.HarvestRandomiseInfluenceArmour, Count: 1, Level: 83},
		{Type: data.HarvestChangeUniqueHarbinger, Count: 1, Level: 83},
		{Type: data.HarvestChangeGemToGem, Count: 1, Level: 77},
		{Type: data.HarvestEnchantArmourLife, Count: 1, Level: 77},
		{Type: data.HarvestEnchantArmourLife, Count: 1, Level: 83},
	})

	runScalingTest(t, "../testdata/horticrafting/rus", true, config.LanguageRussian, []types.ParsedListing{
		{Type: data.HarvestRandomiseInfluenceArmour, Count: 1, Level: 83},
		{Type: data.HarvestChangeUniqueHarbinger, Count: 1, Level: 83},
		{Type: data.HarvestChangeGemToGem, Count: 1, Level: 77},
		{Type: data.HarvestEnchantArmourLife, Count: 1, Level: 77},
		{Type: data.HarvestEnchantArmourLife, Count: 1, Level: 83},
	})

	runScalingTest(t, "../testdata/horticrafting/spa", true, config.LanguageSpanish, []types.ParsedListing{
		{Type: data.HarvestRandomiseInfluenceArmour, Count: 1, Level: 83},
		{Type: data.HarvestChangeUniqueHarbinger, Count: 1, Level: 83},
		{Type: data.HarvestChangeGemToGem, Count: 1, Level: 77},
		{Type: data.HarvestEnchantArmourLife, Count: 1, Level: 77},
		{Type: data.HarvestEnchantArmourLife, Count: 1, Level: 83},
	})
}

func TestScalingGrove(t *testing.T) {
	runScalingTest(t, "../testdata/grove/eng", false, config.LanguageEnglish, []types.ParsedListing{
		{Type: "SacrificeMap1Anarchy", Count: 1, Level: 72},
		{Type: data.HarvestReforgeCaster, Count: 1, Level: 72},
		{Type: data.HarvestReforgeCasterMoreCommon, Count: 1, Level: 72},
		{Type: data.HarvestReforgeFire, Count: 5, Level: 72},
		{Type: data.HarvestReforgeFireMoreCommon, Count: 1, Level: 72},
	})
}

func runScalingTest(t *testing.T, dirPath string, shouldScrollDown bool, lang config.Language, expected []types.ParsedListing) {
	t.Run(string(lang), func(t *testing.T) {
		dir, err := os.ReadDir(dirPath)
		if err != nil {
			t.Error(err)
			return
		}

		for _, entry := range dir {
			if !entry.IsDir() {
				t.Run(entry.Name(), func(t *testing.T) {
					imgPath := filepath.Join(dirPath, entry.Name())

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

					scaling, err := cv.CalculateScaling(context.Background(), img)
					if err != nil {
						t.Error(err)
						return
					}

					config.Cfg.Language = lang
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
	})
}
