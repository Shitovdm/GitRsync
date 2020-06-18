package main

import (
	"fmt"
	"github.com/Shitovdm/GitRsync/public/assets/src/icon"
	"github.com/Shitovdm/GitRsync/src/Components/cmd/prompt"
	"github.com/Shitovdm/GitRsync/src/application"
	"github.com/Shitovdm/GitRsync/src/components/configuration"
	"github.com/Shitovdm/GitRsync/src/components/helpers"
	"github.com/Shitovdm/GitRsync/src/components/logger"
	"github.com/getlantern/systray"
	"io/ioutil"
	"time"
)

func init() {
	configuration.Init("GitRsync")
	logger.Init()
}

func main() {

	application.StartServer()

	//go Application.StartServer()
	//	It`s hiding command prompt during running app (only windows).
	prompt.ChangeConsoleVisibility(false)
	//systray.RunWithAppWindow("GitRsync", 1024, 768, onReady, onExit)
}

func onExit() {
	now := time.Now()
	_ = ioutil.WriteFile(fmt.Sprintf(`./tmp/%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("GitRsync")
	systray.SetTooltip("GitRsync")

	mOpen := systray.AddMenuItem("Open GitRsync", "Open GitRsync UI")
	changeConsoleVisibility := systray.AddMenuItem("Show Terminal", "Hide Terminal")
	systray.AddSeparator()
	mOpenGit := systray.AddMenuItem("Project Page", "Open project page")
	mDocs := systray.AddMenuItem("Documentation", "Open documentation")
	systray.AddSeparator()
	mRestart := systray.AddMenuItem("Restart...", "Restart GitRsync")
	mQuit := systray.AddMenuItem("Quit GitRsync", "Quit GitRsync")

	terminalVisibility := false
	for {
		select {
		case <-mOpen.ClickedCh:
			fmt.Println("Opening application UI...")
			helpers.OpenBrowser("http://localhost:8888")
			break
		case <-changeConsoleVisibility.ClickedCh:
			fmt.Println("Changing console visibility...")
			if terminalVisibility {
				changeConsoleVisibility.SetTitle("Show Terminal")
				prompt.ChangeConsoleVisibility(false)
				terminalVisibility = false
			} else {
				changeConsoleVisibility.SetTitle("Hide Terminal")
				prompt.ChangeConsoleVisibility(true)
				terminalVisibility = true
			}
			break
		case <-mOpenGit.ClickedCh:
			fmt.Println("Opening app GIT page...")
			helpers.OpenBrowser("https://github.com/Shitovdm/GitRsync")
			break
		case <-mDocs.ClickedCh:
			fmt.Println("Opening app specification...")
			helpers.OpenBrowser("http://localhost:8888/docs")
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
