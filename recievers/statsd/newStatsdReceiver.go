package statsd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mono83/slf"
	"github.com/mono83/udpwriter"
	"io"
)

// Config holds configuration for StatsD client
type Config struct {
	Address      string `json:"address" yaml:"address"`
	Microseconds bool   `json:"microseconds" yaml:"microseconds"`
	Buffered     bool   `json:"buffered" yaml:"buffered"`
	Dogstats     bool   `json:"dogstats" yaml:"dogstats"`
	Prefix       string `json:"prefix" yaml:"prefix"`
}

// NewReceiver builds new StatsD receiver
func NewReceiver(c Config) (slf.Receiver, error) {
	uw, err := udpwriter.NewS(c.Address)
	if err != nil {
		return nil, err
	}

	if c.Buffered {
		return nil, errors.New("Buffered output not supported yet")
	}

	return &sdr{target: uw, config: c}, nil
}

type sdr struct {
	target io.Writer
	config Config
}

func (s *sdr) Receive(p slf.Event) {
	switch p.Type {
	case slf.TypeInc:
		s.target.Write([]byte(fmt.Sprintf(
			"%s%s:%d|c%s",
			s.config.Prefix,
			p.Content,
			p.I64,
			s.appendTags(p.Params),
		)))
	case slf.TypeGauge:
		s.target.Write([]byte(fmt.Sprintf(
			"%s%s:%d|g%s",
			s.config.Prefix,
			p.Content,
			p.I64,
			s.appendTags(p.Params),
		)))
	case slf.TypeDuration:
		var d int64
		if s.config.Microseconds {
			// Microseconds precision
			d = p.I64 / 1000
		} else {
			// Milliseconds precision
			d = p.I64 / 1000000
		}

		s.target.Write([]byte(fmt.Sprintf(
			"%s%s:%d|ms%s",
			s.config.Prefix,
			p.Content,
			d,
			s.appendTags(p.Params),
		)))
	}
}

func (s sdr) appendTags(tags []slf.Param) string {
	if len(tags) == 0 || !s.config.Dogstats {
		return ""
	}

	sb := bytes.NewBufferString("|@1.0|#")
	for i, tag := range tags {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("%s:%v", tag.GetKey(), tag.GetRaw()))
	}

	return sb.String()
}
