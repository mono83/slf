package slf

import "time"

// StatsReporter describes common metrics reporter
type StatsReporter interface {
	// IncCounter increments a statsd-like counter with optional params
	IncCounter(name string, value int64, p ...Param)

	// UpdateGauge a statsd-like gauge ("set" of the value) with optional tags
	UpdateGauge(name string, value int64, p ...Param)

	// RecordTimer records a statsd-like timer with optional tags
	RecordTimer(name string, d time.Duration, p ...Param)

	// Timer builds returns timer, bound to this StatsReporter
	Timer(name string, p ...Param) Timer
}
