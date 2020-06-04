package main

import (
	"fmt"
	"github.com/Shitovdm/git-rsync/public/assets/src/icon"
	"github.com/Shitovdm/git-rsync/src/Application"
	"github.com/Shitovdm/git-rsync/src/Components/Configuration"
	"github.com/Shitovdm/git-rsync/src/Components/Helpers"
	"github.com/Shitovdm/git-rsync/src/Components/Logger"
	"github.com/getlantern/systray"
	"io/ioutil"
	"time"
    "github.com/gonutz/ide/w32"
)

func init() {
	Configuration.Init("GitRsync")
	Logger.Init()
}

func main() {
	//Application.StartServer()
	go Application.StartServer()
	hideConsole()
	systray.RunWithAppWindow("GitRsync", 1024, 768, onReady, onExit)
}

func onExit() {
	now := time.Now()
	ioutil.WriteFile(fmt.Sprintf(`./tmp/%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("GitRsync")
	systray.SetTooltip("GitRsync")

	mOpen := systray.AddMenuItem("Open GitRsync", "Open GitRsync UI")
	systray.AddSeparator()
	mOpenGit := systray.AddMenuItem("Project Page", "Open project page")
	mDocs := systray.AddMenuItem("Documentation", "Open documentation")
	systray.AddSeparator()
	mRestart := systray.AddMenuItem("Restart...", "Restart GitRsync")
	mQuit := systray.AddMenuItem("Quit GitRsync", "Quit GitRsync")

	for {
		select {
		case <-mOpen.ClickedCh:
			fmt.Println("Opening application UI...")
			Helpers.OpenBrowser("http://localhost:8888")
			break
		case <-mOpenGit.ClickedCh:
			fmt.Println("Opening app GIT page...")
			Helpers.OpenBrowser("https://github.com/Shitovdm/GitRsync")
			break
		case <-mDocs.ClickedCh:
			fmt.Println("Opening app specification...")
			Helpers.OpenBrowser("http://localhost:8888/docs")
			break
		case <-mRestart.ClickedCh:
			fmt.Println("Restarting application...")
			systray.Quit()
			time.Sleep(2 * time.Second)
			systray.Run(onReady, onExit)
			//return
		case <-mQuit.ClickedCh:
			fmt.Println("Closing application...")
			systray.Quit()
			return
		}
	}
}

func hideConsole() {
	console := w32.GetConsoleWindow()
	if console == 0 {
		return // no console attached
	}
	// If this application is the process that created the console window, then
	// this program was not compiled with the -H=windowsgui flag and on start-up
	// it created a console along with the main application window. In this case
	// hide the console window.
	// See
	// http://stackoverflow.com/questions/9009333/how-to-check-if-the-program-is-run-from-a-console
	_, consoleProcID := w32.GetWindowThreadProcessId(console)
	if w32.GetCurrentProcessId() == consoleProcID {
		w32.ShowWindowAsync(console, w32.SW_HIDE)
	}
}