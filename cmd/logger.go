package main

import (
	"fmt"
	"log"
	"os"
)

type severity int

// Severity levels.
const (
	sDebug severity = iota
	sInfo
	sWarning
	sError
	sFatal
)

//logInterface logger
type logInterface struct {
	LogDebug *log.Logger
	LogInfo  *log.Logger
	LogWarn  *log.Logger
	LogError *log.Logger
	LogFatal *log.Logger
	loglevel severity
	Output   *os.File
}

//Info print with info level
func (l *logInterface) Info(v ...interface{}) {
	if l.loglevel <= sInfo {
		l.output(sInfo, 0, fmt.Sprint(v...))
	} else {
		return
	}
}

//Infof printf with info level
func (l *logInterface) Infof(format string, v ...interface{}) {
	if l.loglevel <= sInfo {
		l.output(sInfo, 0, fmt.Sprintf(format, v...))
	} else {
		return
	}
}

//Debug print with debug level
func (l *logInterface) Debug(v ...interface{}) {
	if l.loglevel <= sDebug {
		l.output(sDebug, 0, fmt.Sprint(v...))
	} else {
		return
	}
}

//Debugf printf with debug level
func (l *logInterface) Debugf(format string, v ...interface{}) {
	if l.loglevel <= sDebug {
		l.output(sDebug, 0, fmt.Sprintf(format, v...))
	} else {
		return
	}
}

//Error print with error level
func (l *logInterface) Error(v ...interface{}) {
	if l.loglevel <= sError {
		l.output(sError, 0, fmt.Sprint(v...))
	} else {
		return
	}
}

//Errorf printf with error level
func (l *logInterface) Errorf(format string, v ...interface{}) {
	if l.loglevel <= sError {
		l.output(sError, 0, fmt.Sprintf(format, v...))
	} else {
		return
	}
}

//Warn print with warn level
func (l *logInterface) Warn(v ...interface{}) {
	if l.loglevel <= sWarning {
		l.output(sWarning, 0, fmt.Sprint(v...))
	} else {
		return
	}
}

//Warnf print with warnf level
func (l *logInterface) Warnf(format string, v ...interface{}) {
	if l.loglevel <= sWarning {
		l.output(sWarning, 0, fmt.Sprintf(format, v...))
	} else {
		return
	}
}

//Fatal print with fatal level
func (l *logInterface) Fatal(v ...interface{}) {
	l.output(sFatal, 0, fmt.Sprint(v...))
	os.Exit(1)
}

//Fatalf printf with fatal level
func (l *logInterface) Fatalf(format string, v ...interface{}) {
	l.output(sFatal, 0, fmt.Sprintf(format, v...))
	os.Exit(1)
}

//output print result
func (l *logInterface) output(s severity, depth int, txt string) {
	switch s {
	case sDebug:
		l.LogDebug.Output(3, txt)
	case sInfo:
		l.LogInfo.Output(3, txt)
	case sWarning:
		l.LogWarn.Output(3, txt)
	case sError:
		l.LogError.Output(3, txt)
	case sFatal:
		l.LogFatal.Output(3, txt)
	default:
		panic(fmt.Sprintln("unrecognized severity:", s))
	}
}

//initLogger initialize logger
func initLogger(ll, ff *string) *logInterface {

	lg := &logInterface{}

	switch *ll {
	case "debug":
		lg.loglevel = sDebug
	case "info":
		lg.loglevel = sInfo
	case "warning":
		lg.loglevel = sWarning
	case "error":
		lg.loglevel = sError
	case "fatal":
		lg.loglevel = sFatal

	}

	var (
		out *os.File
		f   *os.File
		err error
	)
	// if log to file
	if *ff != "" {
		if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
			os.Mkdir(logFilePath, 0755)
		}
		f, err = os.OpenFile(logFilePath+"/"+*ff, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		out = f
	} else {
		out = os.Stderr
	}

	lg.Output = out
	lg.LogDebug = log.New(out, "DEBUG: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	lg.LogInfo = log.New(out, "INFO: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	lg.LogWarn = log.New(out, "WARN: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	lg.LogError = log.New(out, "ERROR: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	lg.LogFatal = log.New(out, "FATAL: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	return lg

}
