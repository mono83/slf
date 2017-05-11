package context

import (
	"github.com/mono83/slf/rays"
	"github.com/mono83/slf/wd"
	"time"
)

// New builds new pooled context
func New(marker, metricsPrefix string) Interface {
	ctx := new(pooled)
	ctx.id = rays.New()
	ctx.start = time.Now().UTC()
	ctx.marker = marker
	ctx.prefix = metricsPrefix
	ctx.createLog()
	register(ctx)
	return ctx
}

type pooled struct {
	id             rays.RayID
	parent         *pooled
	marker, prefix string
	start, done    time.Time
	wd.Watchdog
}

func (p *pooled) createLog() { p.Watchdog = wd.New(p.marker, p.prefix).WithParams(p.id) }

func (p *pooled) ID() rays.RayID       { return p.id }
func (p *pooled) StartedAt() time.Time { return p.start }

func (p *pooled) Close() {
	if p.parent != nil {
		p.parent.Close()
		return
	}
	if p.done.IsZero() {
		p.done = time.Now().UTC()
		unregister(p)
	}
}

func (p *pooled) WithMarker(marker string) Interface {
	return p.With(marker, p.prefix)
}

func (p *pooled) With(marker, metricsPrefix string) Interface {
	if p.marker == marker && p.prefix == metricsPrefix {
		// No changes
		return p
	}

	ctx := new(pooled)
	ctx.id = p.id
	ctx.parent = p
	ctx.start = p.start
	ctx.done = p.done
	ctx.marker = marker
	ctx.prefix = metricsPrefix
	ctx.createLog()
	return ctx
}
