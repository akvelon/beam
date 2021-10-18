package logger

import (
	"cloud.google.com/go/logging"
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	INFO  string = "[INFO]: "
	WARN  string = "[WARN]: "
	ERROR string = "[ERROR]: "
	FATAL string = "[FATAL]: "
	DEBUG string = "[DEBUG]: "
)

var logger *logging.Logger
var client *logging.Client

func init() {
	initLogger()
}

// initLogger determines which logger to use by the environment variable
func initLogger() {
	switch os.Getenv("LOGGER") {
	case "cloudLogger":
		initCloudLogger()
	default:
		logger = nil
	}
}

// initCloudLogger initializes the Google Cloud client and logger
func initCloudLogger() {
	cl, err := logging.NewClient(context.Background(), os.Getenv("PROJECT_ID"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	client = cl
	logger = client.Logger("playground-log")
}

// Info logs a message at level Info.
func Info(args ...interface{}) {
	logMessage(INFO, logging.Info, args...)
}

func Infof(format string, args ...interface{}) {
	logMessageF(INFO, logging.Info, format, args...)
}

// Warn logs a message at level Warning.
func Warn(args ...interface{}) {
	logMessage(WARN, logging.Warning, args...)
}

func Warnf(format string, args ...interface{}) {
	logMessageF(WARN, logging.Warning, format, args...)
}

// Error logs a message at level Error.
func Error(args ...interface{}) {
	logMessage(ERROR, logging.Error, args...)
}

func Errorf(format string, args ...interface{}) {
	logMessageF(ERROR, logging.Error, format, args...)
}

// Debug logs a message at level Debug.
func Debug(args ...interface{}) {
	logMessage(DEBUG, logging.Debug, args...)
}

func Debugf(format string, args ...interface{}) {
	logMessageF(DEBUG, logging.Debug, format, args...)
}

// Fatal logs a message at level Fatal.
// Then the process will exit with status set to 1 if you not use GC loggers.
func Fatal(args ...interface{}) {
	if logger == nil {
		args := append([]interface{}{FATAL}, args...)
		log.Fatal(args...)
	} else {
		logger.Log(logging.Entry{
			Timestamp: time.Now(),
			Severity:  logging.Critical,
			Payload:   fmt.Sprint(args...),
		})
	}
}

func Fatalf(format string, args ...interface{}) {
	if logger == nil {
		log.Fatalf(FATAL+format, args...)
	} else {
		logger.Log(logging.Entry{
			Timestamp: time.Now(),
			Severity:  logging.Critical,
			Payload:   fmt.Sprintf(format, args...),
		})
	}
}

// CloseConn waits for all opened loggers to be flushed and closes the client if you use GC loggers.
func CloseConn() {
	if logger != nil {
		if err := client.Close(); err != nil {
			log.Fatalf("Failed to close client: %v", err)
		}
	}
}

// logMessage logs a message at level strSev if you not use GC loggers.
// If you use GC logger: buffers the Entry with your message and severity for output to the AppEngine logging service.
func logMessage(strSev string, gcpSev logging.Severity, args ...interface{}) {
	if logger == nil {
		args := append([]interface{}{strSev}, args...)
		log.Print(args...)
	} else {
		logger.Log(logging.Entry{
			Timestamp: time.Now(),
			Severity:  gcpSev,
			Payload:   fmt.Sprint(args...),
		})
	}
}

// logMessage logs a message according to a format specifier at level strSev if you not use GC loggers.
// If you use GC logger: buffers the Entry with your message and severity for output to the AppEngine logging service.
func logMessageF(strSev string, gcpSev logging.Severity, format string, args ...interface{}) {
	if logger == nil {
		log.Printf(strSev+format, args...)
	} else {
		logger.Log(logging.Entry{
			Timestamp: time.Now(),
			Severity:  gcpSev,
			Payload:   fmt.Sprintf(format, args...),
		})
	}
}
