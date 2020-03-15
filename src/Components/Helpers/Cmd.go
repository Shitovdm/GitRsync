package Helpers

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

func Exec(command string, rd *io.PipeReader, wr *io.PipeWriter) {
	var wg sync.WaitGroup
	wg.Add(1)
	go execCmd(command, &wg, rd, wr)
	wg.Wait()
}

func execCmd(command string, wg *sync.WaitGroup, rd *io.PipeReader, wr *io.PipeWriter) {
	defer wg.Done()
	fmt.Println("command is ", command)
	args := strings.Fields(command)
	head := args[0]
	args = args[1:]

	//outPipe := os.NewFile(uintptr(syscall.Stdout), "/tmp/outPipe")

	cmd := exec.Command(head, args...)
	cmd.Stdout = wr
	//cmd.Stderr = wr
	//cmd.Stdin = rd
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to exec command! %s\n", err.Error())
		os.Exit(1)
	}

	//fmt.Printf("Result: %v / %v", writer., stderr.String())
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

