package logger

import (
	"errors"
	"fmt"
	"github.com/Shitovdm/GitRsync/src/component/conf"
	"github.com/Shitovdm/GitRsync/src/component/helper"
	"github.com/Shitovdm/GitRsync/src/model"
	"io/ioutil"
	"time"
)

const timeFormatStr = "2006-01-02 15:04:05"

// Describes log file name.
var RuntimeLogFile = "RuntimeLogs.json"
var sessionID = ""
var runtimeLogNote = ""

// Init initializes runtime log for new session.
func Init() {
	sessionID = helper.GenerateUUID()
	runtimeLogNote = ""
}

// GetRuntimeLogFile gets runtime log file.
func GetRuntimeLogFile() string {
	return conf.BuildPlatformPath(RuntimeLogFile)
}

// GetRuntimeLogs gets runtime logs.
func GetRuntimeLogs() []model.RuntimeLog {

	runtimeLogs := make([]model.RuntimeLog, 0)
	err := conf.Load(RuntimeLogFile, &runtimeLogs)
	if err != nil {
		err = conf.Save("Runtime.json", []map[string]interface{}{})
		if err != nil {
			panic(fmt.Sprintf("Error while creating new runtime log file! %s", err.Error()))
		}
		return []model.RuntimeLog{}
	}

	return runtimeLogs
}

// AddRuntimeLog adds runtime log.
func AddRuntimeLog(sessionID string, level string, category string, message string) {

	runtimeLog := model.RuntimeLog{
		SessionID: sessionID,
		Time:      time.Now().Format(timeFormatStr),
		Level:     level,
		Category:  category,
		Message:   message,
	}

	runtimeLogNote = BuildRuntimeLogNote(runtimeLog)
	runtimeLogs := GetRuntimeLogs()
	runtimeLogs = append(runtimeLogs, runtimeLog)
	err := conf.Save(RuntimeLogFile, &runtimeLogs)
	if err != nil {
		panic("Error while saving runtime log!")
	}
}

// ClearRuntimeLogs removes only runtime logs for current session.
func ClearRuntimeLogs() error {

	runtimeLogs := make([]model.RuntimeLog, 0)
	for _, logNote := range GetRuntimeLogs() {
		//	Only fot current session.
		if logNote.SessionID != GetSessionID() {
			runtimeLogs = append(runtimeLogs, logNote)
		}
	}

	err := conf.Save(RuntimeLogFile, &runtimeLogs)
	if err != nil {
		return errors.New("Error while saving runtime log file! ")
	}
	return nil
}

// ClearAllLogs removes all runtime logs.
func ClearAllLogs() error {

	err := ioutil.WriteFile(GetRuntimeLogFile(), []byte(``), 0644)
	if err != nil {
		return errors.New("Error while saving runtime log file! ")
	}
	return nil
}

// GetSessionId gets session ID.
func GetSessionID() string {
	return sessionID
}

// GetRuntimeLogNote gets runtime log note.
func GetRuntimeLogNote() string {
	return runtimeLogNote
}

// ResetRuntimeLogNote resets runtime log note.
func ResetRuntimeLogNote() {
	runtimeLogNote = ""
}

// BuildRuntimeLogNote returns runtime log note.
func BuildRuntimeLogNote(logNote model.RuntimeLog) string {

	runtimeLog := "[" + logNote.Time + "]"
	runtimeLog += "[" + logNote.Level + "]"
	runtimeLog += "[" + logNote.Category + "]"
	runtimeLog += " " + logNote.Message

	return SetLogLevel(logNote.Level, runtimeLog)
}

// SetLogLevel returns different colors by log category.
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

// CountErrorsInRuntimeLog returns count of errors in runtime log.
func CountErrorsInRuntimeLog() int {

	count := 0
	for _, logNote := range GetRuntimeLogs() {
		if logNote.SessionID == GetSessionID() && logNote.Level == "error" {
			count++
		}
	}

	return count
}
