package log

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path/filepath"
	"runtime"
	"wxch-dashboard/config"
)

var logger *zap.Logger

func init() {
	var err error

	sentryHook := zap.Hooks(func(entry zapcore.Entry) error {
		if entry.Level == zapcore.ErrorLevel {
			sentry.CaptureMessage(fmt.Sprintf("%s, Line No: %d :: %s", entry.Caller.File, entry.Caller.Line, entry.Message))
		}
		return nil
	})

	if config.Get().Debug.Verbose {
		developmentConfig := zap.NewDevelopmentConfig()

		if _, file, _, ok := runtime.Caller(0); ok {
			basePath := filepath.Dir(filepath.Dir(filepath.Dir(file))) + "/"
			developmentConfig.EncoderConfig.EncodeCaller = func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
				rel, err := filepath.Rel(basePath, caller.File)
				if err != nil {
					encoder.AppendString(caller.FullPath())
				} else {
					encoder.AppendString(fmt.Sprintf("%s:%d", rel, caller.Line))
				}
			}
		}

		logger, err = developmentConfig.Build(sentryHook)
		if err != nil {
			panic(err)
		}
	} else {
		logger, err = zap.NewProduction(sentryHook)
		if err != nil {
			panic(err)
		}
	}
}

func GetLogger() *zap.Logger {
	return logger
}
