package console

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
)

type (
	consoleLogger struct {
		// slLogLevel uses syslog's definitions which have higher priority
		// levels defined in descending order (0 is highest)
		slLogLevel syslog.Priority
	}
)

// NewConsoleLogger creates a default logger object that prints log messages
// to the console.
func NewConsoleLogger() *consoleLogger {
	return &consoleLogger{
		slLogLevel: syslog.LOG_DEBUG,
	}
}

// NewConsoleLoggerExt() allows the user to create a customized logger
// that will then be used by the vlogger interface.
func NewConsoleLoggerExt(prefix string, flags int) *consoleLogger {
	log.SetPrefix(prefix)
	log.SetFlags(flags)
	return NewConsoleLogger()
}

func (cl *consoleLogger) Debug(msg string) {
	if cl.slLogLevel >= syslog.LOG_DEBUG {
		log.Println("[DEBUG]", msg)
	}
}

func (cl *consoleLogger) Debugf(format string, params ...interface{}) {
	if cl.slLogLevel >= syslog.LOG_DEBUG {
		msg := fmt.Sprintf(format, params...)
		log.Println("[DEBUG]", msg)
	}
}

func (cl *consoleLogger) Info(msg string) {
	if cl.slLogLevel >= syslog.LOG_INFO {
		toSTDOUT(msg)
	}
}

func (cl *consoleLogger) Infof(format string, params ...interface{}) {
	if cl.slLogLevel >= syslog.LOG_INFO {
		msg := fmt.Sprintf(format, params...)
		toSTDOUT(msg)
	}
}

func toSTDOUT(msg string) {
	log.SetOutput(os.Stdout)
	log.Println("[INFO]", msg)
	log.SetOutput(os.Stderr)
}

func (cl *consoleLogger) Warning(msg string) {
	if cl.slLogLevel >= syslog.LOG_WARNING {
		log.Println("[WARNING]", msg)
	}
}

func (cl *consoleLogger) Warningf(format string, params ...interface{}) {
	if cl.slLogLevel >= syslog.LOG_WARNING {
		msg := fmt.Sprintf(format, params...)
		log.Println("[WARNING]", msg)
	}
}

func (cl *consoleLogger) Error(msg string) {
	if cl.slLogLevel >= syslog.LOG_ERR {
		log.Println("[ERROR]", msg)
	}
}

func (cl *consoleLogger) Errorf(format string, params ...interface{}) {
	if cl.slLogLevel >= syslog.LOG_ERR {
		msg := fmt.Sprintf(format, params...)
		log.Println("[ERROR]", msg)
	}
}

func (cl *consoleLogger) Critical(msg string) {
	if cl.slLogLevel >= syslog.LOG_CRIT {
		log.Println("[CRITICAL]", msg)
	}
}

func (cl *consoleLogger) Criticalf(format string, params ...interface{}) {
	if cl.slLogLevel >= syslog.LOG_CRIT {
		msg := fmt.Sprintf(format, params...)
		log.Println("[CRITICAL]", msg)
	}
}

func (cl *consoleLogger) SetLogLevel(slLogLevel syslog.Priority) {
	cl.slLogLevel = slLogLevel
}

func (cl *consoleLogger) GetLogLevel() syslog.Priority {
	return cl.slLogLevel
}

func (cl *consoleLogger) Close() {
}
