package log

import (
	"fmt"
	"github.com/mono83/slf/wd"
	"io"
	"os"
)

var stdWd = wd.New("go", "")

// SetOutput is drop-in replacement for standard log package analogue
func SetOutput(io.Writer) {}

// Flags is drop-in replacement for standard log package analogue
func Flags() int { return 0 }

// SetFlags is drop-in replacement for standard log package analogue
func SetFlags(int) {}

// Prefix is is drop-in replacement for standard log package analogue
func Prefix() string { return "" }

// SetPrefix is drop-in replacement for standard log package analogue
func SetPrefix(string) {}

// Print is drop-in replacement for standard log package analogue
func Print(v ...interface{}) { stdWd.Info(fmt.Sprint(v...)) }

// Printf is drop-in replacement for standard log package analogue
func Printf(format string, v ...interface{}) { stdWd.Info(fmt.Sprintf(format, v...)) }

// Println is drop-in replacement for standard log package analogue
func Println(v ...interface{}) { stdWd.Info(fmt.Sprint(v...)) }

// Fatal is drop-in replacement for standard log package analogue
func Fatal(v ...interface{}) { stdWd.Alert(fmt.Sprint(v...)); os.Exit(1) }

// Fatalf is drop-in replacement for standard log package analogue
func Fatalf(format string, v ...interface{}) { stdWd.Alert(fmt.Sprintf(format, v...)); os.Exit(1) }

// Fatalln is drop-in replacement for standard log package analogue
func Fatalln(v ...interface{}) { stdWd.Alert(fmt.Sprint(v...)); os.Exit(1) }

// Panic is drop-in replacement for standard log package analogue
func Panic(v ...interface{}) { stdWd.Error(fmt.Sprint(v...)); panic(fmt.Sprint(v...)) }

// Panicf is drop-in replacement for standard log package analogue
func Panicf(format string, v ...interface{}) {
	stdWd.Error(fmt.Sprintf(format, v...))
	panic(fmt.Sprintf(format, v...))
}

// Panicln is drop-in replacement for standard log package analogue
func Panicln(v ...interface{}) { stdWd.Error(fmt.Sprint(v...)); panic(fmt.Sprint(v...)) }
