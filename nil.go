package slf

import "time"

// Nil contains No-Op logger, receiver and stats reporter
var Nil = nilstruct{}

type nilstruct struct{}

func (nilstruct) Receive(e Event)                                      {}
func (nilstruct) Trace(message string, p ...Param)                     {}
func (nilstruct) Debug(message string, p ...Param)                     {}
func (nilstruct) Info(message string, p ...Param)                      {}
func (nilstruct) Warning(message string, p ...Param)                   {}
func (nilstruct) Error(message string, p ...Param)                     {}
func (nilstruct) Alert(message string, p ...Param)                     {}
func (nilstruct) Emergency(message string, p ...Param)                 {}
func (nilstruct) IncCounter(name string, value int64, p ...Param)      {}
func (nilstruct) UpdateGauge(name string, value int64, p ...Param)     {}
func (nilstruct) RecordTimer(name string, d time.Duration, p ...Param) {}
func (n nilstruct) Timer(name string, p ...Param) Timer                { return NewTimer(name, p, n) }
