// +build !linux

package systray

import (
	"fmt"
	"github.com/Shitovdm/GitRsync/public/assets/src/icon"
	"github.com/Shitovdm/GitRsync/src/component/cmd/prompt"
	"github.com/Shitovdm/GitRsync/src/component/helper"
	"github.com/getlantern/systray"
	"io/ioutil"
	"time"
)

// Exit describes on exit app action.
func Exit() {
	now := time.Now()
	_ = ioutil.WriteFile(fmt.Sprintf(`./tmp/%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
}

// Ready describes on ready app action.
func Ready() {
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
			helper.OpenBrowser("http://localhost:8888")
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
			helper.OpenBrowser("https://github.com/Shitovdm/GitRsync")
			break
		case <-mDocs.ClickedCh:
			fmt.Println("Opening app specification...")
			helper.OpenBrowser("http://localhost:8888/docs")
			break
		case <-mRestart.ClickedCh:
			fmt.Println("Restarting application...")
			systray.Quit()
			time.Sleep(2 * time.Second)
			systray.Run(Ready, Exit)
			//return
		case <-mQuit.ClickedCh:
			fmt.Println("Closing application...")
			systray.Quit()
			return
		}
	}
}
