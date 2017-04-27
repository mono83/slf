package statsd

import (
	"github.com/mono83/slf"
	"github.com/mono83/slf/recievers"
	"github.com/mono83/udpwriter"
	"io"
	"time"
)

// Config holds configuration for StatsD client
type Config struct {
	Address       string   `json:"address" yaml:"address"`
	Microseconds  bool     `json:"microseconds" yaml:"microseconds"`
	Dogstats      bool     `json:"dogstats" yaml:"dogstats"`
	AllowedParams []string `json:"params" yaml:"params"`
	Prefix        string   `json:"prefix" yaml:"prefix"`
	FlushEvery    int      `json:"flushEvery" yaml:"flushEvery"`
	PacketLines   int      `json:"packetLines" yaml:"packetLines"`
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
	s.allowed = map[string]bool{}
	for _, v := range c.AllowedParams {
		s.allowed[v] = true
	}

	return s, nil
}

type sdr struct {
	slf.Receiver
	target  io.Writer
	config  Config
	allowed map[string]bool
}

func (s *sdr) sendToWriter(list []slf.Event) {
	if len(list) == 0 {
		return
	}

	startIndex := 0
	precision := int64(1000000)
	if s.config.Microseconds {
		precision = 1000
	}
	var buf PacketBuilder
	for startIndex+1 <= len(list) {
		endIndex := startIndex + s.config.PacketLines
		if endIndex >= len(list) {
			endIndex = len(list)
		}

		part := list[startIndex:endIndex]
		if s.config.Dogstats {
			buf = NewPacketWithParamsBuilder(s.config.Prefix, precision, s.allowed)
		} else {
			buf = NewPacketBuilder(s.config.Prefix, precision)
		}

		for _, p := range part {
			buf.WriteEvent(p)
		}

		s.target.Write(buf.Bytes())
		startIndex = endIndex
	}
}
