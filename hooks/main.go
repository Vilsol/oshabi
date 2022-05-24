package hooks

import (
	"context"

	hook "github.com/robotn/gohook"
	"github.com/vilsol/oshabi/app"
	"github.com/vilsol/oshabi/config"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func InitializeHooks(ctx context.Context) {
	hook.Register(hook.KeyDown, []string{"ctrl", "j"}, func(e hook.Event) {
		listings, err := app.ReadFull(ctx)
		if err != nil {
			runtime.EventsEmit(ctx, "error", err.Error())
			return
		}

		if err := config.AddListings(ctx, listings); err != nil {
			runtime.EventsEmit(ctx, "error", err.Error())
			return
		}
	})

	go func() {
		s := hook.Start()
		<-hook.Process(s)

		hook.Start()
	}()
}
