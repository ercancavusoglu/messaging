package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ercancavusoglu/messaging/internal/ports"
	"github.com/sirupsen/logrus"
)

type logrusAdapter struct {
	logger *logrus.Logger
}

// NewLogrusAdapter creates a new logger instance with file and console output
func NewLogrusAdapter(logPath string) (ports.Logger, error) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	// Ensure log directory exists
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return nil, err
	}

	// Open log file
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	// Set output to both file and console
	logger.SetOutput(io.MultiWriter(os.Stdout, file))

	return &logrusAdapter{
		logger: logger,
	}, nil
}

func (l *logrusAdapter) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *logrusAdapter) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *logrusAdapter) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *logrusAdapter) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *logrusAdapter) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *logrusAdapter) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *logrusAdapter) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *logrusAdapter) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *logrusAdapter) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *logrusAdapter) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *logrusAdapter) WithFields(fields map[string]interface{}) ports.Logger {
	return &logrusAdapter{
		logger: l.logger.WithFields(fields).Logger,
	}
}
