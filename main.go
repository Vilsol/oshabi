package main

import (
	"embed"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

var _ logger.Logger = (*wrappedLogger)(nil)

type wrappedLogger struct {
}

func (w wrappedLogger) Print(message string) {
	log.Print(message)
}

func (w wrappedLogger) Trace(message string) {
	log.Trace().Msg(message)
}

func (w wrappedLogger) Debug(message string) {
	log.Debug().Msg(message)
}

func (w wrappedLogger) Info(message string) {
	log.Info().Msg(message)
}

func (w wrappedLogger) Warning(message string) {
	log.Warn().Msg(message)
}

func (w wrappedLogger) Error(message string) {
	log.Error().Msg(message)
}

func (w wrappedLogger) Fatal(message string) {
	log.Fatal().Msg(message)
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	appLogPath := filepath.Join(configDir, "oshabi", "log.log")
	if err := os.MkdirAll(filepath.Dir(appLogPath), 0777); err != nil {
		if !os.IsExist(err) {
			panic(err)
		}
	}

	f, err := os.OpenFile(appLogPath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	output := zerolog.ConsoleWriter{Out: f, TimeFormat: time.RFC3339}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()

	app := NewApp()

	err = wails.Run(&options.App{
		Title:             "Oshabi",
		Width:             1280,
		Height:            650,
		MinWidth:          850,
		MinHeight:         600,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		RGBA:              &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		Assets:            assets,
		LogLevel:          logger.DEBUG,
		Logger:            wrappedLogger{},
		OnStartup:         app.startup,
		OnDomReady:        app.domReady,
		OnShutdown:        app.shutdown,
		Bind: []interface{}{
			app,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
	})

	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
