package mutators

import (
	"github.com/mono83/slf"
	"github.com/mono83/slf/params"
	"github.com/mono83/slf/rays"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddInstanceIDToLog(t *testing.T) {
	assert := assert.New(t)

	mod := AddInstanceIDToLog{}
	_, ok := interface{}(mod).(slf.Mutator)
	assert.True(ok, "AddInstanceIDToLog must implement slf.Mutator")

	// Metrics not affected
	for _, t := range []byte{slf.TypeGauge, slf.TypeInc, slf.TypeDuration} {
		pkt := new(slf.Event)
		pkt.Type = t

		assert.Empty(pkt.Params)
		mod.Modify(pkt)
		assert.Empty(pkt.Params)
	}

	// Logs are affected
	for _, t := range []byte{slf.TypeTrace, slf.TypeDebug, slf.TypeInfo, slf.TypeWarning, slf.TypeError, slf.TypeEmergency, slf.TypeAlert} {
		pkt := new(slf.Event)
		pkt.Type = t

		assert.Empty(pkt.Params)
		mod.Modify(pkt)
		assert.NotEmpty(pkt.Params)
		assert.Len(pkt.Params, 1)
		assert.Equal(rays.InstanceID, pkt.Params[0])
	}

	// Append applies
	pkt := new(slf.Event)
	pkt.Type = slf.TypeInfo
	pkt.Params = []slf.Param{params.Nil{Key: "foo"}}
	mod.Modify(pkt)
	assert.Len(pkt.Params, 2)
	assert.Equal(rays.InstanceID, pkt.Params[1])
	assert.Equal("foo", pkt.Params[0].GetKey())
}
