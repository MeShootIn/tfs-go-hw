package main

import (
	"fmt"
	"tfs-go-hw/lection01/bashcolor"
	"tfs-go-hw/lection01/hourglass"
)

// An example how to use package hourglass
func main() {
	for _, paramStorage := range [...]hourglass.ParamStorage{
		// With default parameters
		hourglass.GetParamStorage(),
		// All parameters differ from the default
		hourglass.GetParamStorage().
			SetCharColor(bashcolor.Red).
			SetBackgroundColor(bashcolor.LightGray).
			SetChar('F').
			SetSize(7),
		// Only 2 parameters differ
		hourglass.GetParamStorage().
			SetCharColor(bashcolor.Green).
			SetSize(3),
	} {
		paramStorage.DisplayHourglass()
		fmt.Println()
	}
}
