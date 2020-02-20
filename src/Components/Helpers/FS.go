package Helpers

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func GetCurrentPath() string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err.Error())
	}
	exPath := filepath.Dir(ex)

	return exPath
}

func IsDirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

func IsFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func WriteFile(filename string, content string) {
	_ = ioutil.WriteFile(filename, []byte(content), 0644)
}
