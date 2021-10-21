package logger

import (
	"cloud.google.com/go/logging"
	"fmt"
	"time"
)

const logId = "playground-log"

type CloudLoggingHandler struct {
	logger *logging.Logger
	client *logging.Client
}

// NewCloudLoggingHandler creates CloudLoggingHandler
func NewCloudLoggingHandler(client *logging.Client) *CloudLoggingHandler {
	return &CloudLoggingHandler{client: client, logger: client.Logger(logId)}
}

func (c CloudLoggingHandler) Info(args ...interface{}) {
	c.logMessage(logging.Info, args...)
}

func (c CloudLoggingHandler) Infof(format string, args ...interface{}) {
	c.logMessage(logging.Info, fmt.Sprintf(format, args...))
}

func (c CloudLoggingHandler) Warn(args ...interface{}) {
	c.logMessage(logging.Warning, args...)
}

func (c CloudLoggingHandler) Warnf(format string, args ...interface{}) {
	c.logMessage(logging.Warning, fmt.Sprintf(format, args...))
}

func (c CloudLoggingHandler) Error(args ...interface{}) {
	c.logMessage(logging.Error, args...)
}

func (c CloudLoggingHandler) Errorf(format string, args ...interface{}) {
	c.logMessage(logging.Error, fmt.Sprintf(format, args...))
}

func (c CloudLoggingHandler) Debug(args ...interface{}) {
	c.logMessage(logging.Debug, args...)
}

func (c CloudLoggingHandler) Debugf(format string, args ...interface{}) {
	c.logMessage(logging.Debug, fmt.Sprintf(format, args...))
}

func (c CloudLoggingHandler) Fatal(args ...interface{}) {
	c.logMessage(logging.Critical, args...)
}

func (c CloudLoggingHandler) Fatalf(format string, args ...interface{}) {
	c.logMessage(logging.Critical, fmt.Sprintf(format, args...))
}

// logMessage buffers the Entry for output to the logging service.
func (c CloudLoggingHandler) logMessage(severity logging.Severity, args ...interface{}) {
	c.logger.Log(logging.Entry{
		Timestamp: time.Now(),
		Severity:  severity,
		Payload:   fmt.Sprint(args...),
	})
}

// CloseConn waits for all opened loggers to be flushed and closes the client.
func (c CloudLoggingHandler) CloseConn() error {
	return c.client.Close()
}
