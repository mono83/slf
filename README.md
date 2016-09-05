# Simple static logging facility for Golang
[![Build Status](https://travis-ci.org/mono83/slf.svg)](https://travis-ci.org/mono83/slf)
[![Go Report Card](https://goreportcard.com/badge/github.com/mono83/slf)](https://goreportcard.com/report/github.com/mono83/slf)
[![GoDoc](https://godoc.org/github.com/mono83/slf?status.svg)](https://godoc.org/github.com/mono83/slf)

## Usage

In main file, configure general SLF logging behavior:

```
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

```

func AnyFunc() {
    log := wd.New("myMarker" /*sender name (marker)*/, "go." /*metrics prefix*/)
    
    log.Info("Starting processing")
}

```