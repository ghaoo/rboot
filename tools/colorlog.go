package tools

import "fmt"

// Code color code.
type Code int

// None color code.
const None Code = -1

// Attributes.
const (
	Reset Code = iota
	Bold
	Dim
	Italic
	Underline
	Blink
	BlinkFast
	Inverse
	Hidden
	Strikethrough
)

// Foreground colors.
const (
	Black Code = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	LightGray
)

// Foreground colors with higher offset.
const (
	DarkGray Code = iota + 90
	LightRed
	LightGreen
	LightYellow
	LightBlue
	LightMagenta
	LightCyan
	White
)

// Background colors.
const (
	BgBlack Code = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgLightGray
)

// Background colors with higher offset.
const (
	BgDarkGray Code = iota + 100
	BgLightRed
	BgLightGreen
	BgLightYellow
	BgLightBlue
	BgLightMagenta
	BgLightCyan
	BgWhite
)

// String sequence for code.
func (c Code) String() string {
	return fmt.Sprintf("\x1b[%dm", c)
}

// Fg256 foreground 256 colors
func Fg256(c uint8) string {
	return fmt.Sprintf("\x1b[38;5;%dm", c)
}

// Bg256 background 256 colors
func Bg256(c uint8) string {
	return fmt.Sprintf("\x1b[48;5;%dm", c)
}
