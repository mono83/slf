package wd

import (
	"github.com/mono83/slf"
	"time"
)

// Watchdog interface combines logger and metrics functionality
type Watchdog interface {
	slf.Logger
	slf.StatsReporter

	// WithParams returns copy of watchdog with pre-setted params
	// Useful for contex-based logging
	WithParams(...slf.Param) Watchdog
}

// New builds and returns new watchdog
func New(name, metricsPrefix string) Watchdog {
	return &watchdog{name: name, prefix: metricsPrefix, pipe: stdDispatcher.Receive}
}

type watchdog struct {
	name, prefix string
	params       []slf.Param
	pipe         func(slf.Event)
}

func (w watchdog) join(p []slf.Param) []slf.Param {
	if len(w.params) == 0 {
		return p
	}

	return append(w.params, p...)
}

func (w watchdog) Trace(message string, p ...slf.Param) {
	w.pipe(newLog(slf.TypeTrace, w.name, message, w.join(p)))
}

func (w watchdog) Debug(message string, p ...slf.Param) {
	w.pipe(newLog(slf.TypeDebug, w.name, message, w.join(p)))
}

func (w watchdog) Info(message string, p ...slf.Param) {
	w.pipe(newLog(slf.TypeInfo, w.name, message, w.join(p)))
}

func (w watchdog) Warning(message string, p ...slf.Param) {
	w.pipe(newLog(slf.TypeWarning, w.name, message, w.join(p)))
}

func (w watchdog) Error(message string, p ...slf.Param) {
	w.pipe(newLog(slf.TypeError, w.name, message, w.join(p)))
}

func (w watchdog) Alert(message string, p ...slf.Param) {
	w.pipe(newLog(slf.TypeAlert, w.name, message, w.join(p)))
}

func (w watchdog) Emergency(message string, p ...slf.Param) {
	w.pipe(newLog(slf.TypeEmergency, w.name, message, w.join(p)))
}

func (w watchdog) IncCounter(name string, value int64, p ...slf.Param) {
	w.pipe(newMetrics(slf.TypeInc, w.name, w.prefix+name, value, w.join(p)))
}

func (w watchdog) UpdateGauge(name string, value int64, p ...slf.Param) {
	w.pipe(newMetrics(slf.TypeGauge, w.name, w.prefix+name, value, w.join(p)))
}

func (w watchdog) RecordTimer(name string, d time.Duration, p ...slf.Param) {
	w.pipe(newMetrics(slf.TypeGauge, w.name, w.prefix+name, int64(d), w.join(p)))
}

func (w watchdog) Timer(name string, p ...slf.Param) slf.Timer {
	return slf.NewTimer(name, w.join(p), w)
}

func (w watchdog) WithParams(p ...slf.Param) Watchdog {
	if len(p) == 0 {
		return w
	}

	return &watchdog{
		name:   w.name,
		prefix: w.prefix,
		pipe:   w.pipe,
		params: append(w.params, p...),
	}
}
