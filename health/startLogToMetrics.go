package health

import (
	"github.com/mono83/slf"
	"github.com/mono83/slf/wd"
)

type metricsAdapter struct {
	rec slf.StatsReporter
}

func (m metricsAdapter) Receive(e slf.Event) {
	switch e.Type {
	case slf.TypeTrace:
		m.rec.IncCounter("trace", 1)
	case slf.TypeDebug:
		m.rec.IncCounter("debug", 1)
	case slf.TypeInfo:
		m.rec.IncCounter("info", 1)
	case slf.TypeWarning:
		m.rec.IncCounter("warning", 1)
	case slf.TypeError:
		m.rec.IncCounter("error", 1)
	case slf.TypeAlert:
		m.rec.IncCounter("alert", 1)
	case slf.TypeEmergency:
		m.rec.IncCounter("emergency", 1)
	}
}

// StartLogToMetrics registers and binds special events receiver,
// that counts log messages per-level and sends counters as
// metrics events
//
// Common usage: StartLogToMetrics(wd.New("", "log.").WithParams(util.HostParam()))
func StartLogToMetrics(rec slf.StatsReporter) {
	ma := &metricsAdapter{rec: rec}
	wd.AddReceiver(ma)
}
