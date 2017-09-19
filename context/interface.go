package context

import (
	"github.com/mono83/slf/rays"
	"github.com/mono83/slf/wd"
	"time"
)

// Interface describes application execution context
type Interface interface {
	wd.Watchdog

	// Closes context and marks it as finished
	Close()

	// StartedAt returns time, this context was created
	StartedAt() time.Time

	// GetID returns context ID, unique for each one
	ID() rays.RayID

	// WithMarker returns clone of context (with same ID), but with new marker
	WithMarker(marker string) Interface

	// With returns clone of context (with same ID), but with new marker and metrics prefix
	With(marker, metricsPrefix string) Interface
}
