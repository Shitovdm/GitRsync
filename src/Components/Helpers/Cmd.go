package Helpers

import (
	"fmt"
	"log"
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

func execCmd(cmd string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("command is ", cmd)
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)
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
