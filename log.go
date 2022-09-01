package THz

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (thz *THz) SetLog(log *zap.Logger) { thz.log = log }
func (thz *THz) SetZapLog(level zapcore.Level) {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: true,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			TimeKey:        "ts",
			NameKey:        "logger",
			CallerKey:      "caller",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	log, err := config.Build()
	if err != nil {
		panic(err)
	}
	thz.log = log
}
