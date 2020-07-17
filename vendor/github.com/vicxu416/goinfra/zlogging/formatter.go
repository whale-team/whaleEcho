package zlogging

import (
	"fmt"
	"strings"
)

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorBold     = 1
	colorDarkGray = 90
)

var levelColorMap = map[string]int{
	"trace": colorMagenta,
	"debug": colorYellow,
	"info":  colorGreen,
	"warn":  colorRed,
	"error": colorRed,
	"fatal": colorRed,
	"panic": colorRed,
}

func colorize(s interface{}, c int) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

// LevelFormater used to formating level
func LevelFormater(i interface{}) string {
	var (
		msg      string
		level    string
		ok       bool
		color    int
		needBold = false
	)

	if level, ok = i.(string); ok {
		if color, ok = levelColorMap[level]; !ok {
			level = "???"
			color = colorBold
		}
		if level == "error" || level == "fatal" || level == "panic" {
			needBold = true
		}
	} else if i == nil {
		level = "???"
		color = colorBold
	} else {
		level = strings.ToUpper(fmt.Sprintf("%s", i))[0:3]
		color = colorBold
	}

	msg = strings.ToUpper(level)
	msg = colorize(msg, color)

	if needBold {
		msg = colorize(msg, colorBold)
	}

	return msg
}
