//go:build windows

package bridge

import (
	"embed"
	"log"
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
)

var systemTray *application.SystemTray

func InitTray(app *application.App, icon []byte, fs embed.FS) {
	icons := [6][2]string{
		{"frontend/dist/icons/tray_normal_light.ico", "data/.cache/icons/tray_normal_light.ico"},
		{"frontend/dist/icons/tray_normal_dark.ico", "data/.cache/icons/tray_normal_dark.ico"},
		{"frontend/dist/icons/tray_proxy_light.ico", "data/.cache/icons/tray_proxy_light.ico"},
		{"frontend/dist/icons/tray_proxy_dark.ico", "data/.cache/icons/tray_proxy_dark.ico"},
		{"frontend/dist/icons/tray_tun_light.ico", "data/.cache/icons/tray_tun_light.ico"},
		{"frontend/dist/icons/tray_tun_dark.ico", "data/.cache/icons/tray_tun_dark.ico"},
	}

	os.MkdirAll(GetPath("data/.cache/icons"), os.ModePerm)

	for _, item := range icons {
		path := GetPath(item[1])
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Printf("InitTray [Icon]: %s", item[1])
			b, _ := fs.ReadFile(item[0])
			os.WriteFile(path, b, os.ModePerm)
		}
	}

	systemTray = app.NewSystemTray()
	b, _ := fs.ReadFile("frontend/dist/wails.png")
	systemTray.SetTemplateIcon(b)

	menu := app.NewMenu()
	menu.AddSubmenu("Test")

	systemTray.SetMenu(menu)
	systemTray.OnClick(func() {
		println("test")
	})
}

func (a *App) UpdateTray(tray TrayContent) {
	println("teray, %v", tray.Icon)
	if tray.Icon != "" {
		// icon, err := os.ReadFile(GetPath(tray.Icon))
		// if err == nil {
		// 	systemTray.SetTemplateIcon(icon)
		// }
	}
	if tray.Title != "" {
		systemTray.SetLabel(tray.Title)
		// runtime.WindowSetTitle(a.Ctx, tray.Title)
	}
	if tray.Tooltip != "" {
		// systray.SetTooltip(tray.Tooltip)
	}
}

func (a *App) ExitApp() {
	a.Ctx.Quit()
}
