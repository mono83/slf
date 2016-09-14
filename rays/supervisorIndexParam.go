package rays

import "strconv"

// SupervisorIndexParam is incremental index of worker under supervisord, if it started under it
type SupervisorIndexParam int

// GetKey is slf.Param implementation
func (r SupervisorIndexParam) GetKey() string { return "supervisorIndex" }

// GetRaw is slf.Param implementation
func (r SupervisorIndexParam) GetRaw() interface{} { return int(r) }
func (r SupervisorIndexParam) String() string      { return strconv.Itoa(int(r)) }
