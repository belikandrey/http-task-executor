package logger

import (
	"http-task-executor/task-service/internal/task-service/config"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger represents application logger.
type Logger interface {
	// Debug - print log with debug level
	Debug(args ...interface{})
	// Debugf - print formatted log with debug level
	Debugf(template string, args ...interface{})
	// Info - print formatted log with info level
	Info(args ...interface{})
	// Infof - print formatted log with info level
	Infof(template string, args ...interface{})
	// Warn - print log with warn level
	Warn(args ...interface{})
	// Warnf - print formatted log with warn level
	Warnf(template string, args ...interface{})
	// Error - print log with error level
	Error(args ...interface{})
	// Errorf - print formatted log with error level
	Errorf(template string, args ...interface{})
	// DPanic - print log with error level and panic
	DPanic(args ...interface{})
	// DPanicf - print formatted log with error level and panic
	DPanicf(template string, args ...interface{})
	// Fatal - print log with error level and os.Exit
	Fatal(args ...interface{})
	// Fatalf - print formatted log with error level and os.Exit
	Fatalf(template string, args ...interface{})
}

// NewLogger creates new instance of Logger.
func NewLogger(config *config.Config) (Logger, error) {

	file, err := os.OpenFile(config.LoggerConfig.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	logWriter, encoderCfg := getEnvBasedOptions(config.Env, file)

	var encoder zapcore.Encoder

	if config.LoggerConfig.Format != "json" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(getLoggerLevel(config.LoggerConfig.Level)))

	logger := zap.New(core, zap.AddCaller())

	sugar := logger.Sugar()

	if err := sugar.Sync(); err != nil {
		return nil, err
	}

	return sugar, nil
}

func getEnvBasedOptions(env string, file *os.File) (zapcore.WriteSyncer, zapcore.EncoderConfig) {
	var logWriter zapcore.WriteSyncer
	var encoderCfg zapcore.EncoderConfig
	switch env {
	case "local":
		logWriter = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(file))
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	case "prod":
		logWriter = zapcore.AddSync(file)
		encoderCfg = zap.NewProductionEncoderConfig()
	default:
		panic("unknown env type: " + env)
	}
	return logWriter, encoderCfg
}

func getLoggerLevel(lvl string) zapcore.Level {
	switch lvl {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "fatal":
		return zap.FatalLevel
	case "panic":
		return zap.PanicLevel
	default:
		return zap.InfoLevel
	}
}
