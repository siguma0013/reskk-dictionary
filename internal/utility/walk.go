package utility

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type FileResult struct {
	Path string
	Errors []error
}

// WalkJsonl はWalkDirのラッパー
//  - filter が true を返すファイルはスキップされる
//  - process にエラーチェックロジックを実装する
func WalkJsonl(
	root string,
	filter func(path string, root string) bool,
	process func(path string, file io.Reader) []error,
) []FileResult {
	var results []FileResult

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		// WalkDir からのエラーを処理
		if err != nil {
			results = append(results, FileResult{Path: path, Errors: []error{err}})
			return nil
		}

		// ディレクトリはスキップ
		if d.IsDir() {
			return nil
		}

		// 拡張子jsonl以外が入力された場合エラー
		if !strings.HasSuffix(d.Name(), "jsonl") {
			results = append(results, FileResult{
				Path: path,
				Errors: []error{fmt.Errorf("対応していない拡張子です")},
			})
			return nil
		}

		// filter でスキップ処理
		if filter != nil && !filter(path, root) {
			return nil
		}

		// ファイルオープン
		file, err := os.Open(path)
		if err != nil {
			results = append(results, FileResult{Path: path, Errors: []error{err}})
			return nil
		}

		// 関数終了時にファイルクローズ呼出を強制
		defer file.Close()

		// 各コマンドの処理を実行
		if err := process(path, file); err != nil {
			results = append(results, FileResult{Path: path, Errors: err})
			return nil
		}

		// エラーなしで結果を作成
		results = append(results, FileResult{Path: path})

		return  nil
	})

	return results
}

// FileDepth は path の深さを返す関数
func FileDepth(path string) int {
	return strings.Count(filepath.Clean(path), string(filepath.Separator))
}
