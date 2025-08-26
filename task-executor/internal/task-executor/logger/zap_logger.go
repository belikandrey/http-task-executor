package logger

import (
	"os"
	"task-executor/internal/task-executor/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

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
