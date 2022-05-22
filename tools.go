//go:build tools

package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/vilsol/oshabi/data"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

//go:generate go run tools.go

func main() {
	downloadHarvestMods()
}

// Download all harvest mods and save them to JSON
func downloadHarvestMods() error {
	existingCrafts := make(map[string]data.HarvestCraft)

	existingFile, err := os.ReadFile("data/crafts.json")
	if err == nil {
		json.Unmarshal(existingFile, &existingCrafts)
	}

	res, err := http.Get("https://poedb.tw/us/Horticrafting")
	if err != nil {
		return errors.Wrap(err, "failed fetching PoeDB Horticrafting page")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.Wrap(err, "PoeDB returned a non-200 response")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return errors.Wrap(err, "failed parsing page as html")
	}

	doc.Find("#HarvestHarvestSeeds tbody tr").Each(func(_ int, s *goquery.Selection) {
		s.Find("li").Each(func(_ int, craft *goquery.Selection) {
			if _, ok := existingCrafts[craft.Text()]; !ok {
				existingCrafts[craft.Text()] = data.HarvestCraft{
					Type:    "",
					Message: "",
				}
			}
		})
	})

	jsonBytes, err := json.MarshalIndent(existingCrafts, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed converting crafts to json")
	}

	if err := os.WriteFile("data/crafts.json", jsonBytes, 0777); err != nil {
		return errors.Wrap(err, "failed writing crafts.json")
	}
}
