// +build windows

package prompt

import "syscall"

// ChangeConsoleVisibility controls console visibility.
func ChangeConsoleVisibility(visibility bool) {
	getConsoleWindow := syscall.NewLazyDLL("kernel32.dll").NewProc("GetConsoleWindow")
	showWindow := syscall.NewLazyDLL("user32.dll").NewProc("ShowWindow")
	if getConsoleWindow.Find() == nil && showWindow.Find() == nil {
		consoleWindow, _, _ := getConsoleWindow.Call()
		if consoleWindow != 0 {
			if visibility {
				_, _, _ = showWindow.Call(consoleWindow, 1)
			} else {
				_, _, _ = showWindow.Call(consoleWindow, 0)
			}
		}
	}
}
