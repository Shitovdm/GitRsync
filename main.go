package main

import (
	"fmt"
	"github.com/Shitovdm/git-rsync/public/assets/src/icon"
	"github.com/Shitovdm/git-rsync/src/Application"
	"github.com/Shitovdm/git-rsync/src/Components/Configuration"
	"github.com/Shitovdm/git-rsync/src/Components/Helpers"
	"github.com/Shitovdm/git-rsync/src/Components/Logger"
	"github.com/getlantern/systray"
	"github.com/gonutz/ide/w32"
	"io/ioutil"
	"time"
)

func init() {
	Configuration.Init("GitRsync")
	Logger.Init()
}

func main() {
	//Application.StartServer()
	go Application.StartServer()
	Application.HideConsole()
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
