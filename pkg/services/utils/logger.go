package utils

import (
	"fmt"
	"log"
)

type Logger struct {
	context string
	logger  *log.Logger
}

func NewLogger(context string) Logger {
	return Logger{
		logger:  log.Default(),
		context: context,
	}
}

func (l Logger) WithContext(context string) Logger {
	return Logger{
		logger:  l.logger,
		context: fmt.Sprintf("%s:%s", l.context, context),
	}
}

func (l Logger) Print(msg any) {
	l.logger.Printf("[%s] %v", l.context, msg)
}

func (l Logger) Println(msg any) {
	l.Print(fmt.Sprintln(msg))
}

func (l Logger) Printf(format string, v ...any) {
	l.Print(fmt.Sprintf(format, v...))
}

func (l Logger) Info(msg any) {
	l.Printf("[INFO] %v", msg)
}

func (l Logger) Infoln(msg any) {
	l.Info(fmt.Sprintln(msg))
}

func (l Logger) Infof(format string, v ...any) {
	l.Info(fmt.Sprintf(format, v...))
}

func (l Logger) Warn(msg any) {
	l.Printf("[WARN] %v", msg)
}

func (l Logger) Warnln(msg any) {
	l.Warn(fmt.Sprintln(msg))
}

func (l Logger) Warnf(format string, v ...any) {
	l.Warn(fmt.Sprintf(format, v...))
}

func (l Logger) Error(msg any) {
	l.logger.Printf("[ERROR] %v", msg)
}

func (l Logger) Errorln(msg any) {
	l.Error(fmt.Sprintln(msg))
}

func (l Logger) Errorf(format string, v ...any) {
	l.Errorln(fmt.Sprintf(format, v...))
}
