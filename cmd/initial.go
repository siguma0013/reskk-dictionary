package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"siguma0013/reskk-dictionary/internal/dictionary"
	"siguma0013/reskk-dictionary/internal/utility"
	"slices"

	"github.com/spf13/cobra"
)

var (
	isInitialCi bool
)

// initialCheckCmd represents the initialCheck command
var initialCheckCmd = &cobra.Command{
	Use:          "initial",
	Short:        "10行分割辞書ファイルが適切に分割されているかチェックするコマンド",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		results := utility.WalkJsonl(filePath, filterInitial, func(path string, file io.Reader) []error {
			allowInitial, ok := dictionary.AllowInitials[filepath.Base(path)]

			if !ok {
				return []error{fmt.Errorf("辞書の10行分割で許可されていないファイル名です")}
			}

			return checkInitial(file, allowInitial)
		})

		if utility.PrintResults(results) {
			return fmt.Errorf("invalid initial found")
		}

		fmt.Println("All JSONL files are valid")

		return nil
	},
}

func init() {
	initialCheckCmd.Flags().BoolVar(&isInitialCi, "ci", false, "use ci")

	rootCmd.AddCommand(initialCheckCmd)
}

func filterInitial(path string, root string) bool {
	pathDepth := utility.FileDepth(path)
	rootDepth := utility.FileDepth(root)

	if isInitialCi {
		return (pathDepth - rootDepth) == 2
	} else {
		return true
	}
}

// checkInitial 辞書ファイルの頭文字チェック関数
func checkInitial(reader io.Reader, allowInitial []string) []error {
	scanner := bufio.NewScanner(reader)
	lineCount := 0

	var results []error

	// 1行づつ繰り返し処理
	for scanner.Scan() {
		lineCount++

		var record dictionary.Entry

		// パース
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			results = append(results, fmt.Errorf("parse error: %d", lineCount))
			continue
		}

		// 頭文字取得
		key := record.Key
		initial := string([]rune(key)[0])

		if !slices.Contains(allowInitial, initial) {
			results = append(results, fmt.Errorf("initial error: %d", lineCount))
			continue
		}
	}

	return results
}
