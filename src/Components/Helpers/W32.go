package Helpers

import (
	"syscall"
	"unsafe"
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	getConsoleWindow         = kernel32.NewProc("GetConsoleWindow")
	getWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	getCurrentProcessId      = kernel32.NewProc("GetCurrentProcessId")
	showWindowAsync          = user32.NewProc("ShowWindowAsync")
)

func GetConsoleWindow() uintptr {
	ret, _, _ := getConsoleWindow.Call()
	return ret
}

func GetWindowThreadProcessId(hwnd uintptr) (uintptr, uint32) {
	var processId uint32
	ret, _, _ := getWindowThreadProcessId.Call(
		hwnd,
		uintptr(unsafe.Pointer(&processId)),
	)
	return ret, processId
}

func GetCurrentProcessId() uint32 {
	id, _, _ := getCurrentProcessId.Call()
	return uint32(id)
}

func ShowWindowAsync(window, commandShow uintptr) bool {
	ret, _, _ := showWindowAsync.Call(window, commandShow)
	return ret != 0
}

func HideConsole() {
	console := GetConsoleWindow()
	if console == 0 {
		return
	}
	_, consoleProcID := GetWindowThreadProcessId(console)
	if GetCurrentProcessId() == consoleProcID {
		ShowWindowAsync(console, 0)
	}
}
