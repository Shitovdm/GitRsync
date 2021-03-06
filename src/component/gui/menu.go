package gui

import (
	"github.com/gin-gonic/gin"
	"strings"
)

// Menu describes left menu bar items.
var Menu = []map[string]string{
	{
		"slug":   "",
		"href":   "/",
		"icon":   "dashboard",
		"text":   "Dashboard",
		"active": "",
	},
	{
		"slug":   "repositories",
		"href":   "/repositories",
		"icon":   "storage",
		"text":   "Repositories",
		"active": "",
	},
	{
		"slug":   "platforms",
		"href":   "/platforms",
		"icon":   "dns",
		"text":   "Platforms",
		"active": "",
	},
	{
		"slug":   "logs",
		"href":   "/logs",
		"icon":   "assignments",
		"text":   "Logs",
		"active": "",
	},
	{
		"slug":   "settings",
		"href":   "/settings",
		"icon":   "settings",
		"text":   "Settings",
		"active": "",
	},
	{
		"slug":   "docs",
		"href":   "/docs",
		"icon":   "comment",
		"text":   "Docs",
		"active": "",
	},
	{
		"slug":   "about",
		"href":   "/about",
		"icon":   "help",
		"text":   "About",
		"active": "",
	},
}

// GetMenu gets menu items.
func GetMenu(c *gin.Context) []map[string]string {
	path := strings.Split(strings.Trim(c.Request.URL.Path, "/"), "/")
	for _, menuItem := range Menu {
		if menuItem["slug"] == path[0] {
			menuItem["active"] = "active"
		} else {
			menuItem["active"] = ""
		}
	}
	return Menu
}
