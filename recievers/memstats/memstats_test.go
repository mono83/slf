package memstats

import (
	"github.com/mono83/slf/wd"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemStats(t *testing.T) {
	assert := assert.New(t)

	mem := New()
	assert.NotNil(mem.values)
	assert.Len(mem.values, 0)

	log := wd.Custom("", "foo.", mem)

	log.IncCounter("example", 5)
	assert.Len(mem.values, 1)

	log.IncCounter("twice", 11)
	log.IncCounter("twice", -4)
	assert.Len(mem.values, 2)

	log.UpdateGauge("scale", 331)
	assert.Len(mem.values, 3)

	v, _ := mem.Get("foo.example")
	assert.Equal(int64(5), v)
	v, _ = mem.Get("foo.twice")
	assert.Equal(int64(7), v)
	v, _ = mem.Get("foo.scale")
	assert.Equal(int64(331), v)

	log.UpdateGauge("scale", -44)
	assert.Len(mem.values, 3)
	v, _ = mem.Get("foo.scale")
	assert.Equal(int64(-44), v)

	vals := mem.Values()
	assert.Len(vals, 3)

	v, _ = vals["foo.example"]
	assert.Equal(int64(5), v)
	v, _ = vals["foo.twice"]
	assert.Equal(int64(7), v)
	v, _ = vals["foo.scale"]
	assert.Equal(int64(-44), v)

	// Modifying MemStats - exported values must be not changed
	log.UpdateGauge("scale", 800)
	v, _ = mem.Get("foo.scale")
	assert.Equal(int64(800), v)
	v, _ = vals["foo.scale"]
	assert.Equal(int64(-44), v)
}
