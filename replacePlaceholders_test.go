package slf

import (
	"errors"
	"github.com/mono83/slf/params"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReplacePlaceholders(t *testing.T) {
	assert.Equal(
		t,
		"[foo] name is [15] with [Test] or <!:p1> or [nil]",
		ReplacePlaceholders(
			":start name is :id with :err or :p1 or :p2",
			[]Param{
				params.Int{Key: "id", Value: 15},
				params.String{Key: "start", Value: "foo"},
				params.Error{Key: "err", Value: errors.New("Test")},
				params.Error{Key: "p2"},
			},
			true,
		),
	)

	assert.Equal(
		t,
		"21.12 24234234 nil true",
		ReplacePlaceholders(
			":This :is :something :other",
			[]Param{
				params.Int64{Key: "is", Value: int64(24234234)},
				params.Nil{Key: "something"},
				params.Bool{Key: "other", Value: true},
				params.Float64{Key: "This", Value: 21.12},
			},
			false,
		),
	)
}
