package main

import (
	"github.com/Shitovdm/GitRsync/src/app"
	"github.com/Shitovdm/GitRsync/src/component/cmd/prompt"
	"github.com/Shitovdm/GitRsync/src/component/conf"
	"github.com/Shitovdm/GitRsync/src/component/logger"
)

func init() {
	conf.Init("GitRsync")
	logger.Init()
}

func main() {
	// Adding application icon to system tray (macos, windows).
	//go systray.Init()

	// It`s hiding command prompt during running app (only windows).
	prompt.ChangeConsoleVisibility(false)
	// Start serve application backend server.
	app.Serve()
}
