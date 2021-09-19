package main

import (
	"fmt"
	"tfs-go-hw/lection01/bashcolor"
	"tfs-go-hw/lection01/hourglass"
)

func main() {
	// With default parameters
	_, _, _, _, hgDefault := hourglass.GetHourglass()
	hgDefault()
	fmt.Println()

	// All parameters differ from the default
	setSize1, setChar1, setCharColor1, setBackgroundColor1, hg1 := hourglass.GetHourglass()
	setCharColor1(bashcolor.Red)
	setBackgroundColor1(bashcolor.LightGray)
	setChar1('F')
	setSize1(9)
	hg1()
	fmt.Println()

	// Only 2 parameters differ
	setSize2, _, setCharColor2, _, hg2 := hourglass.GetHourglass()
	setCharColor2(bashcolor.Green)
	setSize2(3)
	hg2()
}
