/*
Package bashcolor provides work with the console output color.
*/
package bashcolor

import (
	"fmt"
)

// Color describes the color of the console output.
type (
	Color int
)

// Supported colors.
const (
	Black = Color(iota)
	Red
	Green
	Yellow
	Blue
	Purple
	Cyan
	LightGray
)

// Text returns the beginning of the color modification.
func Text(color Color) string {
	return fmt.Sprintf("\033[3%dm", color)
}

// Background returns the beginning of the color modification.
func Background(color Color) string {
	return fmt.Sprintf("\033[4%dm", color)
}

// Reset returns the end of the color modification.
func Reset() string {
	return "\033[0m"
}
