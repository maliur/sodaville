package logger

import "log"

type Slogger struct {
	logger *log.Logger
}

func NewSlogger(logger *log.Logger) *Slogger {
	return &Slogger{logger}
}

func (sl *Slogger) Error(msg string) {
	sl.logger.Printf("[ERROR] %s", msg)
}

func (sl *Slogger) Warn(msg string) {
	sl.logger.Printf("[WARN] %s", msg)
}

func (sl *Slogger) Info(msg string) {
	sl.logger.Printf("[INFO] %s", msg)
}
