package bill

import (
	"fmt"
	"os"
	"strings"
)

// FormatQuotes is needed to format the string containing the file name.
func FormatQuotes(str string) string {
	return strings.TrimSpace(
		strings.Trim(
			strings.TrimSpace(str),
			"\"",
		),
	)
}

func DeferClose(file *os.File) {
	err := file.Close()

	if err != nil {
		fmt.Println(err)
	}
}
