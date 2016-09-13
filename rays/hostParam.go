package rays

// HostParam contains current hostname
type HostParam string

// GetKey is slf.Param implementation
func (r HostParam) GetKey() string { return "host" }

// GetRaw is slf.Param implementation
func (r HostParam) GetRaw() interface{} { return string(r) }
func (r HostParam) String() string      { return string(r) }
