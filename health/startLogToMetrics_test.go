package health

import (
	"github.com/mono83/slf"
	"github.com/mono83/slf/wd"
	"github.com/stretchr/testify/assert"
	"testing"
)

// map stats receiver
type mapSR map[string]int64

func (m mapSR) Receive(e slf.Event) {
	if e.Type == slf.TypeInc {
		m[e.Content] = m[e.Content] + e.I64
	} else if e.Type == slf.TypeGauge {
		m[e.Content] = e.I64
	}
}

func TestStartLogToMetrics(t *testing.T) {
	assert := assert.New(t)

	target := mapSR(map[string]int64{})
	receiver := metricsAdapter{rec: wd.Custom("", "prefix.", target)}

	w := wd.Custom("", "", receiver)

	assert.Len(target, 0)
	w.Trace("This is trace")
	assert.Len(target, 1)
	w.Debug("This is debug")
	assert.Len(target, 2)
	w.Info("This is info")
	assert.Len(target, 3)
	w.Warning("This is warning")
	assert.Len(target, 4)
	w.Error("This is error")
	assert.Len(target, 5)
	w.Alert("This is alert")
	assert.Len(target, 6)
	w.Emergency("This is emergency")
	assert.Len(target, 7)

	assert.Equal(int64(1), target["prefix.trace"])
	assert.Equal(int64(1), target["prefix.debug"])
	assert.Equal(int64(1), target["prefix.info"])
	assert.Equal(int64(1), target["prefix.warning"])
	assert.Equal(int64(1), target["prefix.error"])
	assert.Equal(int64(1), target["prefix.alert"])
	assert.Equal(int64(1), target["prefix.emergency"])
}
