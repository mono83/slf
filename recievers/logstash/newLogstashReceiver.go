package logstash

import (
	"encoding/json"
	"fmt"
	"github.com/mono83/slf"
	"github.com/mono83/slf/filters"
	"github.com/mono83/udpwriter"
	"io"
	"time"
)

// Config holds information for filtered receiver
type Config struct {
	Address         string   `json:"address" yaml:"address"`
	MinLevel        string   `json:"level" yaml:"level"`
	ParamsWhiteList []string `json:"paramsWhiteList" yaml:"paramsWhiteList"`
	ParamsBlackList []string `json:"paramsBlackList" yaml:"paramsBlackList"`
}

// NewFiltered return new logstash events receiver with
// minimum log level configured
func NewFiltered(c Config) (slf.Receiver, error) {
	// Resolving level
	level, ok := slf.ParseType(c.MinLevel)
	if !ok {
		return nil, fmt.Errorf("Unknown level %s", c.MinLevel)
	}

	out, err := buildReceiverForAddr(c.Address)
	if err != nil {
		return nil, err
	}
	if len(c.ParamsWhiteList) > 0 {
		out.filter = slf.NewWhiteListParamsFilter(c.ParamsWhiteList)
	} else {
		out.filter = slf.NewBlackListParamsFilter(c.ParamsBlackList)
	}

	return filters.MinLogLevel(level, out), nil
}

// New builds and returns new logstash events receiver
func New(addr string) (slf.Receiver, error) {
	return buildReceiverForAddr(addr)
}

func buildReceiverForAddr(addr string) (*logstashLogReceiver, error) {
	uw, err := udpwriter.NewS(addr)
	if err != nil {
		return nil, err
	}

	return &logstashLogReceiver{target: uw, filter: slf.NewBlackListParamsFilter(nil)}, nil
}

type logstashLogReceiver struct {
	target io.Writer
	filter slf.ParamsFilter
}

func (l logstashLogReceiver) Receive(p slf.Event) {
	if !p.IsLog() {
		return
	}

	pkt := map[string]interface{}{}

	if len(p.Params) > 0 {
		shownParams := l.filter(p.Params)
		for _, param := range shownParams {
			value := param.GetRaw()
			if e, ok := value.(error); ok && e != nil {
				value = e.Error()
			}
			pkt[param.GetKey()] = value
		}
	}

	pkt["object"] = p.Marker
	pkt["log-level"] = p.StringType()
	pkt["pattern"] = p.Content
	pkt["message"] = slf.ReplacePlaceholders(p.Content, p.Params, false)
	pkt["event-time"] = p.Time.Format(time.RFC3339)

	bts, _ := json.Marshal(pkt)
	l.target.Write(bts)
}
