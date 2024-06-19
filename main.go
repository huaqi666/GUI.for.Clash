package main

import (
	"embed"
	"guiforcores/bridge"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed frontend/dist
var assets embed.FS

var isStartup = true

func main() {
	bridge.InitBridge()

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        "GUI.for.Cores",
		Description: "A GUI program developed by vue3 + wails3.",
		Services: []application.Service{
			application.NewService(&bridge.App{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	systemTray := app.NewSystemTray()
	b, _ := assets.ReadFile("frontend/dist/wails.png")
	systemTray.SetTemplateIcon(b)

	menu := app.NewMenu()
	menu.AddSubmenu("Test")

	systemTray.SetMenu(menu)
	systemTray.OnClick(func() {
		println("test")
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: bridge.Env.AppName,
		Height: func() int {
			if bridge.Config.Height != 0 {
				return bridge.Config.Height
			}
			if bridge.Env.OS == "linux" {
				return 510
			}
			return 540
		}(),
		MinWidth:      600,
		MinHeight:     400,
		Frameless:     bridge.Env.OS == "windows",
		DisableResize: false,
		// StartHidden: func() bool {
		// 	if bridge.Env.FromTaskSch {
		// 		return bridge.Config.WindowStartState == 2
		// 	}
		// 	return false
		// }(),
		// WindowStartState: func() options.WindowStartState {
		// 	if bridge.Env.FromTaskSch {
		// 		return options.WindowStartState(bridge.Config.WindowStartState)
		// 	}
		// 	return 0
		// }(),
		// Windows: &windows.Options{
		// 	WebviewIsTransparent: true,
		// 	WindowIsTranslucent:  true,
		// 	BackdropType:         windows.Acrylic,
		// },
		Width: func() int {
			if bridge.Config.Width != 0 {
				return bridge.Config.Width
			}
			return 800
		}(),
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
			// About:
			// About: &mac.AboutInfo{
			// 	Title:   bridge.Env.AppName,
			// 	Message: "© 2024 GUI.for.Cores",
			// 	Icon:    icon,
			// },
		},
		BackgroundColour: application.NewRGBA(255, 255, 255, 1),
		// Linux: &linux.Options{
		// 	Icon:                icon,
		// 	WindowIsTranslucent: false,
		// },
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
		// OnStartup: func(ctx context.Context) {
		// 	runtime.LogSetLogLevel(ctx, logger.INFO)
		// 	app.Ctx = ctx
		// 	bridge.InitTray(app, icon, assets)
		// 	bridge.InitScheduledTasks()
		// 	bridge.InitNotification(assets)
		// },
		// OnDomReady: func(ctx context.Context) {
		// 	if isStartup {
		// 		runtime.EventsEmit(ctx, "onStartup")
		// 		isStartup = false
		// 	}
		// },
		// OnBeforeClose: func(ctx context.Context) (prevent bool) {
		// 	runtime.EventsEmit(ctx, "beforeClose")
		// 	return true
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
