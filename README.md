# Simple static logging facility for Golang
[![Build Status](https://travis-ci.org/mono83/slf.svg)](https://travis-ci.org/mono83/slf)
[![Go Report Card](https://goreportcard.com/badge/github.com/mono83/slf)](https://goreportcard.com/report/github.com/mono83/slf)
[![GoDoc](https://godoc.org/github.com/mono83/slf?status.svg)](https://godoc.org/github.com/mono83/slf)

## Usage

In main file, configure general SLF logging behavior:

```go
package main

import (
    "github.com/mono83/slf/wd"
    "github.com/mono83/slf/recievers/ansi"
)

func main() {
    // Initialize logger
    wd.AddReceiver(ansi.New(true /*colors*/, true /*show marker*/, false /*async*/))
}

```

Then, create new Watchdog instance and use logging

```go

func AnyFunc() {
    log := wd.New("myMarker" /*sender name (marker)*/, "go." /*metrics prefix*/)
    
    log.Info("Starting processing")
}

```

You can use placeholders 

```go
log.Info("Sending request to :url", wd.StringParam("url", url))
log.Error("Received an error - :err", wd.ErrorParam(err))
```

## Params

SLF uses *Params* (aka placeholders in logging context) to add context into SLF events

1. In logging, params acts as placeholders and it's value injects into message string
2. In special loggers, params can be also delivered as part of context - for example all of params are delivered to logstash server in logstash client (`github.com/mono83/slf/recievers/logstash.New`)
3. In metrics params may act as context to, if client allows it (DogStatsD, for example, allows context delivery)

All params are simple key-value pairs, where key is always `string`

All params can be found under `github.com/mono83/slf/params` package. For more comfortable use, there are builder function in `github.com/mono83/slf/wd` package:

* `wd.IntParam(name string, value int)` - returns integer param
* `wd.CountParam(value int)` - return integer param with name `count`
* `wd.ErrParam(err error)` - returns error param with name `err`
* `wd.StringParam(name, value string)` - returns string param
* `wd.NameParam(value string)` - returns string param with name `name`
* `wd.FloatParam(name string, value float64)` - returns float param

## Useful stuff

### Log to metrics

`github.com/mono83/slf/health` package contains func `StartLogToMetrics`, which provides simple interface to
count logging events without loosing them 

Usecase:

```go
StartLogToMetrics(wd.New("", "log." /*metrics prefix*/).WithParams(util.HostParam()))
```

### Collect Go application health

`github.com/mono83/slf/health.StartHealthMonitor` starts health monitoring goroutine, that will automatically collect various memory and GC information every second and send it as gauge

Usecase:

```go
StartHealthMonitor(wd.New("", "app." /*metrics prefix*/).WithParams(util.HostParam()))
```

Metrics taken:

* `health.gcs` - GCs count
* `health.goroutines` - amount of active goroutines
* `health.sys.malloc`
* `health.sys.free`
* `health.heap.alloc`
* `health.heap.inuse`
* `health.heap.sys`
* `health.heap.objects` - count of objects in heap
* `health.heap.nextgc` - next GC threshold in bytes
* `health.uptime` - uptime in seconds
