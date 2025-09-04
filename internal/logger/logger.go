package logger

import (
	"fmt"
	"time"
)

func Info(msg string, args ...any) {
	fmt.Printf("%s [INFO] %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(msg, args...))
}

func Warn(msg string, args ...any) {
	fmt.Printf("%s [WARN] %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(msg, args...))
}

func Error(msg string, args ...any) {
	fmt.Printf("%s [ERROR] %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(msg, args...))
}
