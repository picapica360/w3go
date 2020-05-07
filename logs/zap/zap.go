package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/picapica360/w3go/config/env"
	"github.com/picapica360/w3go/logs"
)

// DefaultEncoderConfig default profile for zapcore.EncoderConfig.
var DefaultEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "msg",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
	EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
}

// DefaultLevelEnablerFunc implement zapcore.LevelEnabler.
//  the minimal level for development is Debug, and others is Info.
func DefaultLevelEnablerFunc() zap.LevelEnablerFunc {
	return func(level zapcore.Level) bool {
		if env.IsDevelopment() {
			return level >= zapcore.DebugLevel
		}
		return level >= zapcore.InfoLevel
	}
}

var adapters []func() zapcore.Core

// Register logger adapter.
func Register(adapter ...func() zapcore.Core) {
	if adapter == nil {
		panic("[logger]: Register adapter is nil")
	}

	adapters = append(adapters, adapter...)
}

// Build  build the logger.
func Build() logs.Logger {
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	}
	ads := make([]zapcore.Core, len(adapters))
	for i, adapter := range adapters {
		ads[i] = adapter()
	}
	core := zapcore.NewTee(ads...)
	return zap.New(core, options...).Sugar()
}
