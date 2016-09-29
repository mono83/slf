package ansi

import (
	"strconv"
)

var ansiStop = "\033[0m"

func makeColor(frontColor, varColor int) color {
	return color{
		value:    frontColor,
		fgColor:  strconv.Itoa(frontColor),
		varColor: strconv.Itoa(varColor),
	}
}

type color struct {
	value    int
	fgColor  string
	varColor string
}

func (c color) format(colors bool, value string) string {
	if !colors || c.value == 0 {
		return value
	}

	return "\033[38;5;" + c.fgColor + "m" + value
}

func (c color) formatVar(colors bool, value string) string {
	if !colors || c.value == 0 {
		return "[" + value + "]"
	}

	return "\033[38;5;" + c.varColor + "m" + value + ansiStop + "\033[38;5;" + c.fgColor + "m"
}

func (c color) formatErr(colors bool, value error) string {
	if !colors || c.value == 0 {
		return "[" + value.Error() + "]"
	}

	return "\033[48;5;88m\033[97m" + value.Error() + ansiStop + "\033[38;5;" + c.fgColor + "m"
}
