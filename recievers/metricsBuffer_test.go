package slf

import (
	"github.com/mono83/slf"
	"github.com/mono83/slf/wd"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestNewMetricsBuffer(t *testing.T) {
	assert := assert.New(t)
	wg := sync.WaitGroup{}
	wg.Add(1)

	mp := map[string]slf.Event{}

	before := time.Now()
	mb := NewMetricsBuffer(100*time.Millisecond, func(es []slf.Event) {
		// Copy to map to ease assertions
		for _, v := range es {
			mp[getHashKey(v)] = v
		}
		wg.Done()
	})

	// Sending metrics
	w := wd.Custom("any", "foo.", mb)

	w.UpdateGauge("gauge", 100500)
	w.UpdateGauge("gauge2", 675)
	w.UpdateGauge("gauge2", 2)
	w.UpdateGauge("gauge", 675, wd.IntParam("customerId", 15), wd.StringParam("host", "localhost"))
	w.UpdateGauge("gauge", -43, wd.StringParam("host", "localhost"), wd.IntParam("customerId", 15))

	w.IncCounter("cnt", 1)
	w.IncCounter("cnt", 2)
	w.IncCounter("cnt", 5, wd.IntParam("customerId", 15), wd.StringParam("host", "localhost"))
	w.IncCounter("cnt", -4, wd.StringParam("host", "localhost"), wd.IntParam("customerId", 15))

	w.RecordTimer("delta", time.Second)
	w.RecordTimer("delta2", time.Microsecond*15)

	// Waiting
	wg.Wait()
	assert.True(time.Now().Sub(before).Seconds() >= 0.1)

	assert.Equal(int64(100500), mp["foo.gauge"].I64)
	assert.Equal(int64(2), mp["foo.gauge2"].I64)
	assert.Equal(int64(-43), mp["foo.gauge\tcustomerId\t15\thost\tlocalhost"].I64)
	assert.Equal(int64(3), mp["foo.cnt"].I64)
	assert.Equal(int64(1000000000), mp["foo.delta"].I64)
	assert.Equal(int64(15000), mp["foo.delta2"].I64)

}

func TestGetHashKey(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("foo", getHashKey(slf.Event{Content: "foo"}))
	assert.Equal("foo\tbar\t10", getHashKey(slf.Event{Content: "foo", Params: []slf.Param{wd.IntParam("bar", 10)}}))
	assert.Equal("foo\tbar\t10\tbaz\t0.33", getHashKey(slf.Event{Content: "foo", Params: []slf.Param{
		wd.IntParam("bar", 10),
		wd.FloatParam("baz", 0.33),
	}}))
	assert.Equal("foo\tbar\t10\tbaz\t0.33", getHashKey(slf.Event{Content: "foo", Params: []slf.Param{
		wd.FloatParam("baz", 0.33),
		wd.IntParam("bar", 10),
	}}))
}

func BenchmarkGetHashKeyNull(b *testing.B) {
	e := slf.Event{Content: "foo"}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		getHashKey(e)
	}
}

func BenchmarkGetHashKey3(b *testing.B) {
	e := slf.Event{
		Content: "foo",
		Params: []slf.Param{
			wd.StringParam("foo", "bar"),
			wd.IntParam("asda4", 32423),
			wd.FloatParam("aaa", 0.234),
		},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		getHashKey(e)
	}
}

func BenchmarkGetHashKey10(b *testing.B) {
	e := slf.Event{
		Content: "foo",
		Params: []slf.Param{
			wd.StringParam("foo", "bar"),
			wd.StringParam("foo1", "bar"),
			wd.StringParam("foo2", "bar"),
			wd.StringParam("foo3", "bar"),
			wd.StringParam("foo4", "bar"),
			wd.IntParam("asda", 32423),
			wd.IntParam("asda1", 32423),
			wd.IntParam("asda2", 32423),
			wd.IntParam("asda3", 32423),
			wd.IntParam("asda4", 32423),
		},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		getHashKey(e)
	}
}
