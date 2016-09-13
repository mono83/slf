package rays

// RayID contains unique id for process identification
type RayID string

// GetKey is slf.Param implementation
func (r RayID) GetKey() string { return "rayId" }

// GetRaw is slf.Param implementation
func (r RayID) GetRaw() interface{} { return string(r) }
func (r RayID) String() string      { return string(r) }
