package logger

import (
	"io"
	"log"
	"os"
)

// Global logger instance
var Log *log.Logger

// Init initializes the logger
func Init(logFilePath string, alsoConsole bool) error {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	var writer io.Writer = file
	if alsoConsole {
		// write to both file and console
		writer = io.MultiWriter(os.Stdout, file)
	}

	// Create a logger with timestamp prefix
	Log = log.New(writer, "", log.LstdFlags)
	return nil

}

func Info(v ...interface{}) {
	Log.SetPrefix("[INFO] ")
	Log.Println(v...)
}

func Warn(v ...interface{}) {
	Log.SetPrefix("[WARN] ")
	Log.Println(v...)
}

func Error(v ...interface{}) {
	Log.SetPrefix("[ERROR] ")
	Log.Println(v...)
}
