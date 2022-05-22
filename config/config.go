package config

import (
	"context"
	"encoding/json"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/vilsol/oshabi/types"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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
	Display  int                    `json:"display"`
}

var config Config

func InitConfig() error {
	if err := Load(); err != nil {
		return err
	}

	if config.Version == 0 {
		config = Config{
			Version:  1,
			Scaling:  1,
			Prices:   map[string]string{},
			Listings: map[string]map[int]int{},
			Language: LanguageEnglish,
			League:   LeagueSoftcore,
			Messages: map[string]string{},
			Display:  0,
		}

		if err := Save(); err != nil {
			return err
		}
	}

	return nil
}

func Load() error {
	dir, err := GetConfigDir()
	if err != nil {
		return err
	}

	file, err := os.ReadFile(path.Join(dir, "config.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errors.Wrap(err, "failed reading config.json")
	}

	if err := json.Unmarshal(file, &config); err != nil {
		return errors.Wrap(err, "failed parsing config.json")
	}

	return nil
}

func Save() error {
	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed serializing config.json")
	}

	dir, err := GetConfigDir()
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(dir, "config.json"), b, 0777)
	if err != nil {
		return errors.Wrap(err, "failed writing config.json")
	}
	return nil
}

func GetConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()

	if err != nil {
		return "", errors.Wrap(err, "failed to find user config directory")
	}

	finalDir := path.Join(configDir, "oshabi")

	if err := os.MkdirAll(finalDir, 0777); err != nil {
		if !os.IsExist(err) {
			return "", errors.Wrap(err, "failed making config directory")
		}
	}

	return finalDir, nil
}

func AddListings(ctx context.Context, listings []types.ParsedListing) error {
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

	runtime.EventsEmit(ctx, "listings_updated")
	return Save()
}

func ClearListings(ctx context.Context) error {
	config.Listings = make(map[string]map[int]int)
	runtime.EventsEmit(ctx, "listings_updated")
	return Save()
}

func Get() Config {
	return config
}

func SetListing(ctx context.Context, listing string, level int, count int) error {
	if _, ok := config.Listings[listing]; !ok {
		return nil
	}

	if _, ok := config.Listings[listing][level]; !ok {
		return nil
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
	return Save()
}

func SetPrice(ctx context.Context, listing string, price string) error {
	config.Prices[listing] = price
	runtime.EventsEmit(ctx, "config_updated")
	return Save()
}

func SetLeague(ctx context.Context, league string) error {
	config.League = League(league)
	runtime.EventsEmit(ctx, "config_updated")
	return Save()
}

func SetName(ctx context.Context, name string) error {
	config.Name = name
	runtime.EventsEmit(ctx, "config_updated")
	return Save()
}

func SetStream(ctx context.Context, stream bool) error {
	config.Stream = stream
	runtime.EventsEmit(ctx, "config_updated")
	return Save()
}

func SetScaling(ctx context.Context, scaling float64) error {
	config.Scaling = scaling
	runtime.EventsEmit(ctx, "config_updated")
	return Save()
}

func SetDisplay(ctx context.Context, display int) error {
	config.Display = display
	runtime.EventsEmit(ctx, "config_updated")
	return Save()
}
