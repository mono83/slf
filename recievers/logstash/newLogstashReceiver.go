package logstash

import (
	"encoding/json"
	"fmt"
	"github.com/mono83/slf"
	"github.com/mono83/slf/filters"
	"github.com/mono83/slf/rays"
	"github.com/mono83/udpwriter"
	"io"
	"time"
)

// Config holds information for filtered receiver
type Config struct {
	Address  string `json:"address" yaml:"address"`
	MinLevel string `json:"level" yaml:"level"`
}

// NewFiltered return new logstash events receiver with
// minimum log level configured
func NewFiltered(c Config) (slf.Receiver, error) {
	// Resolving level
	level, ok := slf.ParseType(c.MinLevel)
	if !ok {
		return nil, fmt.Errorf("Unknown level %s", c.MinLevel)
	}

	out, err := New(c.Address)
	if err != nil {
		return nil, err
	}

	return filters.MinLogLevel(level, out), nil
}

// New builds and returns new logstash events receiver
func New(addr string) (slf.Receiver, error) {
	uw, err := udpwriter.NewS(addr)
	if err != nil {
		return nil, err
	}

	return lsr{target: uw}, nil
}

type lsr struct {
	target io.Writer
}

func (l lsr) Receive(p slf.Event) {
	if !p.IsLog() {
		return
	}

	pkt := map[string]interface{}{}

	if len(p.Params) > 0 {
		for _, param := range p.Params {
			value := param.GetRaw()
			if e, ok := value.(error); ok && e != nil {
				value = e.Error()
			}
			pkt[param.GetKey()] = value
		}
	}

	pkt["log-level"] = p.StringType()
	pkt["message"] = p.Content
	pkt["hmessage"] = slf.ReplacePlaceholders(p.Content, p.Params, false)
	pkt["event-time"] = p.Time.Format(time.RFC3339)
	pkt["object"] = p.Marker
	pkt["script-id"] = rays.InstanceID.String()

	bts, _ := json.Marshal(pkt)
	l.target.Write(bts)
}
