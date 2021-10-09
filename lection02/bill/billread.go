package bill

import (
	"encoding/json"
	"flag"
	"os"
)

type InvalidInputError struct {
	msg string
}

func (e InvalidInputError) Error() string {
	return e.msg
}

// GetInput searches for a file with input data.
func GetInput() (*os.File, error) {
	type InputGetter func() (*os.File, error)

	flags := func() (*os.File, error) {
		fileArg := flag.String("file", "",
			"File of statements of financial operations of companies (.json format)")
		flag.Parse()
		fileName := FormatQuotes(*fileArg)

		return os.Open(fileName)
	}

	env := func() (*os.File, error) {
		fileName, ok := os.LookupEnv("FILE")

		if !ok {
			return nil, InvalidInputError{msg: "Environment variable was not passed"}
		}

		fileName = FormatQuotes(fileName)

		return os.Open(fileName)
	}

	stdin := func() (*os.File, error) {
		return os.Stdin, nil
	}

	for _, inputGetter := range [...]InputGetter{flags, env, stdin} {
		input, err := inputGetter()

		if err == nil {
			return input, nil
		}
	}

	return nil, InvalidInputError{msg: "Input stream was not passed"}
}

// ReadBills decodes the file with bills.
func ReadBills(input *os.File) ([]Bill, error) {
	var bills []Bill
	err := json.NewDecoder(input).Decode(&bills)

	if err != nil {
		return nil, err
	}

	return bills, nil
}
