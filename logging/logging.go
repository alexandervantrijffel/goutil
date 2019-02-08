package logging

import (
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const LOGDEBUG = "DEBUG: "
const LOGINFO = "INFO: "
const LOGWARNING = "WARN: "
const LOGERROR = "ERROR: "
const LOGFATAL = "FATAL: "

var logger *zap.Logger
var initOnce sync.Once

var appName string
var dbg bool

func InitWith(myAppName string, debugMode bool) {
	dbg = debugMode
	appName = myAppName
	initOnce.Do(func() {
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		cfg.EncoderConfig.EncodeCaller = loggerCallerEntryResolver
		if !dbg {
			cfg = zap.NewProductionConfig()
		}
		var err error
		logger, err = cfg.Build()
		if err != nil {
			fmt.Printf("Failed to build zap logger! %+v", err)
		}
	})
}

func Debugf(format string, v ...interface{}) {
	logIt(LOGDEBUG, format, v...)
}

func Debug(v ...interface{}) {
	logItNoFormat(LOGDEBUG, v...)
}

func Infof(format string, v ...interface{}) {
	logIt(LOGINFO, format, v...)
}

func Info(v ...interface{}) {
	logItNoFormat(LOGINFO, v...)
}

func Warningf(format string, v ...interface{}) {
	logIt(LOGWARNING, format, v...)
}
func Warning(v ...interface{}) {
	logItNoFormat(LOGWARNING, v...)
}
func Errorf(format string, v ...interface{}) {
	logIt(LOGERROR, format, v...)
}
func Error(v ...interface{}) {
	logItNoFormat(LOGERROR, v...)
}
func Errore(err error) {
	logIt(LOGERROR, err.Error())
}

func Fatal(v ...interface{}) {
	logItNoFormat(LOGFATAL, v...)
}

func logItNoFormat(prefix string, v ...interface{}) {
	if len(appName) == 0 {
		panic("Logging not initialized, please call logging.InitWith() first.")
	}
	msg := fmt.Sprint(v...)
	fields := getDefaultFields()
	switch prefix {
	case LOGDEBUG:
		logger.Debug(msg, fields...)
	case LOGINFO:
		logger.Info(msg, fields...)
	case LOGWARNING:
		logger.Warn(msg, fields...)
	case LOGERROR:
		logger.Error(msg, fields...)
	case LOGFATAL:
		logger.Fatal(msg, fields...)
	}
}

func logIt(prefix string, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	logItNoFormat(prefix, msg)
}

func getDefaultFields() (fields []zap.Field) {
	if !dbg {
		fields = append(fields, zap.String("appname", appName))
	}
	return
}

func Flush() {
	flushThisLog(logger)
}

func flushThisLog(l *zap.Logger) {
	err := l.Sync()
	if err != nil {
		log.Print("Failed to sync logger", err)
	}
}

func getLoggingCaller(from int) string {
	var f string
	var l int
	for {
		prevL := -1
		var ok bool
		_, f, l, ok = runtime.Caller(from)
		if !ok {
			f = "?"
			l = -1
		} else {
			f = TrimmedPath(f)
			if (strings.HasPrefix(f, "logging/") || strings.HasPrefix(f, "errorcheck/")) && prevL != l {
				from += 1
				continue
			}
		}
		break
	}
	return fmt.Sprintf("%s:%d", f, l)
}

func loggerCallerEntryResolver(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(getLoggingCaller(8))
}

func getStack() string {
	return " -- " + fmt.Sprintf("%s", debug.Stack())
}

// TrimmedPath returns a package/file:line description of the caller,
// preserving only the leaf directory name and file name.
func TrimmedPath(fullPath string) string {
	// nb. To make sure we trim the path correctly on Windows too, we
	// counter-intuitively need to use '/' and *not* os.PathSeparator here,
	// because the path given originates from Go stdlib, specifically
	// runtime.Caller() which (as of Mar/17) returns forward slashes even on
	// Windows.
	//
	// See https://github.com/golang/go/issues/3335
	// and https://github.com/golang/go/issues/18151
	//
	// for discussion on the issue on Go side.
	//
	// Find the last separator.
	//
	idx := strings.LastIndexByte(fullPath, '/')
	if idx == -1 {
		return fullPath
	}
	// Find the penultimate separator.
	idx = strings.LastIndexByte(fullPath[:idx], '/')
	if idx == -1 {
		return fullPath
	}
	return fullPath[idx+1:]
}
