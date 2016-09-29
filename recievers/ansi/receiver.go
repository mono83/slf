package ansi

import (
	"fmt"
	"github.com/mono83/slf"
	"io"
	"os"
	"strings"
	"time"
)

// Palette
var (
	paletteTime      = makeColor(238, 238)
	paletteInfo      = makeColor(77, 87)
	paletteWarn      = makeColor(226, 208)
	paletteError     = makeColor(202, 196)
	paletteCrit      = makeColor(124, 15)
	paletteTagInfo   = makeColor(37, 37)
	paletteTagError  = makeColor(166, 166)
	paletteTagCrit   = makeColor(196, 196)
	paletteTagMarker = makeColor(91, 91)
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

func (a *ansiPrinter) timeRef(e slf.Event) string {
	var st string
	if delta := e.Time.Sub(a.previousShown); delta < time.Second {
		st = fmt.Sprintf("+ %.6fs", e.Time.Sub(a.previous).Seconds())
	} else {
		st = e.Time.Format("02 15:04:05")
		a.previousShown = e.Time
	}

	return st
}

func (a *ansiPrinter) print(e slf.Event) {
	if !e.IsLog() {
		return
	}

	// Preparing tag
	var tag string
	var tagFunc, msgFunc color
	switch e.Type {
	case slf.TypeTrace:
		tag = " ┊ "
		tagFunc = paletteTime
		msgFunc = paletteTime
	case slf.TypeDebug:
		tag = " ┇ "
		tagFunc = paletteTime
		msgFunc = paletteTime
	case slf.TypeInfo:
		tag = " ┃ "
		tagFunc = paletteTagInfo
		msgFunc = paletteInfo
	case slf.TypeWarning:
		tag = " ┃ "
		tagFunc = paletteTagError
		msgFunc = paletteWarn
	case slf.TypeError:
		tag = " ┣ "
		tagFunc = paletteTagError
		msgFunc = paletteError
	case slf.TypeAlert:
		tag = " ◈ "
		tagFunc = paletteTagCrit
		msgFunc = paletteCrit
	case slf.TypeEmergency:
		tag = " ◉ "
		tagFunc = paletteTagCrit
		msgFunc = paletteCrit
	}

	// Replacing placeholders
	text := getText(e, a.colors, msgFunc)

	// Replacing special chars
	text = strings.Replace(text, "\t", "\\t", -1)
	text = strings.Replace(text, "\n", "\\n", -1)
	text = strings.Replace(text, "\r", "\\r", -1)

	marker := ""
	if a.showMarker {
		marker = paletteTagMarker.format(a.colors, e.Marker) + " "
	}

	stop := ""
	if a.colors {
		stop = ansiStop
	}
	fmt.Fprintf(
		a.to,
		"%s %s %s%s%s\n",
		paletteTime.format(a.colors, a.timeRef(e)),
		tagFunc.format(a.colors, tag),
		marker,
		msgFunc.format(a.colors, text),
		stop,
	)
	a.previous = e.Time
}

func getText(e slf.Event, colors bool, msgFunc color) string {
	text := e.Content
	if len(e.Params) > 0 {
		if colors {
			// Building map
			mp := make(map[string]slf.Param, len(e.Params))
			for _, p := range e.Params {
				mp[p.GetKey()] = p
			}

			text = slf.PlaceholdersRegex.ReplaceAllStringFunc(text, func(x string) string {
				key := x[1:]
				if v, ok := mp[key]; ok {
					if ev, ok := v.GetRaw().(error); ok {
						return msgFunc.formatErr(colors, ev)
					}
					return msgFunc.formatVar(colors, v.String())
				}
				return "<!" + x + ">"
			})
		} else {
			text = slf.ReplacePlaceholders(text, e.Params, true)
		}
	}

	return text
}
