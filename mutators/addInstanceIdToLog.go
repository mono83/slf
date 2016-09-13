package mutators

import (
	"github.com/mono83/slf"
	"github.com/mono83/slf/rays"
)

// AddInstanceIDToLog mutator injects global InstanceID
// to all log packets
type AddInstanceIDToLog struct{}

// Modify injects instance ID to logs only
func (a AddInstanceIDToLog) Modify(e *slf.Event) {
	if e.IsLog() {
		e.Params = append(e.Params, rays.InstanceID)
	}
}
