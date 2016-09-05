package wd

import (
	"github.com/mono83/slf"
	"time"
)

func newLog(t byte, mark, msg string, p []slf.Param) slf.Event {
	return slf.Event{
		Time:    time.Now(),
		Marker:  mark,
		Type:    t,
		Content: msg,
		Params:  p,
	}
}

func newMetrics(t byte, mark, msg string, i int64, p []slf.Param) slf.Event {
	return slf.Event{
		Time:    time.Now(),
		Marker:  mark,
		Type:    t,
		Content: msg,
		Params:  p,
		I64:     i,
	}
}
