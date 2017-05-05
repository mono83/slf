package writer

import (
	"fmt"
	"github.com/mono83/slf"
	"io"
	"os"
	"strings"
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

// Options contains configuration for writer receiver
type Options struct {
	Target     io.Writer // Output target, nil for os.Stdout
	NoColor    bool      // Disable color output
	Async      bool      // Async output
	Marker     bool      // Output log marker
	TimeFormat string    // Time output format, default is "02 15:04:05.000000"
}

// New builds new writer receiver with 256 ANSI colors support
func New(o Options) slf.Receiver {
	a := new(printer)
	a.Options = o
	if a.Target == nil {
		a.Target = os.Stdout
	}
	if len(a.TimeFormat) == 0 {
		a.TimeFormat = "02 15:04:05.000000"
	}
	if a.Async {
		a.delay = make(chan slf.Event)
		go func() {
			for e := range a.delay {
				a.print(e)
			}
		}()
	}

	return a
}

type printer struct {
	Options
	delay chan slf.Event
}

func (p *printer) Receive(e slf.Event) {
	if p.Async {
		p.delay <- e
	} else {
		p.print(e)
	}
}

func (p *printer) timeRef(e slf.Event) string {
	return e.Time.Format(p.TimeFormat)
}

func (p *printer) print(e slf.Event) {
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
	text := getText(e, !p.NoColor, msgFunc)

	// Replacing special chars
	text = strings.Replace(text, "\t", "\\t", -1)
	text = strings.Replace(text, "\n", "\\n", -1)
	text = strings.Replace(text, "\r", "\\r", -1)

	marker := ""
	if p.Marker {
		marker = " " + paletteTagMarker.format(!p.NoColor, "@"+e.Marker)
	}

	stop := ""
	if p.Marker {
		stop = ansiStop
	}
	fmt.Fprintf(
		p.Target,
		"%s%s%s%s%s\n",
		paletteTime.format(!p.NoColor, p.timeRef(e)),
		tagFunc.format(!p.NoColor, tag),
		msgFunc.format(!p.NoColor, text),
		marker,
		stop,
	)
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
