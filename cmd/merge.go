package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"siguma0013/reskk-dictionary/internal/dictionary"
	"slices"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// オプション
var (
	mergeOrderPath  string
	mergeOutputPath string
)

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge JSONL files according",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		orders, err := makeMergeOrder(mergeOrderPath)

		if err != nil {
			return fmt.Errorf("nothing order %s: %w", mergeOrderPath, err)
		}

		mergeData, err := makeMergeData(orders)
		if err != nil {
			return fmt.Errorf("missing merge data")
		}

		// これより出力処理
		outFile, err := os.Create(mergeOutputPath)
		if err != nil {
			return fmt.Errorf("failed to create %s: %w", mergeOutputPath, err)
		}

		defer outFile.Close()

		for key, value := range mergeData {
			jsonObject := dictionary.Entry{
				Key:   key,
				Value: value,
			}

			jsonString, err := json.Marshal(jsonObject)
			if err != nil {
				return fmt.Errorf("missing make json string")
			}

			outFile.Write(jsonString)
			outFile.Write([]byte("\n"))
		}

		return nil
	},
}

func init() {
	mergeCmd.Flags().StringVar(&mergeOrderPath, "input", "merge_order.yml", "input order file")
	mergeCmd.Flags().StringVar(&mergeOutputPath, "output", "merged.jsonl", "output file")
	rootCmd.AddCommand(mergeCmd)
}

// makeMergeOrder はyamlからファイルリストを読み込む
func makeMergeOrder(path string) ([]string, error) {
	orderFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read order file")
	}

	var order struct {
		Files []string `yaml:"files"`
	}

	// YAMLファイルのパース
	if err := yaml.Unmarshal(orderFile, &order); err != nil {
		return nil, fmt.Errorf("failed to read order file")
	}

	// orderList は実際に使用するリスト、スライスのため順序がある
	var orderList []string

	for _, pattern := range order.Files {
		// とりあえず全てGlobで解決を行う
		matches, err := filepath.Glob(pattern)

		// 解決できなければスキップ
		if err != nil {
			continue
		}

		for _, matche := range matches {
			if !fileExists(matche) {
				continue
			}

			// スライスに含まれていなければ追加
			if !slices.Contains(orderList, matche) {
				orderList = append(orderList, matche)
			}
		}
	}

	// リストがゼロであればエラーとする
	if len(orderList) == 0 {
		return nil, fmt.Errorf("no order")
	}

	return orderList, nil
}

// fileExists はファイルの有無を確認する
// ファイルがあるときtrueを返す
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func makeMergeData(orders []string) (map[string][]string, error) {
	outMap := make(map[string][]string)

	for _, path := range orders {
		// jsonlファイルオープン
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		// 関数終了時にファイルクローズ呼出を強制
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			var record dictionary.Entry

			if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
				return nil, fmt.Errorf("parse error: %s", path)
			}

			// 初回データはそのまま投入
			if _, ok := outMap[record.Key]; !ok {
				outMap[record.Key] = record.Value
				continue
			}

			// ここから先は重複データ
			outMap[record.Key] = mergeSlice(outMap[record.Key], record.Value)
		}
	}

	return outMap, nil
}

func mergeSlice(source []string, input []string) []string {
	for _, value := range input {
		if slices.Contains(source, value) {
			continue
		}

		source = append(source, value)
	}

	return source
}
