package helpers

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

func execCommand(command string) {
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