package Logger

import (
	"errors"
	"fmt"
	"github.com/Shitovdm/git-rsync/src/Components/Configuration"
	"github.com/Shitovdm/git-rsync/src/Components/Helpers"
	"github.com/Shitovdm/git-rsync/src/Models"
	"io/ioutil"
	"time"
)

const timeFormatStr = "2006-01-02 15:04:05"

var RuntimeLogFile = "RuntimeLogs.json"
var sessionID = ""
var runtimeLogNote = ""

func Init() {
	sessionID = Helpers.GenerateUuid()
	runtimeLogNote = ""
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

	runtimeLog := Models.RuntimeLog{
		SessionID: sessionID,
		Time:      time.Now().Format(timeFormatStr),
		Level:     level,
		Category:  category,
		Message:   message,
	}

	runtimeLogNote = BuildRuntimeLogNote(runtimeLog)
	runtimeLogs := GetRuntimeLogs()
	runtimeLogs = append(runtimeLogs, runtimeLog)

	fmt.Println(fmt.Sprintf("[%s][%s][%s] %s", runtimeLog.Time, runtimeLog.Level, runtimeLog.Category, runtimeLog.Message))
	err := Configuration.Save(RuntimeLogFile, &runtimeLogs)
	if err != nil {
		panic("Error while saving runtime log!")
	}
}

func ClearRuntimeLogs() error {
	runtimeLogs := make([]Models.RuntimeLog, 0)
	for _, logNote := range GetRuntimeLogs() {
		//	Only fot current session.
		if logNote.SessionID != GetSessionId() {
			runtimeLogs = append(runtimeLogs, logNote)
		}
	}

	err := Configuration.Save(RuntimeLogFile, &runtimeLogs)
	if err != nil {
		return errors.New("Error while saving runtime log file! ")
	}
	return nil
}

func ClearAllLogs() error {
	err := ioutil.WriteFile(GetRuntimeLogFile(), []byte(``), 0644)
	if err != nil {
		return errors.New("Error while saving runtime log file! ")
	}
	return nil
}

func GetSessionId() string {
	return sessionID
}

func GetRuntimeLogNote() string {
	return runtimeLogNote
}

func ResetRuntimeLogNote() {
	runtimeLogNote = ""
}

func BuildRuntimeLogNote(logNote Models.RuntimeLog) string {
	runtimeLog := "[" + logNote.Time + "]" // + "\t"
	//runtimeLog += logNote.SessionID + "\t"
	runtimeLog += "[" + logNote.Level + "]"    // + "\t"
	runtimeLog += "[" + logNote.Category + "]" // + "\t"
	runtimeLog += " " + logNote.Message

	return SetLogLevel(logNote.Level, runtimeLog)
}

func SetLogLevel(level string, str string) string {
	switch level {
	case "info":
		return fmt.Sprintf("\x1b[89m%s\x1b[0m", str)
	case "trace":
		return fmt.Sprintf("\x1b[94m%s\x1b[0m", str)
	case "debug":
		return fmt.Sprintf("\x1b[92m%s\x1b[0m", str)
	case "warning":
		return fmt.Sprintf("\x1b[93m%s\x1b[0m", str)
	case "error":
		return fmt.Sprintf("\x1b[91m%s\x1b[0m", str)
	case "success":
		return fmt.Sprintf("\x1b[92m%s\x1b[0m", str)
	default:
		return fmt.Sprintf("\x1b[90m%s\x1b[0m", str)
	}
}

func CountErrorsInRuntimeLog() int {
	count := 0
	for _, logNote := range GetRuntimeLogs() {
		//	Only fot current session.
		if logNote.SessionID == GetSessionId() && logNote.Level == "error" {
			count++
		}
	}

	return count
}
