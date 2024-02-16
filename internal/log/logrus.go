package log

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type ConfigLogger struct {
	EnableConsole bool
	Level         string
	EnableFile    bool
	FileLocation  string
}

type Logger struct {
	logger *logrus.Logger
	file   *os.File
}

func NewLogger(config ConfigLogger) (*Logger, error) {
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		return nil, err
	}

	lLogger := logrus.New()
	lLogger.SetLevel(level)
	lLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})

	logger := Logger{
		logger: lLogger,
	}

	writers := make([]io.Writer, 0)

	if config.EnableConsole {
		writers = append(writers, os.Stderr)
	}

	if config.EnableFile {
		logger.file, err = os.OpenFile(config.FileLocation, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return nil, err
		}
		writers = append(writers, logger.file)
	}

	lLogger.SetOutput(io.MultiWriter(writers...))

	return &logger, nil
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *Logger) Close() error {
	if l != nil && l.file != nil {
		if err := l.file.Close(); err != nil {
			return err
		}
	}

	return nil
}
