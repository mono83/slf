package statsd

import (
	"github.com/mono83/slf"
	"github.com/mono83/slf/wd"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPacketBuilder(t *testing.T) {
	assert := assert.New(t)

	var pb PacketBuilder
	pb = NewPacketBuilder("foo", 1000)
	assert.Equal(0, pb.Size())
	pb.WriteEvent(slf.Event{Type: slf.TypeGauge, Content: "bar", I64: 150})
	assert.Equal(1, pb.Size())
	assert.Equal("foo.bar:150|g", string(pb.Bytes()))

	pb = NewPacketBuilder("", 1000)
	pb.WriteEvent(slf.Event{Type: slf.TypeDuration, Content: "bar", I64: int64(time.Millisecond * 15)})
	assert.Equal(1, pb.Size())
	assert.Equal("bar:15000|ms", string(pb.Bytes()))

	pb = NewPacketBuilder("xxx.", 1000)
	pb.WriteEvent(slf.Event{Type: slf.TypeInc, Content: "y", I64: 33})
	assert.Equal(1, pb.Size())
	assert.Equal("xxx.y:33|c", string(pb.Bytes()))

	pb = NewPacketBuilder("xxx.", 1000)
	pb.WriteEvent(slf.Event{Type: slf.TypeInc, Content: "y", I64: 33})
	pb.WriteEvent(slf.Event{Type: slf.TypeInc, Content: "z", I64: 3})
	assert.Equal(2, pb.Size())
	assert.Equal("xxx.y:33|c\nxxx.z:3|c", string(pb.Bytes()))
}

func TestPacketBuilderWithParams(t *testing.T) {
	assert := assert.New(t)

	allowed := map[string]bool{}
	allowed["host"] = true
	allowed["port"] = true

	var pb PacketBuilder
	pb = NewPacketWithParamsBuilder("foo", 1000, allowed)
	assert.Equal(0, pb.Size())
	pb.WriteEvent(slf.Event{Type: slf.TypeGauge, Content: "bar", I64: 150, Params: []slf.Param{wd.NameParam("test")}})
	assert.Equal(1, pb.Size())
	assert.Equal("foo.bar:150|g", string(pb.Bytes()))

	pb = NewPacketWithParamsBuilder("foo", 1000, allowed)
	pb.WriteEvent(slf.Event{Type: slf.TypeGauge, Content: "bar", I64: 150, Params: []slf.Param{wd.StringParam("host", "localhost")}})
	assert.Equal(1, pb.Size())
	assert.Equal("foo.bar:150|g|@1.0|#host:localhost", string(pb.Bytes()))

	pb = NewPacketWithParamsBuilder("foo", 1000, allowed)
	pb.WriteEvent(slf.Event{Type: slf.TypeGauge, Content: "bar", I64: 150, Params: []slf.Param{
		wd.NameParam("test"),
		wd.StringParam("host", "host:local"),
		wd.IntParam("port", 3306),
	}})
	assert.Equal(1, pb.Size())
	assert.Equal("foo.bar:150|g|@1.0|#host:host_local,port:3306", string(pb.Bytes()))
}

func BenchmarkCommonPacket5(b *testing.B) {
	allowed := map[string]bool{}
	allowed["host"] = true
	allowed["port"] = true

	events := []slf.Event{
		{Type: slf.TypeGauge, Content: "bar", I64: 150, Params: []slf.Param{
			wd.NameParam("test"),
			wd.StringParam("host", "host:local"),
			wd.IntParam("port", 3306),
		}},
		{Type: slf.TypeInc, Content: "y", I64: 33},
		{Type: slf.TypeInc, Content: "z", I64: 3},
		{Type: slf.TypeDuration, Content: "bar", I64: int64(time.Millisecond * 15)},
		{Type: slf.TypeGauge, Content: "bar", I64: 150},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pb := NewPacketWithParamsBuilder("any", 1000, allowed)
		for _, e := range events {
			pb.WriteEvent(e)
		}
		pb.Bytes()
	}
}

func TestSanitizeParamKey(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("foo", string(SanitizeParamValue("foo")))
	assert.Equal("Ba.R", string(SanitizeParamValue("  Ba.R\n\t")))
	assert.Equal("1234567890", string(SanitizeParamValue("1234567890")))
	assert.Equal("abcdefghijklmnopqrstuvwxyz", string(SanitizeParamValue("abcdefghijklmnopqrstuvwxyz")))
	assert.Equal("_x_", string(SanitizeParamValue("+x-")))
}

func BenchmarkSanitizeParamKeyShort(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		SanitizeParamValue("\tabc 123.")
	}
}

func BenchmarkSanitizeParamKeyLong(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		SanitizeParamValue("\tabc 123. QWE eVen more characters to analyze and replace !")
	}
}
