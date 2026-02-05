package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
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

		walkError := filepath.WalkDir(filePath, func(path string, d fs.DirEntry, err error) error {
			// WalkDir からのエラーはそのまま返す
			if err != nil {
				return err
			}

			// ディレクトリはスキップ
			if d.IsDir() {
				return nil
			}

			// 拡張子jsonl以外はスキップ
			if !strings.HasSuffix(d.Name(), "jsonl") {
				return nil
			}

			// 相対パスの取得
			rel, err := filepath.Rel(filePath, path)
			if err != nil {
				return err
			}

			// 階層が1つ下のものだけ対象
			parts := strings.Split(rel, string(os.PathSeparator))
			if len(parts) != 2 {
				return nil
			}

			allowInitial, ok := allowInitials[d.Name()]
			if !ok {
				fmt.Fprintf(os.Stderr, "contains disallow filename: %v\n", path)
				flagError = true
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

			// 頭文字チェック
			initialsErrors := checkInitials(file, allowInitial)

			// エラーがなければスキップ
			if len(initialsErrors) == 0 {
				return nil
			}

			// フラグ建て
			flagError = true

			// エラーファイル名、エラー数の出力
			fmt.Fprintf(os.Stderr, "%s: %d invalid lines:\n", path, len(initialsErrors))

			// エラーログの出力
			for _, error := range initialsErrors {
				fmt.Fprintf(os.Stderr, "  %v\n", error)
			}

			return nil
		})

		if walkError != nil {
			return walkError
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
