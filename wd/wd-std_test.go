package wd

import (
	"errors"
	"github.com/mono83/slf/params"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntParam(t *testing.T) {
	p := IntParam("foo", 500)

	assert.IsType(t, params.Int{}, p)
	assert.Equal(t, "foo", p.GetKey())
	assert.Equal(t, 500, p.GetRaw())
}

func TestCountParam(t *testing.T) {
	p := CountParam(3231)

	assert.IsType(t, params.Int{}, p)
	assert.Equal(t, "count", p.GetKey())
	assert.Equal(t, 3231, p.GetRaw())
}

func TestErrParam(t *testing.T) {
	err := errors.New("Some text")
	p := ErrParam(err)

	assert.IsType(t, params.Error{}, p)
	assert.Equal(t, "err", p.GetKey())
	assert.Equal(t, err, p.GetRaw())

	p = ErrParam(nil)
	assert.IsType(t, params.Nil{}, p)
}

func TestStringParam(t *testing.T) {
	p := StringParam("bar", "baz")

	assert.IsType(t, params.String{}, p)
	assert.Equal(t, "bar", p.GetKey())
	assert.Equal(t, "baz", p.GetRaw())
}

func TestNameParam(t *testing.T) {
	p := NameParam("Arya")

	assert.IsType(t, params.String{}, p)
	assert.Equal(t, "name", p.GetKey())
	assert.Equal(t, "Arya", p.GetRaw())
}

func TestFloatParam(t *testing.T) {
	p := FloatParam("foo", 0.2123)

	assert.IsType(t, params.Float64{}, p)
	assert.Equal(t, "foo", p.GetKey())
	assert.Equal(t, 0.2123, p.GetRaw())
}
