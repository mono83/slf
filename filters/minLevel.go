package filters

import "github.com/mono83/slf"

// MinLogLevel return filtering receiver with min log level threshold
func MinLogLevel(minLevel byte, r slf.Receiver) slf.Receiver {
	return slf.Filter(r, func(e slf.Event) bool {
		return e.IsLog() && e.Type >= minLevel
	})
}
