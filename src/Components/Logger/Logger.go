package Logger

import (
	"github.com/Shitovdm/git-repo-exporter/src/Components/Configuration"
	"io/ioutil"
)

var logFileName = "runtime.log"

func GetRuntimeLogFile() string {
	return Configuration.BuildPlatformPath(logFileName)
}

func GetRuntimeLogs() ([]byte, error) {

	data, err := ioutil.ReadFile(GetRuntimeLogFile())
	if err != nil {
		return nil, err
	}

	return data, nil
}

func AddRuntimeLog(note string) {

}

func ClearRuntimeLogs() {
	_ = ioutil.WriteFile(GetRuntimeLogFile(), []byte(``), 0644)
}
