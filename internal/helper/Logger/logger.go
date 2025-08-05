package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LOG FORMATTER
type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now()
	hour := timestamp.Hour() & 12

	if hour == 0 {
		hour = 12
	}
	min := fmt.Sprintf("%02d", timestamp.Minute())
	sec := fmt.Sprintf("%02d", timestamp.Second())

	ampm := "AM"

	if timestamp.Hour() >= 12 {
		ampm = "PM"
	}

	date := timestamp.Format("01-01-2001") // MM-DD-YYYY

	timeString := fmt.Sprintf("%s %d:%s:%s %s", date, hour, min, sec, ampm)
	logLine := fmt.Sprintf("%s [%s]: %s \n", timeString, entry.Level.String(), entry.Message)

	return []byte(logLine), nil
}

// INIT LOGGER

func InitLogger() *logrus.Logger {
	log := logrus.New()
	// SET CUSTOM FORMATTER
	log.SetFormatter(new(CustomFormatter))

	//LOG OUTPUT TO CONSOLE
	log.SetOutput(os.Stdout)

	now := time.Now()

	// logDir := "/var/log/SnehalayaaLogs"
	// os.MkdirAll(logDir, 0755)
	// filename := fmt.Sprintf("%s/Log_%02d_%02d_%d.log", logDir, now.Day(), now.Month(), now.Year())

	filename := fmt.Sprintf("Logs/Log_%02d_%02d_%d.log", now.Day(), now.Month(), now.Year())
	os.MkdirAll("Logs", 0755)

	//FILE ROTATION SETUP
	logFile := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10,
		MaxBackups: 0,
		MaxAge:     7,
		Compress:   true,
	}

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetLevel(logrus.InfoLevel) // MATCH INFO LEVEL

	return log
}
