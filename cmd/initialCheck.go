package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"siguma0013/reskk-dictionary/internal/dictionary"
	"slices"

	"github.com/spf13/cobra"
)

// initialCheckCmd represents the initialCheck command
var initialCheckCmd = &cobra.Command{
	Use:   "initialCheck",
	Short: "",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		results, walkErr := WalkJSONL(filePath, func(path string, r io.Reader) ([]error, error) {
			allowInitial, ok := dictionary.AllowInitials[filepath.Base(path)]

			if !ok {
				return nil, fmt.Errorf("辞書の10行分割で許可されていないファイル名です")
			}

			return checkInitials(r, allowInitial), nil
		})

		if walkErr != nil {
			return walkErr
		}

		if printResults(results, os.Stderr, "invalid lines") {
			return fmt.Errorf("invalid initial found")
		}

		fmt.Println("All JSONL files are valid")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initialCheckCmd)
}

// checkInitials 辞書ファイルの頭文字チェック関数
func checkInitials(reader io.Reader, allowInitial []string) []error {
	scanner := bufio.NewScanner(reader)
	lineCount := 0

	var errors []error

	// 1行づつ繰り返し処理
	for scanner.Scan() {
		lineCount++

		var record dictionary.Entry

		// パース
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			errors = append(errors, fmt.Errorf("parse error: %d", lineCount))
			continue
		}

		// 頭文字取得
		key := record.Key
		initial := string([]rune(key)[0])

		if !slices.Contains(allowInitial, initial) {
			errors = append(errors, fmt.Errorf("initial error: %d", lineCount))
			continue
		}
	}

	return errors
}
