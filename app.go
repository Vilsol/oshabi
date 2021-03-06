package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"

	"github.com/vilsol/oshabi/hooks"

	"github.com/go-vgo/robotgo/clipboard"
	"github.com/kbinani/screenshot"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/vilsol/oshabi/app"
	"github.com/vilsol/oshabi/config"
	"github.com/vilsol/oshabi/cv"
	"github.com/vilsol/oshabi/data"
	"github.com/vilsol/oshabi/pricing"
	"github.com/vilsol/oshabi/types"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	var scalingEventF cv.ScalingFunc = func(progress int, total int) {
		runtime.EventsEmit(ctx, "calibration", fmt.Sprintf("%d/%d", progress, total))
	}

	a.ctx = context.WithValue(ctx, cv.ScalingKey{}, scalingEventF)

	if err := data.InitData(); err != nil {
		panic(err)
	}

	if err := config.InitConfig(); err != nil {
		panic(err)
	}

	if err := cv.InitOCR(); err != nil {
		panic(err)
	}

	data.InitCrafts()
	hooks.InitializeHooks(ctx)

	// Pricing Loop
	go func() {
		if err := pricing.UpdatePricing(ctx); err != nil {
			runtime.EventsEmit(ctx, "error", err.Error())
			return
		}

		time.Sleep(time.Hour)
	}()

	_, format, err := mp3.Decode(io.NopCloser(bytes.NewReader(data.AlertMP3)))
	if err != nil {
		panic(errors.Wrap(err, "failed to load alert sound"))
	}

	if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10)); err != nil {
		panic(errors.Wrap(err, "failed to initialize speaker").Error())
	}

	runtime.EventsOn(ctx, "listings_read", func(_ ...interface{}) {
		alertSound, _, err := mp3.Decode(io.NopCloser(bytes.NewReader(data.AlertMP3)))
		if err != nil {
			runtime.EventsEmit(ctx, "error", errors.Wrap(err, "failed to load alert sound").Error())
			return
		}

		volume := &effects.Volume{
			Streamer: alertSound,
			Base:     10,
			Volume:   -1,
		}

		speaker.Play(volume)
	})
}

func (a *App) domReady(_ context.Context) {
}

func (a *App) shutdown(_ context.Context) {
}

func (a *App) Calibrate() error {
	img, err := app.CaptureScreen()
	if err != nil {
		return err
	}

	scale, err := cv.CalculateScaling(a.ctx, img)
	if err != nil {
		return err
	}

	return config.SetScaling(a.ctx, scale)
}

func (a *App) Read() error {
	listings, err := app.ReadFull(a.ctx)
	if err != nil {
		return err
	}

	return config.AddListings(a.ctx, listings)
}

func (a *App) Clear() error {
	return config.ClearListings(a.ctx)
}

func (a *App) UpdatePricing() error {
	return pricing.UpdatePricing(a.ctx)
}

func (a *App) GetListings() []types.ParsedListing {
	converted := make([]types.ParsedListing, 0)
	for harvestType, listing := range config.Get().Listings {
		for level, count := range listing {
			converted = append(converted, types.ParsedListing{
				Type:  types.HarvestType(harvestType),
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
	Display      int               `json:"display"`
	Scaling      float64           `json:"scaling"`
	Shortcut     []string          `json:"shortcut"`
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
		} else if craft.Short != nil {
			if short, ok := craft.Short[config.Get().Language]; ok {
				messages[strType] = short
			}
		}

		if _, ok := messages[strType]; !ok {
			messages[strType] = craft.Translations[config.Get().Language]
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
		Display:      config.Get().Display,
		Scaling:      config.Get().Scaling,
		Shortcut:     config.Get().Shortcut,
	}
}

func (a *App) SetListingCount(listing string, level int, count int) error {
	return config.SetListing(a.ctx, listing, level, count)
}

func (a *App) SetPrice(listing string, price string) error {
	return config.SetPrice(a.ctx, listing, price)
}

func (a *App) Copy() error {
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

	message += " - Oshabi"

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
		name := data.GetCraft(types.HarvestType(listingType)).Message

		if msg, ok := configMessages[listingType]; ok {
			name = msg
		}

		price := "1ex"
		if p, ok := configPrices[listingType]; ok {
			price = p
		} else if p, ok := pricing.Get()[types.HarvestType(listingType)]; ok {
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
		return errors.Wrap(err, "failed writing clipboard")
	}

	return nil
}

func (a *App) SetLeague(league string) error {
	if err := config.SetLeague(a.ctx, league); err != nil {
		return err
	}
	return pricing.UpdatePricing(a.ctx)
}

func (a *App) SetName(name string) error {
	return config.SetName(a.ctx, name)
}

func (a *App) SetStream(stream bool) error {
	return config.SetStream(a.ctx, stream)
}

func (a *App) SetDisplay(display int) error {
	return config.SetDisplay(a.ctx, display)
}

func (a *App) GetDisplayCount() int {
	return screenshot.NumActiveDisplays()
}

func (a *App) GetLanguages() map[string]string {
	return map[string]string{
		string(config.LanguageEnglish):    "English",
		string(config.LanguagePortuguese): "Portugu??s",
		string(config.LanguageRussian):    "??????????????",
		string(config.LanguageThai):       "?????????????????????",
		string(config.LanguageFrench):     "Fran??ais",
		string(config.LanguageGerman):     "Deutsch",
		string(config.LanguageSpanish):    "Espa??ol",
		string(config.LanguageChinese):    "????????????",
		string(config.LanguageKorean):     "?????????",
		string(config.LanguageJapanese):   "?????????",
		string(config.LanguageTaiwanese):  "????????????",
	}
}

func (a *App) SetLanguage(language string) error {
	return config.SetLanguage(a.ctx, language)
}

func (a *App) SetShortcut(shortcut []string) error {
	err := config.SetShortcut(a.ctx, shortcut)
	if err != nil {
		return err
	}
	hooks.UpdateShortcut(a.ctx)
	return nil
}
