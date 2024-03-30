package logger

import (
	"io"

	log "github.com/sirupsen/logrus"
)

type customHook struct {
	levels  []log.Level
	writers []io.Writer
}

func NewCunstomHook(levels []log.Level, writers []io.Writer) *customHook {
	return &customHook{
		levels:  levels,
		writers: writers,
	}
}

func (h *customHook) Levels() []log.Level {
	return h.levels
}

func (h *customHook) Fire(entry *log.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range h.writers {
		_, err = w.Write([]byte(line))
	}

	return err
}
