package statsd

import (
	"bytes"
	"fmt"
	"github.com/mono83/slf"
	"github.com/mono83/slf/recievers"
	"github.com/mono83/udpwriter"
	"io"
	"time"
)

// Config holds configuration for StatsD client
type Config struct {
	Address      string `json:"address" yaml:"address"`
	Microseconds bool   `json:"microseconds" yaml:"microseconds"`
	Dogstats     bool   `json:"dogstats" yaml:"dogstats"`
	Prefix       string `json:"prefix" yaml:"prefix"`
	FlushEvery   int    `json:"flushEvery" yaml:"flushEvery"`
	PacketLines  int    `json:"packetLines" yaml:"packetLines"`
}

// NewReceiver builds new StatsD receiver
func NewReceiver(c Config) (slf.Receiver, error) {
	uw, err := udpwriter.NewS(c.Address)
	if err != nil {
		return nil, err
	}

	if c.PacketLines < 5 {
		c.PacketLines = 5
	}

	if c.FlushEvery < 1 {
		c.FlushEvery = 1
	}

	s := &sdr{target: uw, config: c}
	s.Receiver = recievers.NewMetricsBuffer(time.Duration(int64(c.FlushEvery)*int64(time.Second)), s.sendToWriter)

	return s, nil
}

type sdr struct {
	slf.Receiver
	target io.Writer
	config Config
}

func (s *sdr) sendToWriter(list []slf.Event) {
	if len(list) == 0 {
		return
	}

	startIndex := 0
	for startIndex+1 <= len(list) {
		endIndex := startIndex + s.config.PacketLines
		if endIndex >= len(list) {
			endIndex = len(list)
		}

		part := list[startIndex:endIndex]
		buf := bytes.NewBuffer(nil)

		for _, p := range part {
			switch p.Type {
			case slf.TypeInc:
				buf.Write([]byte(fmt.Sprintf(
					"%s%s:%d|c%s\n",
					s.config.Prefix,
					p.Content,
					p.I64,
					s.appendTags(p.Params),
				)))
			case slf.TypeGauge:
				buf.Write([]byte(fmt.Sprintf(
					"%s%s:%d|g%s\n",
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

				buf.Write([]byte(fmt.Sprintf(
					"%s%s:%d|ms%s\n",
					s.config.Prefix,
					p.Content,
					d,
					s.appendTags(p.Params),
				)))
			}
		}

		s.target.Write(buf.Bytes())
		startIndex = endIndex
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
		sb.WriteString(tag.GetKey())
		sb.WriteString(":")
		sb.WriteString(tag.String())
	}

	return sb.String()
}
