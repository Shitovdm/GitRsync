// +build darwin

package systray

import (
	"github.com/Shitovdm/GitRsync/src/component/conf"
	"github.com/getlantern/systray"
)

// Init creates application icon in system tray.
func Init() {
	systray.RunWithAppWindow(conf.GetApplicationName(), 1024, 768, Ready, Exit)
}
