package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"siguma0013/reskk-dictionary/internal/entry"
)

// sortCmd represents the sort command
var fixFlag bool

var orderMap = buildOrderMap()

var sortCmd = &cobra.Command{
	Use:          "sort",
	Short:        "Check that JSONL 'key's are sorted",
	Long:         "Check JSONL files recursively and report out-of-order keys.",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		results, walkError := WalkJSONL(filePath, func(path string, r io.Reader) ([]error, error) {
			if fixFlag {
				sorted, err := sortData(r, orderMap)
				if err != nil {
					return nil, err
				}
				if sorted == nil {
					return nil, nil
				}

				fmt.Fprintf(os.Stdout, "sorted %v\n", sorted)

				tmpDir := filepath.Dir(path)
				tmp, err := os.CreateTemp(tmpDir, ".tmp-*")
				if err != nil {
					return nil, err
				}
				defer os.Remove(tmp.Name())

				for _, e := range sorted {
					b, _ := json.Marshal(e)
					tmpLine := strings.ReplaceAll(string(b), ":\"", ": \"")
					tmpLine2 := strings.ReplaceAll(tmpLine, ":[", ": [")
					okLine := strings.ReplaceAll(tmpLine2, ",\"", ", \"")
					fmt.Fprintln(tmp, okLine)
				}

				tmp.Sync()
				tmp.Close()

				if err := os.Rename(tmp.Name(), path); err != nil {
					return nil, err
				}

				return nil, nil
			}

			return validateSortedJSONL(r), nil
		})

		if walkError != nil {
			return walkError
		}

		if printResults(results, os.Stderr, "ordering errors") {
			return fmt.Errorf("out-of-order keys found")
		}

		fmt.Println("All JSONL files are sorted")

		return nil
	},
}

var sortOrder = []string{
	"あ", "ぁ", "い", "ぃ", "う", "ぅ", "え", "ぇ", "お", "ぉ",
	"か", "が", "き", "ぎ", "く", "ぐ", "け", "げ", "こ", "ご",
	"さ", "ざ", "し", "じ", "す", "ず", "せ", "ぜ", "そ", "ぞ",
	"た", "だ", "ち", "ぢ", "つ", "っ", "づ", "て", "で", "と", "ど",
	"な", "に", "ぬ", "ね", "の",
	"は", "ば", "ぱ", "ひ", "び", "ぴ", "ふ", "ぶ", "ぷ", "へ", "べ", "ぺ", "ほ", "ぼ", "ぽ",
	"ま", "み", "む", "め", "も",
	"や", "ゃ", "ゆ", "ゅ", "よ", "ょ",
	"ら", "り", "る", "れ", "ろ",
	"わ", "を", "ん",
}

func init() {
	sortCmd.Flags().BoolVar(&fixFlag, "fix", false, "Fix files by sorting keys in place")
	rootCmd.AddCommand(sortCmd)
}

// validateSortedJSONL checks that each successive 'key' is in non-decreasing order
func validateSortedJSONL(reader io.Reader) []error {
	scanner := bufio.NewScanner(reader)
	lineCount := 0

	var errors []error
	var prevKey string

	// 1行づつ繰り返し処理
	for scanner.Scan() {
		lineCount++

		var record entry.Entry

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

func buildOrderMap() map[rune]int {
	orderMap := make(map[rune]int)

	for index, s := range sortOrder {
		r := []rune(s)
		if len(r) > 0 {
			orderMap[r[0]] = index
		}
	}

	return orderMap
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
func sortData(reader io.Reader, orderMap map[rune]int) ([]entry.Entry, error) {
	scanner := json.NewDecoder(reader)

	var records []entry.Entry

	// ファイルパース
	for scanner.More() {
		var record entry.Entry

		if err := scanner.Decode(&record); err != nil {
			return nil, fmt.Errorf("parse error")
		}

		records = append(records, record)
	}

	// 1行だけなら修正不要のため、離脱
	if len(records) <= 1 {
		return nil, nil
	}

	sort.SliceStable(records, func(i, j int) bool {
		return compareKeys(records[i].Key, records[j].Key, orderMap) < 0
	})

	return records, nil
}
