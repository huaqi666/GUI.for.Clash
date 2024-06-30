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
	src := "frontend/dist/icons/"
	dst := "data/.cache/icons/"

	icons := []string{
		"tray_normal_light.png",
		"tray_normal_dark.png",
		"tray_proxy_light.png",
		"tray_proxy_dark.png",
		"tray_tun_light.png",
		"tray_tun_dark.png",
	}

	os.MkdirAll(GetPath(dst), os.ModePerm)

	for _, icon := range icons {
		path := GetPath(dst + icon)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Printf("InitTray [Icon]: %s", src+icon)
			b, _ := fs.ReadFile(src + icon)
			os.WriteFile(path, b, os.ModePerm)
		}
	}

	appTray = app.NewSystemTray()
	appMenu = app.NewMenu()

	appTray.SetIcon(icon)
	appTray.SetDarkModeIcon(icon)
	appTray.SetMenu(appMenu)
	appTray.OnClick(func() {
		win := app.GetWindowByName("Main")
		win.UnMinimise()
		win.Show()
	})
}

func (a *App) UpdateTray(tray TrayContent) {
	if tray.Icon != "" {
		icon, err := os.ReadFile(GetPath(tray.Icon))
		if err == nil {
			appTray.SetIcon(icon)
			appTray.SetDarkModeIcon(icon)
		}
	}
	if tray.Title != "" {
		appTray.SetLabel(tray.Title)
		a.Ctx.GetWindowByName("Main").SetTitle(tray.Title)
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

	if len(menu.Children) != 0 {
		subMenu := parent.AddSubmenu(menu.Text)
		for _, child := range menu.Children {
			createMenuItem(child, a, subMenu)
		}
		return
	}

	onClick := func(ctx *application.Context) {
		a.Ctx.Events.Emit(&application.WailsEvent{Name: menu.Event})
	}

	switch menu.Type {
	case "item":
		parent.Add(menu.Text).SetChecked(menu.Checked).OnClick(onClick)
	case "radio":
		parent.AddRadio(menu.Text, menu.Checked).OnClick(onClick)
	case "separator":
		appMenu.AddSeparator()
	}

}
