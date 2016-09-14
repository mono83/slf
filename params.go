package slf

// Param represents parameter, sent to logger or metrics
type Param interface {
	// GetKey returns param key name
	GetKey() string
	// GetRaw returns raw value
	GetRaw() interface{}
	// String returns string representation of param value
	String() string
}
