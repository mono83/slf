package wd

import (
	"github.com/mono83/slf"
	"github.com/mono83/slf/params"
)

var stdDispatcher = &slf.Dispatcher{}

// AddReceiver registers receiver in default receivers pool
func AddReceiver(r slf.Receiver) {
	stdDispatcher.AddReceiver(r)
}

// AddMutator registers mutator
func AddMutator(m slf.Mutator) {
	stdDispatcher.AddMutator(m)
}

// IntParam return integer Param
func IntParam(name string, value int) slf.Param {
	return params.Int{Key: name, Value: value}
}

// Int64Param return int64 Param
func Int64Param(name string, value int64) slf.Param {
	return params.Int64{Key: name, Value: value}
}

// ID64Param return int64 Param with key "id"
func ID64Param(value int64) slf.Param {
	return params.Int64{Key: "id", Value: value}
}

// ErrParam returns Param for errors with key "err"
func ErrParam(err error) slf.Param {
	if err == nil {
		return params.Nil{Key: "err"}
	}

	return params.Error{Key: "err", Value: err}
}

// StringParam returns Param for strings
func StringParam(name, value string) slf.Param {
	return params.String{Key: name, Value: value}
}

// NameParam is alias for StringParam with key "name"
func NameParam(value string) slf.Param {
	return params.String{Key: "name", Value: value}
}

// CountParam is alias for IntParam with key "count"
func CountParam(value int) slf.Param {
	return params.Int{Key: "count", Value: value}
}

// FloatParam return Param for float values
func FloatParam(name string, value float64) slf.Param {
	return params.Float64{Key: name, Value: value}
}
