package logs

import (
	"fmt"
	"log"
	"os"
)

const (
	info = "[INFO]    "
	warn = "[WARNING] "
	errs = "[ERROR]   "

	fileLogs = "./internal/logs/logs.log" // where save log
)

var (
	// WarningLogger for log WARNING
	WarningLogger *log.Logger
	// InfoLogger for log INFO
	InfoLogger *log.Logger
	// ErrorLogger for lof ERROR
	ErrorLogger *log.Logger
)

// InitLoggers to project
func InitLoggers() {
	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile(fileLogs, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("CANT START LOGGER reason: %s", err.Error())
	}
	InfoLogger = log.New(file, info, log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, warn, log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, errs, log.Ldate|log.Ltime|log.Lshortfile)

	fmt.Fprintf(file, "----------------------------------------\n")
}

// Succes job in INFO log with comment
func Succes(format string, v ...interface{}) {
	InfoLogger.Printf("[SUCCES] %s", fmt.Sprintf(format, v...))
}

// Mess job in INFO log with comment
func Mess(format string, v ...interface{}) {
	InfoLogger.Printf("[MESS] %s", fmt.Sprintf(format, v...))
}

// DB job in INFO log with comment
func DB(format string, v ...interface{}) {
	InfoLogger.Printf("[DB] %s", fmt.Sprintf(format, v...))
}

// Err job in ERROR log with comment
func Err(format string, v ...interface{}) {
	ErrorLogger.Printf("%s", fmt.Sprintf(format, v...))
}

// DBErr job in ERROR log with comment
func DBErr(format string, v ...interface{}) {
	ErrorLogger.Printf("[DBERR] %s", fmt.Sprintf(format, v...))
}

// Warn job in WARN log with comment
func Warn(format string, v ...interface{}) {
	WarningLogger.Printf("%s", fmt.Sprintf(format, v...))
}
