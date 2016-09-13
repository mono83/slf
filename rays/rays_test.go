package rays

import (
	"github.com/mono83/slf"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	r := New()
	r2 := New()
	_, ok := interface{}(r).(slf.Param)
	assert.True(ok, "RayID must implement slf.Param")
	assert.Equal("rayId", r.GetKey())
	assert.True(len(r.String()) > 30)
	assert.Equal("-", r.String()[6:7])
	assert.Equal(4, strings.Count(r.String(), "-"))

	assert.NotEqual(r.String(), r2.String())
}

func TestInstanceId(t *testing.T) {
	assert := assert.New(t)

	_, ok := interface{}(InstanceID).(slf.Param)
	assert.True(ok, "InstanceID must implement slf.Param")
	assert.Equal("instanceId", InstanceID.GetKey())
	assert.NotEmpty(InstanceID.String())
}

func TestHost(t *testing.T) {
	assert := assert.New(t)

	_, ok := interface{}(Host).(slf.Param)
	assert.True(ok, "Host must implement slf.Param")
	assert.Equal("host", Host.GetKey())
	assert.NotEmpty(Host.String())
}

func TestPid(t *testing.T) {
	assert := assert.New(t)

	_, ok := interface{}(PID).(slf.Param)
	assert.True(ok, "PID must implement slf.Param")
	assert.Equal("pid", PID.GetKey())
	assert.NotEmpty(PID.Value())
}
