package memstats

import (
	"github.com/mono83/slf"
	"sync"
)

// New returns metrics receiver, that holds all data
// in memory map.
func New() *MemStats {
	return &MemStats{
		values: map[string]int64{},
	}
}

// MemStats structure is Receiver, that stores incremental and gauge data
// Params are ignored
type MemStats struct {
	m sync.Mutex

	values map[string]int64
}

// Receive handles incoming event
func (m *MemStats) Receive(e slf.Event) {
	if e.Type == slf.TypeInc {
		m.m.Lock()
		prev, _ := m.values[e.Content]
		m.values[e.Content] = prev + e.I64
		m.m.Unlock()
	} else if e.Type == slf.TypeGauge {
		m.m.Lock()
		m.values[e.Content] = e.I64
		m.m.Unlock()
	}
}

// Get return metrics value associated with key
func (m *MemStats) Get(key string) (value int64, found bool) {
	m.m.Lock()
	defer m.m.Unlock()

	value, found = m.values[key]
	return
}

// Values returns copy of metrics values, currently available in MemStats
func (m *MemStats) Values() map[string]int64 {
	m.m.Lock()
	defer m.m.Unlock()

	clone := make(map[string]int64, len(m.values))
	for k, v := range m.values {
		clone[k] = v
	}

	return clone
}
