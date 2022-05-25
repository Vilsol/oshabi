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
	LanguageChinese    = Language("chi_sim")
	LanguageKorean     = Language("kor")
	LanguageJapanese   = Language("jpn")
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
	Shortcut []string               `json:"shortcut"`
}

var Cfg Config

func InitConfig() error {
	if err := Load(); err != nil {
		return err
	}

	if Cfg.Version == 0 {
		Cfg = Config{
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

	if len(Cfg.Shortcut) == 0 {
		Cfg.Shortcut = []string{"ctrl", "j"}
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

	if err := json.Unmarshal(file, &Cfg); err != nil {
		return errors.Wrap(err, "failed parsing config.json")
	}

	return nil
}

func Save() error {
	b, err := json.MarshalIndent(Cfg, "", "  ")
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
		if _, ok := Cfg.Listings[strType]; !ok {
			Cfg.Listings[strType] = make(map[int]int)
		}

		if _, ok := Cfg.Listings[strType][listing.Level]; !ok {
			Cfg.Listings[strType][listing.Level] = 0
		}

		Cfg.Listings[strType][listing.Level] += listing.Count
	}

	runtime.EventsEmit(ctx, "listings_updated")
	return Save()
}

func ClearListings(ctx context.Context) error {
	Cfg.Listings = make(map[string]map[int]int)
	runtime.EventsEmit(ctx, "listings_updated")
	return Save()
}

func Get() Config {
	return Cfg
}

func SetListing(ctx context.Context, listing string, level int, count int) error {
	if _, ok := Cfg.Listings[listing]; !ok {
		return nil
	}

	if _, ok := Cfg.Listings[listing][level]; !ok {
		return nil
	}

	if count == 0 {
		delete(Cfg.Listings[listing], level)

		if len(Cfg.Listings[listing]) == 0 {
			delete(Cfg.Listings, listing)
		}
	} else {
		Cfg.Listings[listing][level] = count
	}

	runtime.EventsEmit(ctx, "listings_updated")
	return Save()
}

func SetPrice(ctx context.Context, listing string, price string) error {
	Cfg.Prices[listing] = price
	runtime.EventsEmit(ctx, "config_updated")
	return Save()
}

func SetLeague(ctx context.Context, league string) error {
	Cfg.League = League(league)
	runtime.EventsEmit(ctx, "config_updated")
	return Save()
}

func SetName(ctx context.Context, name string) error {
	Cfg.Name = name
	runtime.EventsEmit(ctx, "config_updated")
	return Save()
}

func SetStream(ctx context.Context, stream bool) error {
	Cfg.Stream = stream
	runtime.EventsEmit(ctx, "config_updated")
	return Save()
}

func SetScaling(ctx context.Context, scaling float64) error {
	Cfg.Scaling = scaling
	runtime.EventsEmit(ctx, "config_updated")
	return Save()
}

func SetDisplay(ctx context.Context, display int) error {
	Cfg.Display = display
	runtime.EventsEmit(ctx, "config_updated")
	return Save()
}

func SetLanguage(ctx context.Context, language string) error {
	Cfg.Language = Language(language)
	runtime.EventsEmit(ctx, "config_updated")
	return Save()
}

func SetShortcut(ctx context.Context, shortcut []string) error {
	Cfg.Shortcut = shortcut
	runtime.EventsEmit(ctx, "config_updated")
	return Save()
}
