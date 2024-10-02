package logger

import (
	"io"
	"log"
	"os"
)

const DEFAULT_LOG_FILE = "/tmp/project00.log"

func InitLog(logFileName string) {
	logFile, err := os.OpenFile(
		logFileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Fatalf("[FATAL] Failed to open log file: %v", err)
	}
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	Info("Successfully initiated log")
}

func Debug(msg string) {
	log.Printf("[DEBUG] %s", msg)
}

func Info(msg string) {
	log.Printf("[INFO] %s", msg)
}

func Warn(msg string) {
	log.Printf("[WARN] %s", msg)
}

func Error(msg string) {
	log.Printf("[ERROR] %s", msg)
}
