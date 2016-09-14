package slf

import (
	"github.com/mono83/slf"
	"sort"
	"strings"
	"sync"
	"time"
)

type metricsBuffer struct {
	m sync.Mutex

	gauges    map[string]slf.Event
	counters  map[string]slf.Event
	durations []slf.Event

	interval time.Duration
	target   func([]slf.Event)
}

// NewMetricsBuffer builds new metrics buffer for desired flush interval
func NewMetricsBuffer(flushInterval time.Duration, flushTo func([]slf.Event)) slf.Receiver {
	if flushInterval.Nanoseconds() == 0 {
		flushInterval = time.Second * 10
	}

	mb := &metricsBuffer{
		gauges:   map[string]slf.Event{},
		counters: map[string]slf.Event{},
		interval: flushInterval,
		target:   flushTo,
	}

	go mb.flush()
	return mb
}

func (m *metricsBuffer) flush() {
	for {
		time.Sleep(m.interval)

		// Take lock
		m.m.Lock()
		// Copy durations to local variable
		lineBuffer := m.durations
		m.durations = nil
		// Copy counters to local variable
		localCounters := m.counters
		m.counters = map[string]slf.Event{}
		// Copy gauges values under lock - they are persistent
		for _, v := range m.gauges {
			lineBuffer = append(lineBuffer, v)
		}
		// Release lock
		m.m.Unlock()
		// Copy counters to local values - without lock, from local variable
		for _, v := range localCounters {
			lineBuffer = append(lineBuffer, v)
		}

		if m.target != nil {
			m.target(lineBuffer)
		}
	}
}

func (m *metricsBuffer) Receive(e slf.Event) {
	switch e.Type {
	case slf.TypeGauge:
		key := getHashKey(e)
		m.m.Lock()
		defer m.m.Unlock()

		m.gauges[key] = e
	case slf.TypeInc:
		key := getHashKey(e)
		m.m.Lock()
		defer m.m.Unlock()

		save, ok := m.counters[key]
		if ok {
			save.I64 += e.I64
		} else {
			save = e
		}

		m.counters[key] = save
	case slf.TypeDuration:
		m.m.Lock()
		defer m.m.Unlock()

		m.durations = append(m.durations, e)
	}
}

// getHashKey calculates hash, used as key in metrics buffer maps
func getHashKey(e slf.Event) string {
	if len(e.Params) == 0 {
		return e.Content
	}

	names := []string{}
	for _, key := range e.Params {
		names = append(names, key.GetKey()+"\t"+key.String())
	}
	sort.Strings(names)

	return e.Content + "\t" + strings.Join(names, "\t")
}
