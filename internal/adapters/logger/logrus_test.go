package logger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogrusAdapter(t *testing.T) {
	// Test dizini oluştur
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	// Logger oluştur
	logger, err := NewLogrusAdapter(logPath)
	assert.NoError(t, err)
	assert.NotNil(t, logger)

	// Log dosyasının oluşturulduğunu kontrol et
	_, err = os.Stat(logPath)
	assert.NoError(t, err)
}

func TestLogrusAdapter_LogMethods(t *testing.T) {
	// Test dizini oluştur
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	// Logger oluştur
	logger, err := NewLogrusAdapter(logPath)
	assert.NoError(t, err)

	// Test mesajları
	testCases := []struct {
		name    string
		logFunc func()
		level   string
		message string
	}{
		{
			name: "Info log",
			logFunc: func() {
				logger.Info("info message")
			},
			level:   "info",
			message: "info message",
		},
		{
			name: "Infof log",
			logFunc: func() {
				logger.Infof("infof message %s", "test")
			},
			level:   "info",
			message: "infof message test",
		},
		{
			name: "Error log",
			logFunc: func() {
				logger.Error("error message")
			},
			level:   "error",
			message: "error message",
		},
		{
			name: "Errorf log",
			logFunc: func() {
				logger.Errorf("errorf message %s", "test")
			},
			level:   "error",
			message: "errorf message test",
		},
		{
			name: "Debug log",
			logFunc: func() {
				logger.Debug("debug message")
			},
			level:   "debug",
			message: "debug message",
		},
		{
			name: "Debugf log",
			logFunc: func() {
				logger.Debugf("debugf message %s", "test")
			},
			level:   "debug",
			message: "debugf message test",
		},
		{
			name: "Warn log",
			logFunc: func() {
				logger.Warn("warn message")
			},
			level:   "warning",
			message: "warn message",
		},
		{
			name: "Warnf log",
			logFunc: func() {
				logger.Warnf("warnf message %s", "test")
			},
			level:   "warning",
			message: "warnf message test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Log dosyasını temizle
			err := os.Truncate(logPath, 0)
			assert.NoError(t, err)

			// Log fonksiyonunu çağır
			tc.logFunc()

			// Log dosyasını oku
			content, err := os.ReadFile(logPath)
			assert.NoError(t, err)

			// JSON'ı parse et
			var logEntry map[string]interface{}
			err = json.Unmarshal(content, &logEntry)
			assert.NoError(t, err)

			// Log seviyesi ve mesajı kontrol et
			assert.Equal(t, tc.level, logEntry["level"])
			assert.Equal(t, tc.message, logEntry["msg"])
			assert.NotEmpty(t, logEntry["time"])
		})
	}
}

func TestLogrusAdapter_InvalidLogPath(t *testing.T) {
	// Geçersiz log dizini
	logPath := "/invalid/path/test.log"

	// Logger oluşturmayı dene
	logger, err := NewLogrusAdapter(logPath)
	assert.Error(t, err)
	assert.Nil(t, logger)
}
