//go:build tools

package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/vilsol/oshabi/config"
	"github.com/vilsol/oshabi/data"
	"github.com/vilsol/oshabi/types"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//go:generate go run tools.go

func main() {
	if err := downloadHarvestMods(); err != nil {
		panic(err)
	}
}

var languageMapping = map[config.Language]string{
	config.LanguageEnglish:    "us",
	config.LanguagePortuguese: "pt",
	config.LanguageRussian:    "ru",
	config.LanguageThai:       "th",
	config.LanguageFrench:     "fr",
	config.LanguageGerman:     "de",
	config.LanguageSpanish:    "sp",
	config.LanguageChinese:    "cn",
	config.LanguageKorean:     "kr",
	config.LanguageJapanese:   "jp",
}

const urlPattern = "https://poedb.tw/%s/Horticrafting"

// Download all harvest mods and save them to JSON
func downloadHarvestMods() error {
	existingCrafts := make(map[types.HarvestType]data.HarvestCraft)

	existingFile, err := os.ReadFile("data/crafts.json")
	if err == nil {
		_ = json.Unmarshal(existingFile, &existingCrafts)
	}

	remapped := make(map[string]types.HarvestType)
	for _, craft := range existingCrafts {
		remapped[craft.Translations[config.LanguageEnglish]] = craft.Type
	}

	doc, err := fetchDecode(config.LanguageEnglish)
	if err != nil {
		return err
	}

	numberMapping := make(map[int]types.HarvestType)
	n := 0

	doc.Find("#HarvestHarvestSeeds tbody tr").Each(func(_ int, s *goquery.Selection) {
		s.Find("li").Each(func(_ int, craft *goquery.Selection) {
			if existingType, ok := remapped[craft.Text()]; ok {
				numberMapping[n] = existingType
			} else {
				existingCrafts[types.HarvestType(craft.Text())] = data.HarvestCraft{
					Type:    "",
					Message: "",
					Pricing: "",
					Translations: map[config.Language]string{
						config.LanguageEnglish: craft.Text(),
					},
				}
				numberMapping[n] = types.HarvestType(craft.Text())
			}
			n++
		})
	})

	for lang := range languageMapping {
		if lang != config.LanguageEnglish {
			doc, err := fetchDecode(lang)
			if err != nil {
				return err
			}

			i := 0
			doc.Find("#HarvestHarvestSeeds tbody tr").Each(func(_ int, s *goquery.Selection) {
				s.Find("li").Each(func(_ int, craft *goquery.Selection) {
					clean := craft.Text()
					// Remove zero width spaces
					clean = strings.Replace(clean, "\u200B", "", -1)
					existingCrafts[numberMapping[i]].Translations[lang] = clean
					i++
				})
			})
		}
	}

	jsonBytes, err := json.MarshalIndent(existingCrafts, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed converting crafts to json")
	}

	if err := os.WriteFile("data/crafts.json", jsonBytes, 0777); err != nil {
		return errors.Wrap(err, "failed writing crafts.json")
	}

	return nil
}

func fetchDecode(language config.Language) (*goquery.Document, error) {
	res, err := http.Get(fmt.Sprintf(urlPattern, languageMapping[language]))
	if err != nil {
		return nil, errors.Wrap(err, "failed fetching PoeDB Horticrafting page")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.Wrap(err, "PoeDB returned a non-200 response")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed parsing page as html")
	}

	return doc, nil
}
