package logger

import (
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
}

func getEncoder(isJson bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time" // To convert Timestamp to Time
	if isJson {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case Info:
		return zapcore.InfoLevel
	case Warn:
		return zapcore.WarnLevel
	case Debug:
		return zapcore.DebugLevel
	case Error:
		return zapcore.ErrorLevel
	case Fatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func NewZapLogger(config Configuration) (Logger, error) {
	cores := []zapcore.Core{}

	level := getZapLevel(config.ConsoleLevel)
	writer := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(getEncoder(config.ConsoleJSONFormat), writer, level)
	cores = append(cores, core)

	combinedCore := zapcore.NewTee(cores...)
	zc := apmzap.Core{
		Tracer:             apm.DefaultTracer,
		FatalFlushTimeout: 0,
	}

	zc.WrapCore(combinedCore)

	logger := zap.New(combinedCore, zap.AddCallerSkip(2), zap.AddCaller()).Sugar()
	return &zapLogger{sugaredLogger: logger}, nil
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format, args...)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugaredLogger.Warnf(format, args...)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugaredLogger.Errorf(format, args...)
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format, args...)
}

func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format, args...)
}

func (l *zapLogger) With(args ...interface{}) *zap.SugaredLogger {
	return l.sugaredLogger.With(args)
}

func (l *zapLogger) WithFields(fields Fields) Logger {
	var f = make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}
	newLogger := l.sugaredLogger.With(f...)
	return &zapLogger{newLogger}
}
