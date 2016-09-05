package ansi

import (
	"fmt"
	"github.com/mono83/slf"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

var placeholdersRegex = regexp.MustCompile(":[0-9a-zA-Z\\-_]+")

func ansi256(c int) func(string) string {
	sc := fmt.Sprintf("%d", c)
	return func(s string) string {
		return "\033[38;5;" + sc + "m" + s + "\033[0m"
	}
}

// Palette
var (
	noColor          = func(v string) string { return v }
	paletteTime      = ansi256(238)
	paletteInfo      = ansi256(81)
	paletteWarn      = ansi256(226)
	paletteError     = ansi256(196)
	paletteTagInfo   = ansi256(37)
	paletteTagError  = ansi256(166)
	paletteTagCrit   = ansi256(196)
	paletteTagMarker = ansi256(91)
)

// New builds new stdout receiver with 256 ANSI colors support
func New(colors, showMarker, async bool) slf.Receiver {
	a := new(ansiPrinter)
	a.colors = colors
	a.showMarker = showMarker
	a.to = os.Stdout
	if async {
		a.delay = make(chan slf.Event)
		go func() {
			for e := range a.delay {
				a.print(e)
			}
		}()
	}

	return a
}

type ansiPrinter struct {
	to                      io.Writer
	colors, showMarker      bool
	delay                   chan slf.Event
	previous, previousShown time.Time
}

func (a *ansiPrinter) Receive(e slf.Event) {
	if a.delay == nil {
		a.print(e)
	} else {
		a.delay <- e
	}
}

func (a ansiPrinter) sel(colorFunc func(string) string) func(string) string {
	if a.colors {
		return colorFunc
	}

	return noColor
}

func (a *ansiPrinter) print(e slf.Event) {
	if !e.IsLog() {
		return
	}

	// Formatting time
	var st string
	if delta := e.Time.Sub(a.previousShown); delta < time.Second {
		st = fmt.Sprintf("+ %.6fs", e.Time.Sub(a.previous).Seconds())
	} else {
		st = e.Time.Format("02 15:04:05")
		a.previousShown = e.Time
	}

	// Preparing tag
	var tag string
	var tagFunc = noColor
	var msgFunc = noColor
	switch e.Type {
	case slf.TypeDebug:
		tag = " ▪ "
		tagFunc = paletteTime
		msgFunc = paletteTime
	case slf.TypeInfo:
		tag = " ▪ "
		tagFunc = paletteTagInfo
		msgFunc = paletteInfo
	case slf.TypeWarning:
		tag = "▪▪▪"
		tagFunc = paletteTagError
		msgFunc = paletteWarn
	case slf.TypeError:
		tag = " ✗ "
		tagFunc = paletteTagError
		msgFunc = paletteError
	case slf.TypeAlert, slf.TypeEmergency:
		tag = "✗✗✗"
		tagFunc = paletteTagCrit
		msgFunc = paletteError
	default:
		tag = "   "
		msgFunc = paletteTime
	}

	// Replacing placeholders
	text := e.Content
	if len(e.Params) > 0 {
		// Building map
		mp := make(map[string]slf.Param, len(e.Params))
		for _, p := range e.Params {
			mp[p.GetKey()] = p
		}

		text = placeholdersRegex.ReplaceAllStringFunc(text, func(x string) string {
			key := x[1:]
			if v, ok := mp[key]; ok {
				return fmt.Sprintf("[%v]", v.GetRaw())
			}
			return "<!" + x + ">"
		})
	}

	// Replacing special chars
	text = strings.Replace(text, "\t", "\\t", -1)
	text = strings.Replace(text, "\n", "\\n", -1)
	text = strings.Replace(text, "\r", "\\r", -1)

	marker := ""
	if a.showMarker {
		marker = a.sel(paletteTagMarker)(e.Marker) + " "
	}

	fmt.Fprintf(
		a.to,
		"%s %s %s%s\n",
		a.sel(paletteTime)(st),
		a.sel(tagFunc)(tag),
		marker,
		a.sel(msgFunc)(text),
	)
	a.previous = e.Time
}
