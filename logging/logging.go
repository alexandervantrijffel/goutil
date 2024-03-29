package logging

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/TheZeroSlave/zapsentry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const LOGDEBUG = "DEBUG: "
const LOGINFO = "INFO: "
const LOGWARNING = "WARN: "
const LOGERROR = "ERROR: "
const LOGFATAL = "FATAL: "

var logger *zap.Logger

var appName string
var dbg bool
var loggerInitialized bool

func init() {
	InitWith("unset", true, "")
}

func InitWith(myAppName string, debugMode bool, sentryDsn ...string) {
	dbg = debugMode
	appName = myAppName
	var cfg zap.Config
	if dbg {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig = zapcore.EncoderConfig{
			TimeKey:        "", // time disabled
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "", // caller disabled
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
	} else {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	}
	cfg.OutputPaths = []string{"stdout"}
	cfg.EncoderConfig.EncodeCaller = loggerCallerEntryResolver
	cfg.Encoding = "json" // not console
	var err error
	logger, err = cfg.Build()
	if err != nil {
		log.Fatalf("Failed to build zap logger! %+v", err)
		return
	}
	if !debugMode && len(sentryDsn) > 0{
		initSentry(sentryDsn[0])
	}
	loggerInitialized = true
}

func initSentry(dsn string) {
	level := zapcore.ErrorLevel
	cfg := zapsentry.Configuration{
		Level: level,
		Tags: map[string]string{
			"component": "system",
		},
	}
	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromDSN(dsn))
	//in case of err it will return noop core. so we can safely attach it
	if err != nil {
		logger.Warn("failed to init zap")
	}
	logger = zapsentry.AttachCoreToLogger(core, logger)
	logger.Info("sentry forwarding enabled", zap.Stringer("level", level))
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
	msg := fmt.Sprint(v...)
	fields := getDefaultFields()
	switch prefix {
	case LOGDEBUG:
		if !loggerInitialized {
			fmt.Printf(prefix+" (logger not initialized)\n", v...)
			return
		}
		logger.Debug(msg, fields...)
	case LOGINFO:
		if !loggerInitialized {
			fmt.Printf(prefix+" (logger not initialized)\n", v...)
			return
		}
		logger.Info(msg, fields...)
	case LOGWARNING:
		if !loggerInitialized {
			fmt.Printf(prefix+" (logger not initialized)\n", v...)
			return
		}
		logger.Warn(msg, fields...)
	case LOGERROR:
		if !loggerInitialized {
			fmt.Printf(prefix+" (logger not initialized)\n", v...)
			return
		}
		logger.Error(msg, fields...)
	case LOGFATAL:
		if !loggerInitialized {
			fmt.Printf(prefix+" (logger not initialized)\n", v...)
			return
		}
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
	if loggerInitialized {
		flushThisLog(logger)
	}
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
	enc.AppendString(getLoggingCaller(13))
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
