/*
 * hopper - A gRPC API for collecting IoT device event messages
 * Copyright (C) 2022 Brian Reece

 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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
	l.Printf("[ERROR] %v", msg)
}

func (l Logger) Errorln(msg any) {
	l.Error(fmt.Sprintln(msg))
}

func (l Logger) Errorf(format string, v ...any) {
	l.Errorln(fmt.Sprintf(format, v...))
}
