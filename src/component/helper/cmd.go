package helper

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

// Exec execute the shell command.
func Exec(command string) {
	var wg sync.WaitGroup
	wg.Add(1)
	go execCmd(command, &wg)
	wg.Wait()
}

// execCmd runs custom shell command.
func execCmd(command string, wg *sync.WaitGroup) {
	defer wg.Done()
	args := strings.Fields(command)
	head := args[0]
	args = args[1:]
	cmd := exec.Command(head, args...)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to exec command! %s\n", err.Error())
		os.Exit(1)
	}
}

// OpenBrowser opens web browser in different platforms.
func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
