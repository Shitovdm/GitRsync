package Logger

import (
	"fmt"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Configuration"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Helpers"
	"github.com/Shitovdm/git-repo-exporter/src/Models"
	"io/ioutil"
	"time"
)

const timeFormatStr = "2006-01-02 15:04:05"

var RuntimeLogFile = "Runtime.json"
var sessionID = ""

func Init() {
	sessionID = Helpers.GenerateUuid()
}

func GetRuntimeLogFile() string {
	return Configuration.BuildPlatformPath(RuntimeLogFile)
}

func GetRuntimeLogs() []Models.RuntimeLog {
	runtimeLogs := make([]Models.RuntimeLog, 0)
	err := Configuration.Load(RuntimeLogFile, &runtimeLogs)
	if err != nil {
		err = Configuration.Save("Runtime.json", []map[string]interface{}{})
		if err != nil {
			panic(fmt.Sprintf("Error while creating new runtime log file! %s", err.Error()))
		}
		return []Models.RuntimeLog{}
	}

	return runtimeLogs
}

func AddRuntimeLog(sessionID string, level string, category string, message string) {
	runtimeLogs := GetRuntimeLogs()
	runtimeLogs = append(runtimeLogs, Models.RuntimeLog{
		SessionID: sessionID,
		Time:      time.Now().Format(timeFormatStr),
		Level:     level,
		Category:  category,
		Message:   message,
	})

	err := Configuration.Save(RuntimeLogFile, &runtimeLogs)
	if err != nil {
		panic("Error while saving runtime log!")
	}
}

func ClearRuntimeLogs() {
	_ = ioutil.WriteFile(GetRuntimeLogFile(), []byte(``), 0644)
}

func SetLogLevel(level string, str string) string {

	switch level {
	case "info":
		return fmt.Sprintf("\x1b[92m%s\x1b[0m", str)
	case "error":
		return fmt.Sprintf("\x1b[91m%s\x1b[0m", str)
	default:
		return fmt.Sprintf("\x1b[91m%s\x1b[0m", str)
	}


}
