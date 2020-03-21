package Helpers

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

func Exec(command string) {
	var wg sync.WaitGroup
	wg.Add(1)
	go execCmd(command, &wg)
	wg.Wait()
}

func execCmd(command string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Command is:", command)
	args := strings.Fields(command)
	head := args[0]
	args = args[1:]
	cmd := exec.Command(head, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to exec command! %s\n", err.Error())
		os.Exit(1)
	}
}

func Copy() {

}

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
