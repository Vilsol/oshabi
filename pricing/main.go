package pricing

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/vilsol/oshabi/config"
	"github.com/vilsol/oshabi/data"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type CraftPricing struct {
	Exalt         float64 `json:"exalt"`
	Chaos         int     `json:"chaos"`
	LowConfidence bool    `json:"lowConfidence"`
}

var pricing = make(map[data.HarvestType]CraftPricing)

type RawData struct {
	Timestamp int64 `json:"timestamp"`
	Data      []struct {
		Name          string  `json:"name"`
		Exalt         float64 `json:"exalt"`
		Chaos         int     `json:"chaos"`
		LowConfidence bool    `json:"lowConfidence"`
	} `json:"data"`
}

func UpdatePricing(ctx context.Context) error {
	league := config.Get().League

	response, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/The-Forbidden-Trove/tft-data-prices/master/%s/harvest.json", league))
	if err != nil {
		return errors.Wrap(err, "failed fetching harvest prices")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "failed reading harvest price body")
	}

	var raw RawData
	if err := json.Unmarshal(body, &raw); err != nil {
		return errors.Wrap(err, "failed parsing harvest price body to json")
	}

	resultPrices := make(map[data.HarvestType]CraftPricing)
	for _, i := range raw.Data {
		craft := data.GetCraftByPricing(i.Name)
		if craft != nil {
			resultPrices[craft.Craft.Type] = CraftPricing{
				Exalt:         i.Exalt,
				Chaos:         i.Chaos,
				LowConfidence: i.LowConfidence,
			}
		} else {
			runtime.EventsEmit(ctx, "warning", "Did not find craft for price name:", i.Name)
		}
	}

	pricing = resultPrices

	runtime.EventsEmit(ctx, "config_updated")

	return nil
}

func GetPrice(craft data.HarvestType) *CraftPricing {
	if price, ok := pricing[craft]; ok {
		return &price
	}
	return nil
}

func Get() map[data.HarvestType]CraftPricing {
	return pricing
}
