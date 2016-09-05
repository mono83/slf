package params

import "time"

// Bool type Param
type Bool struct {
	Key   string
	Value bool
}

// GetKey returns key of Param
func (b Bool) GetKey() string { return b.Key }

// GetRaw returns value of Param as interface{}
func (b Bool) GetRaw() interface{} { return b.Value }

// Int type Param
type Int struct {
	Key   string
	Value int
}

// GetKey returns key of Param
func (i Int) GetKey() string { return i.Key }

// GetRaw returns value of Param as interface{}
func (i Int) GetRaw() interface{} { return i.Value }

// Int64 type Param
type Int64 struct {
	Key   string
	Value int
}

// GetKey returns key of Param
func (i Int64) GetKey() string { return i.Key }

// GetRaw returns value of Param as interface{}
func (i Int64) GetRaw() interface{} { return i.Value }

// Float64 type Param
type Float64 struct {
	Key   string
	Value float64
}

// GetKey returns key of Param
func (f Float64) GetKey() string { return f.Key }

// GetRaw returns value of Param as interface{}
func (f Float64) GetRaw() interface{} { return f.Value }

// String type Param
type String struct {
	Key, Value string
}

// GetKey returns key of Param
func (s String) GetKey() string { return s.Key }

// GetRaw returns value of Param as interface{}
func (s String) GetRaw() interface{} { return s.Value }

// Error type Param
type Error struct {
	Key   string
	Value error
}

// GetKey returns key of Param
func (e Error) GetKey() string { return e.Key }

// GetRaw returns value of Param as interface{}
func (e Error) GetRaw() interface{} { return e.Value }

// Nil is Param with nil value
type Nil struct {
	Key string
}

// GetKey returns key of Param
func (n Nil) GetKey() string { return n.Key }

// GetRaw returns value of Param as interface{}
func (n Nil) GetRaw() interface{} { return nil }

// Duration type Param
type Duration struct {
	Key   string
	Value time.Duration
}

// GetKey returns key of Param
func (d Duration) GetKey() string { return d.Key }

// GetRaw returns value of Param as interface{}
func (d Duration) GetRaw() interface{} { return d.Value }
