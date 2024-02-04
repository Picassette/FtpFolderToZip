package common

import (
	"fmt"
	"time"
)

/*
* Print message , level = info/warn/error/debug/critical
 */
func PrintMsg(msg string, level string) {
	var msgType string
	currentTime := time.Now().Format("2006-01-02T15:04:05")
	switch level {
	case "info":
		msgType = "INFO"
	case "warning":
		msgType = "WARNING"
	case "error":
		msgType = "ERROR"
	case "debug":
		msgType = "DEBUG"
	case "critical":
		msgType = "CRITICAL"
	default:
		fmt.Printf("%s | [CRITICAL] : INVALID MESSAGE LEVEL %s", currentTime, level)
		return
	}
	fmt.Printf("%s | [%s] : %s\n", msgType, currentTime, msg)
}
