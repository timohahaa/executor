package logger

import (
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func New(logFileName, logPath string) *log.Logger {
	l := log.New()
	l.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339,
		PrettyPrint:     true,
	})
	l.SetReportCaller(true)

	// еще добавим файлы в кастом хуках
	l.SetOutput(os.Stdout)
	logFile, err := os.OpenFile(logPath+logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error creating a log file: %v", err)
		return nil
	}

	l.SetOutput(io.Discard)

	l.AddHook(&customHook{
		levels:  log.AllLevels,
		writers: []io.Writer{os.Stdout, logFile},
	})

	return l
}
