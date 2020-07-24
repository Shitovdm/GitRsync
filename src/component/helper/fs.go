package helper

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

// GetCurrentPath returns current location path.
func GetCurrentPath() string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err.Error())
	}
	exPath := filepath.Dir(ex)

	return exPath
}

// IsDirExists returns dir exists flag.
func IsDirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

// CreateNewDir creates new folder in current file system.
func CreateNewDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// RemoveDir removes folder from current file system.
func RemoveDir(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	return nil
}

// IsFileExists returns is file exists flag.
func IsFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

// WriteFile writes file data.
func WriteFile(filename string, content string) {
	_ = ioutil.WriteFile(filename, []byte(content), 0644)
}

// Move moves file to new location.
func Move(oldLocation string, newLocation string) error {

	err := os.Rename(oldLocation, newLocation)
	if err != nil {
		return err
	}

	return nil
}

// CopyFile creates file copy.
func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

// CopyDirContent copies dir content.
func CopyDirContent(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	srcinfo, err = os.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, srcinfo.Mode())
	if err != nil {
		return err
	}

	fds, err = ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, fd := range fds {
		_ = path.Join(src, fd.Name())
		_ = path.Join(dst, fd.Name())
	}

	return nil
}

// ExploreDir opens dir in Explorer.
func ExploreDir(path string) {
	cmd := exec.Command(`explorer`, `/select,`, path)
	_ = cmd.Run()
}
