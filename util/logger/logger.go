package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	LogDirectoryPath = "logs"
)

type Logger struct {
	LogDirectory string
}

func New() *Logger {
	err := os.Mkdir(LogDirectoryPath, 0755)
	if err != nil {
		return nil
	}

	return &Logger{
		LogDirectory: LogDirectoryPath,
	}
}

func (l *Logger) Info() *log.Logger {
	file := readOrCreateLogFile()

	return log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (l *Logger) Warning() *log.Logger {
	file := readOrCreateLogFile()

	return log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (l *Logger) Error() *log.Logger {
	file := readOrCreateLogFile()

	return log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func readOrCreateLogFile() *os.File {
	year, month, day := time.Now().Date()
	fileName := fmt.Sprintf("%v-%v-%v.log", day, month.String(), year)
	file, _ := os.OpenFile(LogDirectoryPath+"/"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	return file
}
