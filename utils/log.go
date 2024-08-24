package utils

import (
	"log"
	"os"
)

func InitLog(logFile string) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
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
