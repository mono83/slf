package mutators

import "github.com/mono83/slf"

// HideInfo is mutator, that converts INFO log levels
// to DEBUG
type HideInfo struct {
	Predicate func(slf.Event) bool
}

// Modify changes log level from INFO to DEBUG
func (h HideInfo) Modify(e *slf.Event) {
	if h.Predicate != nil && !h.Predicate(*e) {
		return
	}

	if e.Type == slf.TypeInfo {
		e.Type = slf.TypeDebug
	}
}

// HideInfoForMarker builds HideInfo mutator, than works
// only with log events, marked with particular marker
func HideInfoForMarker(markerName string) HideInfo {
	return HideInfo{
		Predicate: func(e slf.Event) bool {
			return e.Marker == markerName
		},
	}
}
