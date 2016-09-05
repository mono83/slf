package slf

import "time"

// List of
const (
	TypeTrace     byte = 1
	TypeDebug     byte = 2
	TypeInfo      byte = 3
	TypeWarning   byte = 4
	TypeError     byte = 5
	TypeAlert     byte = 6
	TypeEmergency byte = 7

	TypeInc      byte = 100
	TypeGauge    byte = 101
	TypeDuration byte = 102
)

// Event represents event information
type Event struct {
	Time   time.Time // Event time
	Marker string    // Event owner

	Type    byte    // Event type
	Content string  // Content - either message or metrics name
	Params  []Param // Params slice
	I64     int64   // Increment or duration or gauge value
}

// IsMetrics return true if event is metrics event
func (e Event) IsMetrics() bool {
	return e.Type == TypeInc || e.Type == TypeGauge || e.Type == TypeDuration
}

// IsLog returns true if event is logging event
func (e Event) IsLog() bool {
	return e.Type == TypeInfo || e.Type == TypeDebug || e.Type == TypeWarning || e.Type == TypeError || e.Type == TypeTrace || e.Type == TypeAlert || e.Type == TypeEmergency
}

// StringType returns string representation of log event type
// Will return empty string on metrics or invalid type
func (e Event) StringType() string {
	switch e.Type {
	case TypeTrace:
		return "trace"
	case TypeDebug:
		return "debug"
	case TypeInfo:
		return "info"
	case TypeWarning:
		return "warning"
	case TypeError:
		return "error"
	case TypeAlert:
		return "alert"
	case TypeEmergency:
		return "emergency"
	}

	return ""
}

// ParseType parses incoming string type representation into byte
func ParseType(s string) (byte, bool) {
	switch s {
	case "trace":
		return TypeTrace, true
	case "debug":
		return TypeDebug, true
	case "info", "notice":
		return TypeInfo, true
	case "warn", "warning":
		return TypeWarning, true
	case "err", "error":
		return TypeError, true
	case "alert", "crititcal":
		return TypeAlert, true
	case "emergency":
		return TypeEmergency, true
	}

	return TypeInfo, false
}
