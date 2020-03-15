package main

import (
	"bufio"
	"fmt"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Configuration"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Helpers"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Logger"
	"io"
	"time"
)

func init() {
	Configuration.Init("GitRsync")
	Logger.Init()
}

func main() {

	//resource := "https://github.com/Shitovdm/rpc.git"
	//path := Configuration.BuildPlatformPath(fmt.Sprintf("/projects/%s", "rpc-test-sync"))

	rd, wr := io.Pipe()
	/*go func() {
		for i := 0; i < 10; i++ {
			_, _ = w.Write([]byte(fmt.Sprintf("ping %d \n", i)))
		}
		_ = w.Close()
	}()*/

	/*for i := 0; i < 10; i++ {
		_, _ = wr.Write([]byte(fmt.Sprintf("ping %d", i)))
	}*/

	//go Helpers.Exec(fmt.Sprintf("git clone %s %s", resource, path), rd, wr)
	go Helpers.Exec("curl ya.ru", rd, wr)

	/*go func() {
		scanner := bufio.NewScanner(rd)
		for scanner.Scan() {
			stdoutLine := scanner.Text()
			fmt.Println("line: " + stdoutLine)
		}
	}()*/


	go func() {
		for {
			time.Sleep(10 - time.Millisecond)
			scanner := bufio.NewScanner(rd)
			for scanner.Scan() {
				stdoutLine := scanner.Text()
				fmt.Println("line: " + stdoutLine)
			}

			/*buf := make([]byte, 100)
			n, err := rd.Read(buf)
			if err != nil {

			}
			fmt.Printf("%q \n", buf[:n])*/
		}
	}()

	/*buf := make([]byte, 100)
	n, err := r.Read(buf)
	if err != nil {

	}
	fmt.Printf("%q \n", buf[:n])

	buf2 := make([]byte, 100)
	n2, err := r.Read(buf2)
	if err != nil {

	}
	fmt.Printf("%q \n", buf2[:n2])
	_ = r.Close()*/

	time.Sleep(60* time.Second)
	//Application.StartServer()
}
