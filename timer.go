package slf

import (
	"github.com/mono83/slf/params"
	"time"
)

// Timer is interface to metrics timer
type Timer interface {
	// Done stops timer and returns elapsed time
	// If timer was already stopped, just returns elapsed time
	Done() time.Duration

	// Params returns current timer value as param
	Param(string) Param
}

// NewTimer builds new Timer
func NewTimer(name string, p []Param, target StatsReporter) Timer {
	return &timer{
		started: time.Now(),
		target:  target,
		name:    name,
		p:       p,
	}
}

type timer struct {
	started, finished time.Time
	target            StatsReporter

	name string
	p    []Param
}

func (t *timer) Done() time.Duration {
	if t.finished.IsZero() {
		// Not finished yet
		t.finished = time.Now()
		t.target.RecordTimer(t.name, t.finished.Sub(t.started), t.p...)
	}

	return t.finished.Sub(t.started)
}

func (t timer) Param(n string) Param {
	if t.finished.IsZero() {
		return params.Nil{Key: n}
	}

	return params.Duration{Key: n, Value: t.Done()}
}
