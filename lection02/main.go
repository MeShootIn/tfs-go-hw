package main

import (
	"fmt"
	"lection02/bill"
)

const outputFileName = "out.json"

// An example how to use package bill
func main() {
	input, err := bill.GetInput()
	defer bill.DeferClose(input)

	if err != nil {
		fmt.Println(err)

		return
	}

	bills, err := bill.ReadBills(input)

	if err != nil {
		fmt.Println(err)

		return
	}

	err = bill.ProcessBills(bills, outputFileName)

	if err != nil {
		fmt.Println(err)
	}
}
