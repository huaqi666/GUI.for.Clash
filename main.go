package main

import (
	"embed"
	"guiforcores/bridge"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed frontend/dist
var assets embed.FS

//go:embed frontend/dist/favicon.ico
var icon []byte

var isStartup = true

func main() {
	appService := &bridge.App{}

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        "GUI.for.Cores",
		Description: "A GUI program developed by vue3 + wails3.",
		Services: []application.Service{
			application.NewService(appService),
		},
		Assets: application.AssetOptions{
			Handler:    application.AssetFileServerFS(assets),
			Middleware: appService.BridgeHTTPApi,
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	appService.Ctx = app
	bridge.InitBridge()
	bridge.InitTray(app, icon, assets)
	bridge.InitNotification(assets)
	bridge.InitScheduledTasks()

	app.On(events.Common.ApplicationStarted, func(event *application.Event) {
		println(isStartup)
		if isStartup {
			app.Events.Emit(&application.WailsEvent{
				Name: "onStartup",
			})
			isStartup = false
		}
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:                  bridge.Env.AppName,
		MinWidth:               600,
		MinHeight:              400,
		Frameless:              bridge.Env.OS == "windows",
		DisableResize:          false,
		BackgroundType:         application.BackgroundTypeTranslucent,
		BackgroundColour:       application.NewRGBA(255, 255, 255, 1),
		OpenInspectorOnStartup: true,
		Width: func() int {
			if bridge.Config.Width != 0 {
				return bridge.Config.Width
			}
			return 800
		}(),
		Height: func() int {
			if bridge.Config.Height != 0 {
				return bridge.Config.Height
			}
			if bridge.Env.OS == "linux" {
				return 510
			}
			return 540
		}(),
		StartState: func() application.WindowState {
			if bridge.Env.FromTaskSch {
				return application.WindowState(bridge.Config.WindowStartState)
			}
			return 0
		}(),
		Hidden: func() bool {
			if bridge.Env.FromTaskSch {
				return bridge.Config.WindowStartState == 2
			}
			return false
		}(),
		Windows: application.WindowsWindow{
			BackdropType: application.Acrylic,
		},
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		Linux: application.LinuxWindow{
			Icon:                icon,
			WindowIsTranslucent: true,
			WebviewGpuPolicy:    application.WebviewGpuPolicyNever,
		},
		ShouldClose: func(window *application.WebviewWindow) bool {
			appService.Ctx.Events.Emit(&application.WailsEvent{
				Name: "beforeClose",
			})
			return true
		},
		// SingleInstanceLock: &options.SingleInstanceLock{
		// 	UniqueId: func() string {
		// 		if bridge.Config.MultipleInstance {
		// 			return uuid.New().String()
		// 		}
		// 		return bridge.Env.AppName
		// 	}(),
		// 	OnSecondInstanceLaunch: func(data options.SecondInstanceData) {
		// 		runtime.Show(app.Ctx)
		// 		runtime.EventsEmit(app.Ctx, "launchArgs", data.Args)
		// 	},
		// },
		URL: "/",
	})

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
