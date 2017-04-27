package statsd

import (
	"bytes"
	"github.com/mono83/slf"
	"strconv"
	"strings"
)

// NewPacketBuilder builds packet builder without params support
func NewPacketBuilder(prefix string, precision int64) PacketBuilder {
	// Preparing prefix
	if len(prefix) > 0 && prefix[len(prefix)-1] != '.' {
		prefix += "."
	}

	return &packetBuilder{
		Buffer:        bytes.NewBuffer(nil),
		prefix:        []byte(prefix),
		precision:     precision,
		paramsAllowed: false,
		wasParams:     false,
	}
}

// NewPacketWithParamsBuilder builds packet builder with params support
func NewPacketWithParamsBuilder(prefix string, precision int64, allowed map[string]bool) PacketBuilder {
	// Preparing prefix
	if len(prefix) > 0 && prefix[len(prefix)-1] != '.' {
		prefix += "."
	}

	return &packetBuilder{
		Buffer:        bytes.NewBuffer(nil),
		prefix:        []byte(prefix),
		precision:     precision,
		paramsAllowed: true,
		paramsMap:     allowed,
		wasParams:     false,
	}
}

// PacketBuilder is special type of buffer, used to build StatsD-compatible packets.
type PacketBuilder interface {
	WriteEvent(e slf.Event)
	Size() int
	Bytes() []byte
}

type packetBuilder struct {
	*bytes.Buffer
	count         int             // Amount of placed events
	prefix        []byte          // Metrics prefix
	precision     int64           // 1 for nanos, 1000 for micros and 1000000 for millis (for duration events)
	paramsAllowed bool            // True if params printing in Dogstats format allowed
	paramsMap     map[string]bool // Whitelist of allowed params
	wasParams     bool            // True if params output was started at current line
}

func (pb *packetBuilder) Size() int {
	return pb.count
}

func (pb *packetBuilder) WriteEvent(e slf.Event) {
	if len(e.Content) > 0 {
		if pb.Size() > 0 {
			pb.WriteRune('\n')
		}
		if len(pb.prefix) > 0 {
			pb.Write(pb.prefix)
		}
		pb.WriteString(e.Content)
		pb.WriteRune(':')
		switch e.Type {
		case slf.TypeGauge:
			pb.WriteString(strconv.FormatInt(e.I64, 10))
			pb.WriteRune('|')
			pb.WriteRune('g')
		case slf.TypeDuration:
			pb.WriteString(strconv.FormatInt(e.I64/pb.precision, 10))
			pb.WriteString("|ms")
		default:
			pb.WriteString(strconv.FormatInt(e.I64, 10))
			pb.WriteRune('|')
			pb.WriteRune('c')
		}

		if pb.paramsAllowed && len(e.Params) > 0 {
			for _, param := range e.Params {
				if pb.paramsMap[param.GetKey()] {
					pb.WriteParam(param)
				}
			}
		}
		pb.count++
	}
}

func (pb *packetBuilder) WriteParam(p slf.Param) {
	if p != nil {
		if pb.wasParams {
			pb.WriteRune(',')
		} else {
			pb.WriteString("|@1.0|#")
			pb.wasParams = true
		}

		pb.WriteString(p.GetKey())
		pb.WriteRune(':')
		pb.Write(SanitizeParamValue(p.String()))
	}
}

var sanitizeReplacement = byte('_')

// SanitizeParamValue replaces special characters from param value
func SanitizeParamValue(value string) []byte {
	if len(value) == 0 {
		return []byte{}
	}

	bts := []byte(strings.TrimSpace(value))
	for i, v := range bts {
		if !(v == 46 || (v >= 48 && v <= 57) || (v >= 65 && v <= 90) || (v >= 97 && v <= 122)) {
			bts[i] = sanitizeReplacement
		}
	}

	return bts
}
