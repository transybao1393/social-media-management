package logger

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	logger *logrus.Logger
	fields Fields
}

func NewLogrusLogger() Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	defaultFields := Fields{
		"service":  "social-media-management",
		"hostname": hostname,
	}

	return &LogrusLogger{
		logger: l,
		fields: defaultFields,
	}
}

func rootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

func (l *LogrusLogger) Fields(data Fields) Logger {
	formatter := time.Now().Format("2006-01-02")

	//- With this solution, every time we call the logger.Fields() method, it will find the path of this file when runtime call and point to correct folder
	logPathLocation2 := filepath.Join(rootDir(), "/logs/")

	fmt.Printf("logPathLocation2 %s\n", logPathLocation2)
	if _, err := os.Stat(logPathLocation2); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(logPathLocation2, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	logFilePath := logPathLocation2 + "/" + formatter + ".log"
	fmt.Printf("logFilePath %s\n", logFilePath)
	logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("Failed to open log file %s for output: %s", logFilePath, err)
	}

	l.logger.SetOutput(io.MultiWriter(os.Stderr, logFile))
	return &LogrusLogger{
		logger: l.logger,
		fields: data,
	}
}

func (l *LogrusLogger) Debug(msg string) {
	l.logger.WithFields(logrus.Fields(l.fields)).Debug(msg)
}

func (l *LogrusLogger) Debugf(msg string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields(l.fields)).Debugf(msg, args...)
}

func (l *LogrusLogger) Info(msg string) {
	l.logger.WithFields(logrus.Fields(l.fields)).Info(msg)
}

func (l *LogrusLogger) Infof(msg string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields(l.fields)).Infof(msg, args...)
}

func (l *LogrusLogger) Warn(msg string) {
	l.logger.WithFields(logrus.Fields(l.fields)).Warn(msg)
}

func (l *LogrusLogger) Warnf(msg string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields(l.fields)).Warnf(msg, args...)
}

func (l *LogrusLogger) Error(err error, msg string) {
	l.logger.WithFields(logrus.Fields(l.fields)).WithError(err).Error(msg)
}

func (l *LogrusLogger) Errorf(err error, msg string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields(l.fields)).WithError(err).Errorf(msg, args...)
}

func (l *LogrusLogger) Fatal(err error, msg string) {
	l.logger.WithFields(logrus.Fields(l.fields)).WithError(err).Fatal(msg)
}

func (l *LogrusLogger) Fatalf(err error, msg string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields(l.fields)).WithError(err).Fatalf(msg, args...)
}

func (l *LogrusLogger) Printf(format string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields(l.fields)).Printf(format, args...)
}

func (l *LogrusLogger) Println(args ...interface{}) {
	l.logger.WithFields(logrus.Fields(l.fields)).Println(args...)
}
