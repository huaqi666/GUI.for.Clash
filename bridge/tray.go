//go:build windows

package bridge

import (
	"embed"
	"log"
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
)

var (
	appTray *application.SystemTray
	appMenu *application.Menu
)

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

	appTray = app.NewSystemTray()
	appMenu = app.NewMenu()

	appTray.SetMenu(appMenu)

	appTray.OnClick(func() { app.Events.Emit(&application.WailsEvent{Name: "tray:click"}) })
	appTray.OnDoubleClick(func() { app.Events.Emit(&application.WailsEvent{Name: "tray:dblclick"}) })
	appTray.OnRightDoubleClick(func() { app.Events.Emit(&application.WailsEvent{Name: "tray:rdblclick"}) })
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
		// tray.SetLabel(tray.Title)
		// runtime.WindowSetTitle(a.Ctx, tray.Title)
	}
	if tray.Tooltip != "" {
		// systray.SetTooltip(tray.Tooltip)
	}
}

func (a *App) UpdateTrayMenus(menus []MenuItem) {
	log.Printf("UpdateTrayMenus")

	appMenu = a.Ctx.NewMenu()

	for _, menu := range menus {
		createMenuItem(menu, a, appMenu)
	}

	appTray.SetMenu(appMenu)
}

func createMenuItem(menu MenuItem, a *App, parent *application.Menu) {
	if menu.Hidden {
		return
	}

	switch menu.Type {
	case "item":
		if len(menu.Children) == 0 {
			parent.Add(menu.Text).OnClick(func(ctx *application.Context) {
				a.Ctx.Events.Emit(&application.WailsEvent{Name: menu.Event})
			})
			return
		}
		subMenu := parent.AddSubmenu(menu.Text)
		for _, child := range menu.Children {
			createMenuItem(child, a, subMenu)
		}
	case "separator":
		appMenu.AddSeparator()
	}
}

func (a *App) ExitApp() {
	a.Ctx.Quit()
}
