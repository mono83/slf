# Simple static logging facility for Golang
[![Build Status](https://travis-ci.org/mono83/slf.svg)](https://travis-ci.org/mono83/slf)
[![Go Report Card](https://goreportcard.com/badge/github.com/mono83/slf)](https://goreportcard.com/report/github.com/mono83/slf)
[![GoDoc](https://godoc.org/github.com/mono83/slf?status.svg)](https://godoc.org/github.com/mono83/slf)

## Brief

SLF provides easy and convenient way to organize logging and metrics reporting inside your application.

```go
import (
    "github.com/mono83/slf/wd"
)

// Create logger with marker example
log := wd.NewLogger("example")
log.Info("Processing some staff")
```

Allowed logging levels are: `Trace`, `Debug`, `Info`, `Warning`, `Error`, `Alert`

## Got logs but no output? Assign receivers!

By default, SLF will ignore any received. To manage output or forwarding to external logging facilities just add
a receiver. Most common is stdout receiver, which can be assigned right in your `func main()`

```go

import (
    "github.com/mono83/slf/wd"
    "github.com/mono83/slf/recievers/ansi"
)

func main() {
        // Add ANSI standard output receiver 
        wd.AddReceiver(ansi.New(true /*colors*/, true /*show marker*/, false /*async*/))
}
```

## More concision with placeholders

SLF supports additional parameters, that can be applied to log. These parameters can be used as placeholders values and 
some receiver can contain additional logic for it (for example, Logstash receiver can send placeholder values as separated
columns):
 
```go
log.Info("Processing user :id", params.Int{Key: "id", Value: 300})
```

Raw params aren't most convenient things in the world, so `wd.` package contains some builders for frequent params:

| Func | Args | Description |
| ---- | ---- | ----------- |
|`wd.IntParam` | `string`, `int` | |
|`wd.Int64Param` | `string`, `int` | |
|`wd.CountParam` | `int` | Builds integer param with key `count` |
|`wd.ID64Param` | `int` | Builds 64-bit integer param with key `id` |
|`wd.FloatParam` | `string`, `float64` | |
|`wd.ErrParam` | `error` | Builds param containing Go `error`s with key `err` |
|`wd.StringParam` | `string`, `string` | |
|`wd.NameParam` | `string`, `strings` | Builds string param with key `name` |

