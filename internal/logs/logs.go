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

	suc    = "[SUCCES] "
	mes    = "[MESS]   "
	db     = "[DB]     "
	dberr  = "[DBERR]  "
	dbwarn = "[DBWARN] "

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

	// file, err := os.OpenFile(fileLogs, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	// if err != nil {
	// 	log.Fatalf("CANT START LOGGER reason: %s", err.Error())
	// }

	InfoLogger = log.New(os.Stdout, info, log.Ldate|log.Ltime)
	WarningLogger = log.New(os.Stdout, warn, log.Ldate|log.Ltime)
	ErrorLogger = log.New(os.Stdout, errs, log.Ldate|log.Ltime)

	fmt.Fprintf(os.Stdout, "----------------------------------------\n")
}

//! INFO

// Succes job in INFO log with comment
func Succes(format string, v ...interface{}) {
	InfoLogger.Printf("%s|| %s", suc, fmt.Sprintf(format, v...))
}

// Mess job in INFO log with comment
func Mess(format string, v ...interface{}) {
	InfoLogger.Printf("%s|| %s", mes, fmt.Sprintf(format, v...))
}

// DB job in INFO log with comment
func DB(format string, v ...interface{}) {
	InfoLogger.Printf("%s|| %s", db, fmt.Sprintf(format, v...))
}

//! ERR

// Err job in ERROR log with comment
func Err(format string, v ...interface{}) {
	ErrorLogger.Printf("%s", fmt.Sprintf(format, v...))
}

// DBErr db job in ERROR log with comment
func DBErr(format string, v ...interface{}) {
	ErrorLogger.Printf("%s|| %s", dberr, fmt.Sprintf(format, v...))
}

//! WARN

// Warn job in WARN log with comment
func Warn(format string, v ...interface{}) {
	WarningLogger.Printf("%s", fmt.Sprintf(format, v...))
}

// WarnNot when start not
func WarnNot(t string) {
	WarningLogger.Printf("[WARNING]|| %s at %s", "start notice", t)
}

// DBWarn db job in WARN log with comment
func DBWarn(format string, v ...interface{}) {
	WarningLogger.Printf("%s|| %s", dbwarn, fmt.Sprintf(format, v...))
}
