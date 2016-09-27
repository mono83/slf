package slf

import (
	"github.com/mono83/slf/params"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWhiteListParamsFilter(t *testing.T) {
	assert := assert.New(t)

	set := []Param{
		params.String{Key: "one", Value: "goo"},
		params.String{Key: "two", Value: "some"},
		params.Int{Key: "int", Value: 34},
	}

	wl1 := NewWhiteListParamsFilter(nil)
	assert.NotNil(wl1(set))
	assert.Len(wl1(set), 0)
	assert.NotNil(wl1(nil))
	assert.Len(wl1(nil), 0)

	wl2 := NewWhiteListParamsFilter([]string{})
	assert.NotNil(wl2(set))
	assert.Len(wl2(set), 0)
	assert.NotNil(wl2(nil))
	assert.Len(wl2(nil), 0)

	wl3 := NewWhiteListParamsFilter([]string{"other"})
	assert.NotNil(wl3(set))
	assert.Len(wl3(set), 0)
	assert.NotNil(wl3(nil))
	assert.Len(wl3(nil), 0)

	wl4 := NewWhiteListParamsFilter([]string{"two", "int"})
	assert.NotNil(wl4(set))
	assert.Len(wl4(set), 2)
	assert.Equal("some", wl4(set)[0].GetRaw())
	assert.Equal(34, wl4(set)[1].GetRaw())
}

func TestNewBlackListParamsFilter(t *testing.T) {
	assert := assert.New(t)

	set := []Param{
		params.String{Key: "one", Value: "goo"},
		params.String{Key: "two", Value: "some"},
		params.Int{Key: "int", Value: 34},
	}

	bl1 := NewBlackListParamsFilter(nil)
	assert.Len(bl1(set), 3)

	bl2 := NewBlackListParamsFilter([]string{})
	assert.Len(bl2(set), 3)

	bl3 := NewBlackListParamsFilter([]string{"other"})
	assert.Len(bl3(set), 3)

	bl4 := NewBlackListParamsFilter([]string{"two", "int"})
	assert.Len(bl4(set), 1)
	assert.Equal("goo", bl4(set)[0].GetRaw())
}
