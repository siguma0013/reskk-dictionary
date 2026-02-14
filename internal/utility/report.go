package utility

import (
	"fmt"
	"os"
)

func PrintResults(results []FileResult) bool {
	errorFlag := false

	for _, result := range results {
		if len(result.Errors) == 0 {
			fmt.Fprintf(os.Stdout, "[OK]%s\n", result.Path)
			continue
		}

		errorFlag = true

		fmt.Fprintf(os.Stdout, "[NG]%s\n", result.Path)

		for _, err := range result.Errors  {
			fmt.Fprintf(os.Stdout, " - %v\n", err)
		}
	}

	return errorFlag
}
