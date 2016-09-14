package rays

import "strconv"

// PIDParam contains current PID
type PIDParam int

// GetKey is slf.Param implementation
func (r PIDParam) GetKey() string { return "pid" }

// GetRaw is slf.Param implementation
func (r PIDParam) GetRaw() interface{} { return string(r) }

// Value return PID value as integer
func (r PIDParam) Value() int { return int(r) }

// String returns string representation of param value
func (r PIDParam) String() string { return strconv.Itoa(int(r)) }
