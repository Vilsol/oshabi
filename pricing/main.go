package pricing

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/vilsol/oshabi/config"
	"github.com/vilsol/oshabi/data"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"net/http"
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

func UpdatePricing(ctx context.Context) {
	league := config.Get().League

	response, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/The-Forbidden-Trove/tft-data-prices/master/%s/harvest.json", league))
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var raw RawData
	if err := json.Unmarshal(body, &raw); err != nil {
		panic(err)
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
			fmt.Println("Did not find craft for price:", i.Name)
		}
	}

	pricing = resultPrices

	runtime.EventsEmit(ctx, "config_updated")
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
