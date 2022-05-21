package config

import (
	"context"
	"encoding/json"
	"github.com/vilsol/oshabi/types"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"path"
)

type Language string

const (
	LanguageEnglish    = Language("eng")
	LanguagePortuguese = Language("por")
	LanguageRussian    = Language("rus")
	LanguageThai       = Language("tha")
	LanguageFrench     = Language("fra")
	LanguageGerman     = Language("deu")
	LanguageSpanish    = Language("spa")
)

type League string

const (
	LeagueStandard = League("std")
	LeagueSoftcore = League("lsc")
	LeagueHardcore = League("lhc")
)

type Config struct {
	Version  int                    `json:"version"`
	Scaling  float64                `json:"scaling"`
	Prices   map[string]string      `json:"prices"`
	Listings map[string]map[int]int `json:"listings"`
	Language Language               `json:"language"`
	League   League                 `json:"league"`
	Messages map[string]string      `json:"messages"`
	Name     string                 `json:"name"`
	Stream   bool                   `json:"stream"`
}

var config Config

func InitConfig() {
	Load()

	if config.Version == 0 {
		config = Config{
			Version:  1,
			Scaling:  1,
			Prices:   map[string]string{},
			Listings: map[string]map[int]int{},
			Language: LanguageEnglish,
			League:   LeagueSoftcore,
			Messages: map[string]string{},
		}

		Save()
	}
}

func Load() {
	dir := GetConfigDir()
	file, err := os.ReadFile(path.Join(dir, "config.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}
	if err := json.Unmarshal(file, &config); err != nil {
		panic(err)
	}
}

func Save() {
	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		panic(err)
	}

	dir := GetConfigDir()
	err = os.WriteFile(path.Join(dir, "config.json"), b, 0777)
	if err != nil {
		panic(err)
	}
}

func GetConfigDir() string {
	configDir, err := os.UserConfigDir()

	if err != nil {
		panic(err)
	}

	finalDir := path.Join(configDir, "oshabi")

	if err := os.MkdirAll(finalDir, 0777); err != nil {
		if !os.IsExist(err) {
			panic(err)
		}
	}

	return finalDir
}

func AddListings(ctx context.Context, listings []types.ParsedListing) {
	for _, listing := range listings {
		strType := string(listing.Type)
		if _, ok := config.Listings[strType]; !ok {
			config.Listings[strType] = make(map[int]int)
		}

		if _, ok := config.Listings[strType][listing.Level]; !ok {
			config.Listings[strType][listing.Level] = 0
		}

		config.Listings[strType][listing.Level] += listing.Count
	}

	Save()

	runtime.EventsEmit(ctx, "listings_updated")
}

func ClearListings(ctx context.Context) {
	config.Listings = make(map[string]map[int]int)
	Save()
	runtime.EventsEmit(ctx, "listings_updated")
}

func Get() Config {
	return config
}

func SetListing(ctx context.Context, listing string, level int, count int) {
	if _, ok := config.Listings[listing]; !ok {
		return
	}

	if _, ok := config.Listings[listing][level]; !ok {
		return
	}

	if count == 0 {
		delete(config.Listings[listing], level)

		if len(config.Listings[listing]) == 0 {
			delete(config.Listings, listing)
		}
	} else {
		config.Listings[listing][level] = count
	}

	runtime.EventsEmit(ctx, "listings_updated")

	Save()
}

func SetPrice(ctx context.Context, listing string, price string) {
	config.Prices[listing] = price
	runtime.EventsEmit(ctx, "config_updated")
	Save()
}

func SetLeague(ctx context.Context, league string) {
	config.League = League(league)
	runtime.EventsEmit(ctx, "config_updated")
	Save()
}

func SetName(ctx context.Context, name string) {
	config.Name = name
	runtime.EventsEmit(ctx, "config_updated")
	Save()
}

func SetStream(ctx context.Context, stream bool) {
	config.Stream = stream
	runtime.EventsEmit(ctx, "config_updated")
	Save()
}
