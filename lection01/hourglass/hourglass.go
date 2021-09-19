/*
Package hourglass implements a function Hourglass that prints the ASCII image of an hourglass to stdout.
*/
package hourglass

import (
	"fmt"
	"tfs-go-hw/lection01/bashcolor"
)

// Functions returned by GetHourglass.
type (
	SizeSetter            func(size int)
	CharSetter            func(char rune)
	CharColorSetter       func(color bashcolor.Color)
	BackgroundColorSetter func(color bashcolor.Color)
	Hourglass             func()
)

// GetHourglass returns four setters for setting text properties like foreground and background colors and
// the main function for displaying the hourglass.
func GetHourglass() (
	setSize SizeSetter,
	setChar CharSetter,
	setCharColor CharColorSetter,
	setBackgroundColor BackgroundColorSetter,
	hg Hourglass,
) {
	// Default values
	var (
		Size            = 15
		Char            = 'X'
		CharColor       = bashcolor.Blue
		BackgroundColor = bashcolor.Black
	)

	// Setters
	setSize = func(size int) {
		Size = size
	}

	setChar = func(char rune) {
		Char = char
	}

	setCharColor = func(charColor bashcolor.Color) {
		CharColor = charColor
	}

	setBackgroundColor = func(backgroundColor bashcolor.Color) {
		BackgroundColor = backgroundColor
	}
	//

	lineBreak := func() {
		fmt.Println()
	}

	// Main function for drawing Hourglass
	hg = func() {
		charString := fmt.Sprintf("%c", Char)

		printColoredString := func(str string, backgroundColor bashcolor.Color, charColor bashcolor.Color) {
			fmt.Print(bashcolor.Background(backgroundColor), bashcolor.Text(charColor), str, bashcolor.Reset())
		}

		printBase := func() {
			for i := 1; i <= Size; i++ {
				printColoredString(charString, BackgroundColor, CharColor)
			}
		}

		printLine := func(glasses []bool, line int) {
			for column, glass := range glasses {
				if glass {
					printColoredString(charString, BackgroundColor, CharColor)
				} else {
					// Sand at the top
					if 1 <= line && line <= (Size-3)/2 && line+1 <= column && column <= Size-line-2 {
						printColoredString(" ", bashcolor.Yellow, CharColor)
					} else {
						printColoredString(" ", BackgroundColor, CharColor)
					}
				}
			}
		}

		printBase()
		lineBreak()

		chars := make([]bool, Size)

		for line := 1; line <= Size-2; line++ {
			chars[line] = true
			chars[Size-line-1] = true

			printLine(chars, line)

			chars[line] = false
			chars[Size-line-1] = false

			lineBreak()
		}

		printBase()
		lineBreak()
	}

	return setSize, setChar, setCharColor, setBackgroundColor, hg
}
