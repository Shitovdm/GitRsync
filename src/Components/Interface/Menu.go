package Interface

import (
	"github.com/gin-gonic/gin"
	"strings"
)

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
		"slug":   "about",
		"href":   "/about",
		"icon":   "help",
		"text":   "About",
		"active": "",
	},
}

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
