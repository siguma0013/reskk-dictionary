package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"siguma0013/reskk-dictionary/internal/entry"
)

var formatCheckCmd = &cobra.Command{
	Use:          "format",
	Short:        "Check JSONL format correctness under a directory",
	Long:         `Check JSONL files recursively and report invalid JSON lines with file and line numbers.`,
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error { // 実行時に呼ばれる関数（エラーを返せる）
		filePath := args[0]
		flagError := false

		results, walkError := WalkJSONL(filePath, func(path string, r io.Reader) ([]error, error) {
			return validateJSONL(r), nil
		})

		if walkError != nil {
			return walkError
		}

		for _, res := range results {
			if res.CriticalError != nil && len(res.Errors) == 0 {
				fmt.Fprintf(os.Stderr, "failed to open %s: %v\n", res.Path, res.CriticalError)
				flagError = true
				continue
			}

			if len(res.Errors) == 0 {
				continue
			}

			flagError = true
			fmt.Fprintf(os.Stderr, "%s: %d invalid lines:\n", res.Path, len(res.Errors))
			for _, err := range res.Errors {
				fmt.Fprintf(os.Stderr, "  %v\n", err)
			}
		}

		if flagError {
			return fmt.Errorf("invalid JSONL format found")
		}

		fmt.Println("All JSONL files are valid")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(formatCheckCmd)
}

// validateJSONL 辞書ファイルのフォーマットチェック本体
func validateJSONL(reader io.Reader) []error {
	scanner := bufio.NewScanner(reader)
	lineCount := 0

	var errors []error

	// 1行づつ繰り返し処理
	for scanner.Scan() {
		lineCount++

		line := scanner.Text()

		// 空行がある時、エラー
		if line == "" {
			errors = append(errors, fmt.Errorf("line %d: empty line", lineCount))
			continue
		}

		// 前後にスペースがある時、エラー
		if strings.TrimSpace(line) != line {
			errors = append(errors, fmt.Errorf("line %d: include trailing space", lineCount))
			continue
		}

		decoder := json.NewDecoder(strings.NewReader(line))
		decoder.DisallowUnknownFields()

		var record entry.Entry

		if decodeError := decoder.Decode(&record); decodeError != nil {
			errors = append(errors, fmt.Errorf("line %d: schema error", lineCount))
			continue
		}

		// keyの有無
		if record.Key == "" {
			errors = append(errors, fmt.Errorf("line %d: empty key", lineCount))
			continue
		}

		// valueの有無
		if len(record.Value) == 0 {
			errors = append(errors, fmt.Errorf("line %d: empty value", lineCount))
			continue
		}

		for _, rule := range entry.FormatRules {
			if rule.Regexp.MatchString(line) {
				errors = append(errors, fmt.Errorf("line %d: %v", lineCount, rule.Message))
			}
		}
	}

	if scannerError := scanner.Err(); scannerError != nil { // Scanner 自身のエラー（IO エラー等）をチェック
		errors = append(errors, fmt.Errorf("scanner error: %w", scannerError))
	}

	return errors
}
