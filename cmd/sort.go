package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"siguma0013/reskk-dictionary/internal/entry"
)

// sortCmd represents the sort command
var sortCmd = &cobra.Command{
	Use:          "sort",
	Short:        "Check that JSONL 'key's are sorted",
	Long:         "Check JSONL files recursively and report out-of-order keys.",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		flagError := false

		walkError := filepath.WalkDir(filePath, func(path string, d fs.DirEntry, err error) error {
			// WalkDir からのエラーはそのまま返す
			if err != nil {
				return err
			}

			// ディレクトリはリセット
			if d.IsDir() {
				return nil
			}

			// 拡張子jsonl以外はスキップ
			if !strings.HasSuffix(d.Name(), "jsonl") {
				return nil
			}

			// ファイルオープン
			file, fileError := os.Open(path)

			// ファイルオープン失敗はログ出力してスキップ
			if fileError != nil {
				fmt.Fprintf(os.Stderr, "failed to open %s: %v\n", path, fileError)
				flagError = true
				return nil
			}

			// 関数終了時にファイルクローズ
			defer file.Close()

			validateErrors := validateSortedJSONL(file)

			// エラーがなければスキップ
			if len(validateErrors) == 0 {
				return nil
			}

			// フラグ建て
			flagError = true

			// エラーファイル名、エラー数の出力
			fmt.Fprintf(os.Stderr, "%s: %d ordering errors:\n", path, len(validateErrors))

			// エラーログの出力
			for _, e := range validateErrors {
				fmt.Fprintf(os.Stderr, "  %v\n", e)
			}

			return nil
		})

		if walkError != nil {
			return walkError
		}

		if flagError {
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
	rootCmd.AddCommand(sortCmd)
}

// validateSortedJSONL checks that each successive 'key' is in non-decreasing order
func validateSortedJSONL(reader io.Reader) []error {
	scanner := bufio.NewScanner(reader)
	lineCount := 0

	var errors []error
	var prevKey string

	orderMap := buildOrderMap()

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
