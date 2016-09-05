package health

import (
	"github.com/mono83/slf/wd"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStartHealthMonitor(t *testing.T) {
	assert := assert.New(t)

	target := mapSR(map[string]int64{})
	sendStats(wd.Custom("", "health.", target))

	assert.Len(target, 10)

	has := func(key string) {
		_, ok := target[key]
		assert.True(ok, "Key not found "+key)
	}

	has("health.gcs")
	has("health.goroutines")
	has("health.sys.malloc")
	has("health.sys.free")
	has("health.heap.alloc")
	has("health.heap.inuse")
	has("health.heap.sys")
	has("health.heap.objects")
	has("health.heap.nextgc")
	has("health.uptime")
}
