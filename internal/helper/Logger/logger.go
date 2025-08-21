package logger

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LOG FORMATTER
type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s [%s]: %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		entry.Level.String(),
		entry.Message)), nil
}

// INIT LOGGER

func InitLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(new(CustomFormatter))

	now := time.Now()
	filename := fmt.Sprintf("Logs/Log_%02d_%02d_%d.log", now.Day(), now.Month(), now.Year())
	os.MkdirAll("Logs", 0755)

	// create file once with required permissions
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		file.Close()
		// explicitly set permissions
		syscall.Chmod(filename, 0644)
	}

	// lumberjack rotation
	logFile := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10, // MB
		MaxBackups: 0,
		MaxAge:     7, // days
		Compress:   true,
	}

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetLevel(logrus.InfoLevel)

	return log
}
