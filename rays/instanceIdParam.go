package rays

// InstanceIDParam contains current application unique id
type InstanceIDParam string

// GetKey is slf.Param implementation
func (r InstanceIDParam) GetKey() string { return "instanceId" }

// GetRaw is slf.Param implementation
func (r InstanceIDParam) GetRaw() interface{} { return string(r) }
func (r InstanceIDParam) String() string      { return string(r) }
