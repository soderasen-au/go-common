package loggers

import (
	"io"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

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

type LogFileCloser func() error

func GetLoggerWithCloser(fn string) (*zerolog.Logger, LogFileCloser, error) {
	writers := []io.Writer{&NullWriter{}}
	if !EnableConsolColor {
		DefaultConsoleWriter.NoColor = true
	}
	if EnableConsolWriter {
		writers = append(writers, DefaultConsoleWriter)
	}
	var fnWriteCloser io.WriteCloser
	if EnableFileWriter && len(fn) > 0 {
		fnWriteCloser = &lumberjack.Logger{
			Filename:   fn,
			MaxSize:    globalConfig.MaxSizeMB,
			MaxAge:     globalConfig.MaxAgeDays,
			MaxBackups: globalConfig.MaxBackups,
			LocalTime:  globalConfig.UseLocalTime,
			Compress:   globalConfig.Compress,
		}
		writers = append(writers, fnWriteCloser)
	}
	multi := zerolog.MultiLevelWriter(writers...)
	logger := zerolog.New(multi).With().Timestamp().Logger()
	return &logger, func() error {
		if fnWriteCloser != nil {
			return fnWriteCloser.Close()
		}
		return nil
	}, nil
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
