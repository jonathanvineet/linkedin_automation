package logger

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

// Initialize sets up the structured logger
func Initialize(level string, logFilePath string) error {
	Log = logrus.New()

	// Set log level
	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		parsedLevel = logrus.InfoLevel
	}
	Log.SetLevel(parsedLevel)

	// Set JSON formatter for structured logging
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     false,
	})

	// Create logs directory if it doesn't exist
	if logFilePath != "" {
		logDir := filepath.Dir(logFilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}

		// Create or open log file
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}

		// Write to both file and stdout
		Log.SetOutput(file)
	} else {
		Log.SetOutput(os.Stdout)
	}

	Log.Info("Logger initialized successfully")
	return nil
}

// GetLogger returns the singleton logger instance
func GetLogger() *logrus.Logger {
	if Log == nil {
		Log = logrus.New()
		Log.SetLevel(logrus.InfoLevel)
		Log.SetFormatter(&logrus.JSONFormatter{})
	}
	return Log
}

// LogAction logs an automation action with context
func LogAction(action string, persona string, delayMs int, target string, fields map[string]interface{}) {
	if Log == nil {
		return
	}

	entry := Log.WithFields(logrus.Fields{
		"action":   action,
		"persona":  persona,
		"delay_ms": delayMs,
		"target":   target,
	})

	for k, v := range fields {
		entry = entry.WithField(k, v)
	}

	entry.Info("Automation action executed")
}
