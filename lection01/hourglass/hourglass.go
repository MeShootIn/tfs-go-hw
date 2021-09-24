/*
Package hourglass implements functions for displaying the ASCII image of an hourglass to stdout.
*/
package hourglass

import (
	"fmt"
	"tfs-go-hw/lection01/bashcolor"
)

type (
	// ParamStorage is a function that "stores" output parameters
	ParamStorage func() (int, rune, bashcolor.Color, bashcolor.Color)
)

// displayHourglass is a main function for displaying hourglass.
func displayHourglass(size int, char rune, charColor bashcolor.Color, backgroundColor bashcolor.Color) {
	charString := fmt.Sprintf("%c", char)

	lineBreak := func() {
		fmt.Println()
	}

	printColoredString := func(str string, backgroundColor bashcolor.Color, charColor bashcolor.Color) {
		fmt.Print(bashcolor.Background(backgroundColor), bashcolor.Text(charColor), str, bashcolor.Reset())
	}

	printBase := func() {
		for i := 1; i <= size; i++ {
			printColoredString(charString, backgroundColor, charColor)
		}
	}

	printLine := func(glasses []bool, line int) {
		for column, glass := range glasses {
			if glass {
				printColoredString(charString, backgroundColor, charColor)
			} else {
				// Sand at the top
				if 1 <= line && line <= (size-3)/2 && line+1 <= column && column <= size-line-2 {
					printColoredString(" ", bashcolor.Yellow, charColor)
				} else {
					printColoredString(" ", backgroundColor, charColor)
				}
			}
		}
	}

	printBase()
	lineBreak()

	chars := make([]bool, size)

	for line := 1; line <= size-2; line++ {
		chars[line] = true
		chars[size-line-1] = true

		printLine(chars, line)

		chars[line] = false
		chars[size-line-1] = false

		lineBreak()
	}

	printBase()
	lineBreak()
}

// DisplayHourglass method displays the ASCII image of an hourglass to stdout with specific parameters.
func (p ParamStorage) DisplayHourglass() {
	displayHourglass(p())
}

// SetSize method returns a new ParamStorage with a new size.
func (p ParamStorage) SetSize(size int) ParamStorage {
	_, Char, CharColor, BackgroundColor := p()

	return func() (int, rune, bashcolor.Color, bashcolor.Color) {
		return size, Char, CharColor, BackgroundColor
	}
}

// SetChar method returns a new ParamStorage with a new char.
func (p ParamStorage) SetChar(char rune) ParamStorage {
	Size, _, CharColor, BackgroundColor := p()

	return func() (int, rune, bashcolor.Color, bashcolor.Color) {
		return Size, char, CharColor, BackgroundColor
	}
}

// SetCharColor method returns a new ParamStorage with a new char color.
func (p ParamStorage) SetCharColor(charColor bashcolor.Color) ParamStorage {
	Size, Char, _, BackgroundColor := p()

	return func() (int, rune, bashcolor.Color, bashcolor.Color) {
		return Size, Char, charColor, BackgroundColor
	}
}

// SetBackgroundColor method returns a new ParamStorage with a new background color.
func (p ParamStorage) SetBackgroundColor(backgroundColor bashcolor.Color) ParamStorage {
	Size, Char, CharColor, _ := p()

	return func() (int, rune, bashcolor.Color, bashcolor.Color) {
		return Size, Char, CharColor, backgroundColor
	}
}

// GetParamStorage returns a ParamStorage with default parameter values.
func GetParamStorage() ParamStorage {
	return func() (int, rune, bashcolor.Color, bashcolor.Color) {
		return 15, 'X', bashcolor.Blue, bashcolor.Black
	}
}
