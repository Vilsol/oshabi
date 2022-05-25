package hooks

import (
	"context"
	"time"

	hook "github.com/robotn/gohook"
	"github.com/vilsol/oshabi/app"
	"github.com/vilsol/oshabi/config"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var scanning = false

func InitializeHooks(ctx context.Context) {
	lastTime := time.Now()
	hook.Register(hook.KeyDown, config.Get().Shortcut, func(e hook.Event) {
		if time.Since(lastTime) < time.Second*2 {
			return
		}
		lastTime = time.Now()

		if scanning {
			return
		}
		scanning = true

		go func() {
			defer func() {
				scanning = false
			}()

			listings, err := app.ReadFull(ctx)
			if err != nil {
				runtime.EventsEmit(ctx, "error", err.Error())
				return
			}

			if err := config.AddListings(ctx, listings); err != nil {
				runtime.EventsEmit(ctx, "error", err.Error())
				return
			}
		}()
	})

	go func() {
		s := hook.Start()
		<-hook.Process(s)

		hook.Start()
	}()
}

var unhooking = false

func UpdateShortcut(ctx context.Context) {
	if unhooking {
		return
	}

	unhooking = true
	hook.End()
	time.Sleep(time.Millisecond * 500)
	InitializeHooks(ctx)
	unhooking = false
}
