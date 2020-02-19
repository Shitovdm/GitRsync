package main

import (
	"fmt"
	"github.com/Shitovdm/git-repo-exporter/helpers"
)

//var wg sync.WaitGroup // 1

/**
Step 1/8 : ARG  BUILDER=${BUILDER}
Step 2/8 : FROM ${BUILDER} as builder
 ---> 755e306a9382
Step 3/8 : FROM scratch
*/

func main() {

	command := "git clone https://github.com/jung-kurt/gofpdf.git"
	helpers.execCommand(command)

	if helpers.IsDirExists("") {
		fmt.Println("true")
	}

}
