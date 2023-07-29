package loggers

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var (
	DefaultConsoleWriter = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	EnableConsolColor  = false
	EnableConsolWriter = true
	EnableFileWriter   = true
	globalConfig       *Config

	NullLogger      *zerolog.Logger
	CoreDebugLogger *zerolog.Logger
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logger := zerolog.New(&NullWriter{})
	NullLogger = &logger
	CoreDebugLogger = NullLogger

	globalConfig = NewConfig()
}

func SetGlobalConfig(c Config) *Config {
	globalConfig = &c
	EnableConsolColor = c.EnableConsolColor
	EnableConsolWriter = c.EnableConsolWriter
	EnableFileWriter = c.EnableFileWriter
	level, err := zerolog.ParseLevel(c.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	return globalConfig
}

func GetLogger(fn string) (*zerolog.Logger, error) {
	writers := []io.Writer{&NullWriter{}}
	if !EnableConsolColor {
		DefaultConsoleWriter.NoColor = true
	}
	if EnableConsolWriter {
		writers = append(writers, DefaultConsoleWriter)
	}
	if EnableFileWriter && len(fn) > 0 {
		writers = append(writers, &lumberjack.Logger{
			Filename:   fn,
			MaxSize:    globalConfig.MaxSizeMB,
			MaxAge:     globalConfig.MaxAgeDays,
			MaxBackups: globalConfig.MaxBackups,
			LocalTime:  globalConfig.UseLocalTime,
			Compress:   globalConfig.Compress,
		})
	}
	multi := zerolog.MultiLevelWriter(writers...)
	logger := zerolog.New(multi).With().Timestamp().Logger()
	return &logger, nil
}

func NewLogger(c Config) *zerolog.Logger {
	writers := []io.Writer{&NullWriter{}}
	if !c.EnableConsolColor {
		DefaultConsoleWriter.NoColor = true
	}
	if c.EnableConsolWriter {
		writers = append(writers, DefaultConsoleWriter)
	}
	if c.EnableFileWriter {
		writers = append(writers, &lumberjack.Logger{
			Filename:   c.FileName,
			MaxSize:    c.MaxSizeMB,
			MaxAge:     c.MaxAgeDays,
			MaxBackups: c.MaxBackups,
			LocalTime:  c.UseLocalTime,
			Compress:   c.Compress,
		})
	}
	multi := zerolog.MultiLevelWriter(writers...)

	level, err := zerolog.ParseLevel(c.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	logger := zerolog.New(multi).Level(level).With().Timestamp().Logger()
	return &logger
}

func NewGlobalLogger() *zerolog.Logger {
	return NewLogger(*globalConfig)
}

func ResetCoreDebugLogger() {
	CoreDebugLogger = NullLogger
}
