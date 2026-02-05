package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

		// 入力されたパスに応じて再帰処理
		walkError := filepath.WalkDir(filePath, func(path string, meta os.DirEntry, err error) error {
			// WalkDir からのエラーはそのまま返す
			if err != nil {
				return err
			}

			// ディレクトリはスキップ
			if meta.IsDir() {
				return nil
			}

			// 拡張子jsonl以外はスキップ
			if !strings.HasSuffix(meta.Name(), "jsonl") {
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

			// フォーマットチェック実行
			validateErrors := validateJSONL(file)

			// エラーがなければスキップ
			if len(validateErrors) == 0 {
				return nil
			}

			// フラグ建て
			flagError = true

			// エラーファイル名、エラー数の出力
			fmt.Fprintf(os.Stderr, "%s: %d invalid lines:\n", path, len(validateErrors))

			// エラーログの出力
			for _, error := range validateErrors {
				fmt.Fprintf(os.Stderr, "  %v\n", error)
			}

			return nil
		})

		if walkError != nil {
			return walkError
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
