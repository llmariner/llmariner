package ui

import "github.com/mgutz/ansi"

var (
	bold    = ansi.ColorFunc("default+b")
	red     = ansi.ColorFunc("red")
	yellow  = ansi.ColorFunc("yellow")
	cyan    = ansi.ColorFunc("cyan")
	green   = ansi.ColorFunc("green")
	redBg   = ansi.ColorFunc("white:red")
	greenBg = ansi.ColorFunc("white:green")
)

// ColorFormatter is a string formatter with color.
type ColorFormatter struct {
	enabled bool
}

// Bold returns bold decorated string.
func (c *ColorFormatter) Bold(s string) string {
	if !c.enabled {
		return s
	}
	return bold(s)
}

// Red returns red-colored string.
func (c *ColorFormatter) Red(s string) string {
	if !c.enabled {
		return s
	}
	return red(s)
}

// Yellow returns yellow-colored string.
func (c *ColorFormatter) Yellow(s string) string {
	if !c.enabled {
		return s
	}
	return yellow(s)
}

// Cyan returns cyan-colored string.
func (c *ColorFormatter) Cyan(s string) string {
	if !c.enabled {
		return s
	}
	return cyan(s)
}

// Green returns green-colored string.
func (c *ColorFormatter) Green(s string) string {
	if !c.enabled {
		return s
	}
	return green(s)
}

// RedBg returns string with red background.
func (c *ColorFormatter) RedBg(s string) string {
	if !c.enabled {
		return s
	}
	return redBg(s)
}

// GreenBg returns string with green background.
func (c *ColorFormatter) GreenBg(s string) string {
	if !c.enabled {
		return s
	}
	return greenBg(s)
}
