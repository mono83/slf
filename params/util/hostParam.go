package util

import (
	"github.com/mono83/slf"
	"github.com/mono83/slf/params"
	"os"
)

// HostParam return StringParam with current hostname as value and "host" as key
// If app was unable to determine hostname, NilParam will be returned
func HostParam() slf.Param {
	if host == "" {
		return params.Nil{Key: "host"}
	}

	return params.String{Key: "host", Value: host}
}

var host = ""

func init() {
	// Reading host
	if h, err := os.Hostname(); err == nil {
		host = h
	}
}
