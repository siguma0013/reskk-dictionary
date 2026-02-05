package cmd

import (
	"fmt"
	"io"
)

// printResults は WalkJSONL の結果を標準エラーへ出力し、何らかのエラーがあったかを返します。
// listLabel は各ファイルのエラー一覧表示に用いる文言（例: "invalid lines" や "ordering errors"）です。
func printResults(results []fileResult, errOut io.Writer, listLabel string) bool {
	flagError := false

	for _, res := range results {
		if res.CriticalError != nil && len(res.Errors) == 0 {
			fmt.Fprintf(errOut, "[NG] %s: %v\n", res.Path, res.CriticalError)
			flagError = true
			continue
		}

		if len(res.Errors) == 0 {
			fmt.Fprintf(errOut, "[OK] %s\n", res.Path)
			continue
		}

		flagError = true
		fmt.Fprintf(errOut, "[NG] %s: %d %s:\n", res.Path, len(res.Errors), listLabel)
		for _, e := range res.Errors {
			fmt.Fprintf(errOut, "  %v\n", e)
		}
	}

	return flagError
}
