// +build windows

package Prompt

import "syscall"

func ChangeConsoleVisibility(visibility bool) {
	getConsoleWindow := syscall.NewLazyDLL("kernel32.dll").NewProc("GetConsoleWindow")
	showWindow := syscall.NewLazyDLL("user32.dll").NewProc("ShowWindow")
	if getConsoleWindow.Find() == nil && showWindow.Find() == nil {
		hwnd, _, _ := getConsoleWindow.Call()
		if hwnd != 0 {
			if visibility {
				_, _, _ = showWindow.Call(hwnd, 1)
			} else {
				_, _, _ = showWindow.Call(hwnd, 0)
			}
		}
	}
}
