package tcc

import (
	"github.com/mushroomsir/logger/alog"
)

// LoggerInterface ...
type LoggerInterface interface {
	Err(err error)
	Warning(err error)
}

// Logger ...
type Logger struct {
}

// Err ...
func (a *Logger) Err(err error) {
	alog.Errf("tcc, %s", err.Error())
}

// Warning ...
func (a *Logger) Warning(err error) {
	alog.Warningf("tcc, %s", err.Error())
}
