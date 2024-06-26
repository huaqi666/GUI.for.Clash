package main

import (
	"embed"
	"guiforcores/bridge"
	"log"
	"log/slog"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed frontend/dist
var assets embed.FS

//go:embed frontend/dist/favicon.png
var icon []byte

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
		Icon:        icon,
		LogLevel:    slog.LevelWarn,
		Services: []application.Service{
			application.NewService(appService),
		},
		Assets: application.AssetOptions{
			Handler:        application.AssetFileServerFS(assets),
			Middleware:     appService.BridgeHTTPApi,
			DisableLogging: true,
		},
	})

	appService.Ctx = app
	bridge.InitApp()
	bridge.InitTray(app, icon, assets)
	bridge.InitNotification(assets)
	bridge.InitScheduledTasks()

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Name:                   "Main",
		URL:                    "/",
		MinWidth:               600,
		MinHeight:              400,
		Centered:               true,
		DisableResize:          false,
		OpenInspectorOnStartup: true,
		Title:                  bridge.Env.AppName,
		Width:                  bridge.Config.Width,
		Height:                 bridge.Config.Height,
		Frameless:              bridge.Env.OS == "windows",
		Hidden:                 bridge.Env.FromTaskSch && bridge.Config.WindowStartState == 2,
		BackgroundType:         application.BackgroundTypeTranslucent,
		BackgroundColour:       application.NewRGBA(255, 255, 255, 1),
		StartState:             application.WindowState(bridge.Config.WindowStartState),
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
				Name: "onBeforeExitApp",
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
	})

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
