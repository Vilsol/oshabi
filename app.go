package main

import (
	"context"
	"fmt"
	"github.com/go-vgo/robotgo/clipboard"
	"github.com/olekukonko/tablewriter"
	"github.com/vilsol/oshabi/app"
	"github.com/vilsol/oshabi/config"
	"github.com/vilsol/oshabi/cv"
	"github.com/vilsol/oshabi/data"
	"github.com/vilsol/oshabi/pricing"
	"github.com/vilsol/oshabi/types"
	"strconv"
	"strings"
	"time"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	config.InitConfig()
	cv.InitOCR()
	data.InitCrafts()
	app.InitializeApp(ctx)

	// Pricing Loop
	go func() {
		pricing.UpdatePricing(a.ctx)
		time.Sleep(time.Hour)
	}()
}

func (a App) domReady(ctx context.Context) {
}

func (a *App) shutdown(ctx context.Context) {
}

func (a *App) Calibrate() {
	img := app.CaptureScreen()

	scale, err := cv.CalculateScaling(img)
	if err != nil {
		fmt.Println("Error calibrating", err)
		return
	}

	fmt.Println("Final scale:", scale)
	// TODO actually use scaling
}

func (a *App) Read() bool {
	config.AddListings(a.ctx, app.ReadFull(nil))
	return true
}

func (a *App) Clear() bool {
	config.ClearListings(a.ctx)
	return true
}

func (a *App) UpdatePricing() {
	pricing.UpdatePricing(a.ctx)
}

func (a *App) GetListings() []types.ParsedListing {
	converted := make([]types.ParsedListing, 0)
	for harvestType, listing := range config.Get().Listings {
		for level, count := range listing {
			converted = append(converted, types.ParsedListing{
				Type:  data.HarvestType(harvestType),
				Count: count,
				Level: level,
			})
		}
	}
	return converted
}

type ConvertedConfig struct {
	MarketPrices map[string]string `json:"market_prices"`
	Prices       map[string]string `json:"prices"`
	Language     config.Language   `json:"language"`
	League       config.League     `json:"league"`
	Messages     map[string]string `json:"messages"`
	Name         string            `json:"name"`
	Stream       bool              `json:"stream"`
}

func (a *App) GetConfig() ConvertedConfig {
	configMessages := config.Get().Messages
	configPrices := config.Get().Prices
	rawMarket := pricing.Get()

	prices := make(map[string]string)
	messages := make(map[string]string)
	for _, craft := range data.AllCrafts() {
		strType := string(craft.Type)
		if msg, ok := configMessages[strType]; ok {
			messages[strType] = msg
		} else {
			messages[strType] = craft.Message
		}

		if p, ok := configPrices[strType]; ok {
			prices[strType] = p
		} else if p, ok := rawMarket[craft.Type]; ok {
			if p.Exalt > 1 {
				prices[strType] = fmt.Sprintf("%fex", p.Exalt)
			} else {
				prices[strType] = fmt.Sprintf("%dc", p.Chaos)
			}
		} else {
			prices[strType] = "1ex"
		}
	}

	marketPrices := make(map[string]string)
	for harvestType, price := range rawMarket {
		if price.Exalt > 1 {
			marketPrices[string(harvestType)] = fmt.Sprintf("%fex", price.Exalt)
		} else {
			marketPrices[string(harvestType)] = fmt.Sprintf("%dc", price.Chaos)
		}
	}

	return ConvertedConfig{
		MarketPrices: marketPrices,
		Prices:       prices,
		Language:     config.Get().Language,
		League:       config.Get().League,
		Messages:     messages,
		Name:         config.Get().Name,
		Stream:       config.Get().Stream,
	}
}

func (a *App) SetListingCount(listing string, level int, count int) {
	config.SetListing(a.ctx, listing, level, count)
}

func (a *App) SetPrice(listing string, price string) {
	config.SetPrice(a.ctx, listing, price)
}

func (a *App) Copy() {
	message := "**WTS "

	switch config.Get().League {
	case config.LeagueStandard:
		message += "Standard"
	case config.LeagueSoftcore:
		message += "Softcore"
	case config.LeagueHardcore:
		message += "Hardcore"
	}

	message += "**"

	if config.Get().Name != "" {
		message += " - IGN: **" + config.Get().Name + "**"
	}

	// TODO Add mark
	// message += " - Oshabi"

	message += "\n"

	if config.Get().Stream {
		message += "*Can stream if requested*\n"
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetBorder(false)
	table.SetAutoWrapText(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetHeaderLine(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetNoWhiteSpace(true)
	table.SetTablePadding(" ")

	configMessages := config.Get().Messages
	configPrices := config.Get().Prices
	for listingType, listing := range config.Get().Listings {
		name := data.GetCraftByType(data.HarvestType(listingType)).Craft.Message

		if msg, ok := configMessages[listingType]; ok {
			name = msg
		}

		price := "1ex"
		if p, ok := configPrices[listingType]; ok {
			price = p
		} else if p, ok := pricing.Get()[data.HarvestType(listingType)]; ok {
			if p.Exalt > 1 {
				price = fmt.Sprintf("%fex", p.Exalt)
			} else {
				price = fmt.Sprintf("%dc", p.Chaos)
			}
		}

		for level, count := range listing {
			table.Append([]string{
				"`" + strconv.Itoa(count) + "x",
				name,
				"[" + strconv.Itoa(level) + "]",
				"< " + price + " >`",
			})
		}
	}

	table.Render()
	message += tableString.String()

	message = strings.TrimSpace(message)

	if err := clipboard.WriteAll(message); err != nil {
		panic(err)
	}
}

func (a *App) SetLeague(league string) {
	config.SetLeague(a.ctx, league)
	pricing.UpdatePricing(a.ctx)
}

func (a *App) SetName(name string) {
	config.SetName(a.ctx, name)
}

func (a *App) SetStream(stream bool) {
	config.SetStream(a.ctx, stream)
}
