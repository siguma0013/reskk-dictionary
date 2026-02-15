package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"siguma0013/reskk-dictionary/internal/dictionary"
	"siguma0013/reskk-dictionary/internal/utility"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// sortCmd represents the sort command
var (
	isSortCi  bool
	isSortFix bool
)

var sortCmd = &cobra.Command{
	Use:          "sort",
	Short:        "辞書ファイルのソート確認&修正をするコマンド",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		var results []utility.FileResult

		if isSortFix {
			results = utility.WalkJsonl(filePath, nil, sortJsonl)
		} else {
			results = utility.WalkJsonl(filePath, sortFilter, func(path string, file io.Reader) []error {
				return checkSorted(file)
			})
		}

		if utility.PrintResults(results) {
			return fmt.Errorf("out-of-order keys found")
		}

		fmt.Println("All JSONL files are sorted")

		return nil
	},
}

func init() {
	sortCmd.Flags().BoolVar(&isSortCi, "ci", false, "use ci")
	sortCmd.Flags().BoolVar(&isSortFix, "fix", false, "Fix files by sorting keys in place")
	rootCmd.AddCommand(sortCmd)
}

func sortFilter(path string, _ string) bool {
	if !isSortCi {
		return true
	}

	fmt.Printf("is number.jsonl: %s\n", filepath.Base(path))

	return filepath.Base(path) != "number.jsonl"
}

func sortJsonl(path string, reader io.Reader) []error {
	// ソート済みデータの作成
	sorted, err := sortData(reader, dictionary.SortOrder())

	if err != nil {
		return []error{err}
	}

	// 出力ディレクトリの特定
	outputDir := filepath.Dir(path)

	// 一時ファイル作成
	tmp, err := os.CreateTemp(outputDir, ".tmp-*")
	if err != nil {
		return []error{err}
	}

	// 関数終了時に一時ファイルを削除
	defer os.Remove(tmp.Name())

	for _, e := range sorted {
		b, _ := json.Marshal(e)
		tmpLine := strings.ReplaceAll(string(b), ":\"", ": \"")
		tmpLine2 := strings.ReplaceAll(tmpLine, ":[", ": [")
		okLine := strings.ReplaceAll(tmpLine2, ",\"", ", \"")
		fmt.Fprintln(tmp, okLine)
	}

	// ファイル置換のために書き込みが完全終了してから処理移行
	tmp.Sync()
	tmp.Close()

	if err := os.Rename(tmp.Name(), path); err != nil {
		return []error{err}
	}

	return nil
}

// checkSorted checks that each successive 'key' is in non-decreasing order
func checkSorted(reader io.Reader) []error {
	scanner := bufio.NewScanner(reader)
	lineCount := 0

	var errors []error
	var prevKey string

	var orderMap = dictionary.SortOrder()

	// 1行づつ繰り返し処理
	for scanner.Scan() {
		lineCount++

		var record dictionary.Entry

		// パース
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			errors = append(errors, fmt.Errorf("parse error: %d", lineCount))
			continue
		}

		if prevKey != "" && compareKeys(prevKey, record.Key, orderMap) > 0 {
			errors = append(errors, fmt.Errorf("line %d: key %q is out of order after %q", lineCount, record.Key, prevKey))
		}

		prevKey = record.Key
	}

	if scannerError := scanner.Err(); scannerError != nil {
		errors = append(errors, fmt.Errorf("scanner error: %w", scannerError))
	}

	return errors
}

// compareKeys returns -1 if prevKey < currentKey, 0 if equal, 1 if prevKey > currentKey according to kana order map
func compareKeys(prevKey string, currentKey string, orderMap map[rune]int) int {
	prevRunes := []rune(prevKey)
	currentRunes := []rune(currentKey)

	for i := 0; i < len(prevRunes) && i < len(currentRunes); i++ {
		prevRune, prevOk := orderMap[prevRunes[i]]
		currentRune, currentOk := orderMap[currentRunes[i]]

		if !prevOk || !currentOk {
			if prevRunes[i] == currentRunes[i] {
				continue
			}
			if prevRunes[i] < currentRunes[i] {
				return -1
			}
			return 1
		}

		if prevRune < currentRune {
			return -1
		}

		if prevRune > currentRune {
			return 1
		}
	}

	if len(prevRunes) == len(currentRunes) {
		return 0
	}

	if len(prevRunes) < len(currentRunes) {
		return -1
	}

	return 1
}

// sortData ソート済みデータ作成関数
func sortData(reader io.Reader, orderMap map[rune]int) ([]dictionary.Entry, error) {
	scanner := json.NewDecoder(reader)

	var records []dictionary.Entry

	// ファイルパース
	for scanner.More() {
		var record dictionary.Entry

		if err := scanner.Decode(&record); err != nil {
			return nil, fmt.Errorf("parse error")
		}

		records = append(records, record)
	}

	// 1行だけなら修正不要のため、離脱
	if len(records) <= 1 {
		return records, nil
	}

	sort.SliceStable(records, func(i, j int) bool {
		return compareKeys(records[i].Key, records[j].Key, orderMap) < 0
	})

	return records, nil
}
