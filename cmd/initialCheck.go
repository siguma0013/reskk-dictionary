package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"siguma0013/reskk-dictionary/internal/entry"
	"slices"
	"strings"

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
		flagError := false

		results, walkErr := WalkJSONL(filePath, func(path string, r io.Reader) ([]error, error) {
			// 相対パスの取得
			rel, err := filepath.Rel(filePath, path)
			if err != nil {
				return nil, err
			}

			// 階層が1つ下のものだけ対象
			parts := strings.Split(rel, string(os.PathSeparator))
			if len(parts) != 2 {
				return nil, nil // 対象外
			}

			allowInitial, ok := allowInitials[filepath.Base(path)]
			if !ok {
				return []error{fmt.Errorf("contains disallow filename: %v", path)}, nil
			}

			return checkInitials(r, allowInitial), nil
		})

		if walkErr != nil {
			return walkErr
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
			return fmt.Errorf("invalid initial found")
		}

		fmt.Println("All JSONL files are valid")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initialCheckCmd)
}

var allowInitials = map[string][]string{
	"01-a.jsonl":  {"あ", "い", "う", "え", "お"},
	"02-ka.jsonl": {"か", "が", "き", "ぎ", "く", "ぐ", "け", "げ", "こ", "ご"},
	"03-sa.jsonl": {"さ", "ざ", "し", "じ", "す", "ず", "せ", "ぜ", "そ", "ぞ"},
	"04-ta.jsonl": {"た", "だ", "ち", "ぢ", "つ", "づ", "て", "で", "と", "ど"},
	"05-na.jsonl": {"な", "に", "ぬ", "ね", "の"},
	"06-ha.jsonl": {"は", "ば", "ぱ", "ひ", "び", "ぴ", "ふ", "ぶ", "ぷ", "へ", "べ", "ぺ", "ほ", "ぼ", "ぽ"},
	"07-ma.jsonl": {"ま", "み", "む", "め", "も"},
	"08-ya.jsonl": {"や", "ゆ", "よ"},
	"09-ra.jsonl": {"ら", "り", "る", "れ", "ろ"},
	"10-wa.jsonl": {"わ", "を", "ん"},
}

// checkInitials 辞書ファイルの頭文字チェック関数
func checkInitials(reader io.Reader, allowInitial []string) []error {
	scanner := bufio.NewScanner(reader)
	lineCount := 0

	var errors []error

	// 1行づつ繰り返し処理
	for scanner.Scan() {
		lineCount++

		var record entry.Entry

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
