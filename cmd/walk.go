package cmd

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// fileResult は1つのファイルの処理結果を表します。
// Path: ファイルパス
// Errors: ファイルの検査が返した行単位などのエラー一覧
// OpenErr: ファイルオープンや処理中に発生した致命的なエラー
type fileResult struct {
	Path          string
	Errors        []error
	CriticalError error
}

// WalkJSONL は root 以下の .jsonl ファイルを列挙し、
// 各ファイルに対して process を呼び出して結果を集約します。
// process は (path, io.Reader) を受け取り、行単位のエラー一覧と処理そのもののエラーを返す仕様です。
func WalkJSONL(
	root string,
	process func(path string, r io.Reader) ([]error, error),
) ([]fileResult, error) {
	var results []fileResult

	walkErr := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		// WalkDir からのエラーはそのまま返す
		if err != nil {
			return err
		}

		// ディレクトリはスキップ
		if d.IsDir() {
			return nil
		}

		// 拡張子jsonl以外が入力された場合エラー
		if !strings.HasSuffix(d.Name(), "jsonl") {
			results = append(results, fileResult{Path: path, CriticalError: fmt.Errorf("対応していない拡張子です")})
			return nil
		}

		// ファイルオープン
		file, err := os.Open(path)
		if err != nil {
			results = append(results, fileResult{Path: path, CriticalError: err})
			return nil
		}

		// 関数終了時にファイルクローズ呼出を強制
		defer file.Close()

		// 各コマンドの処理を実行
		validationErrors, processError := process(path, file)

		// 処理結果を作成、エラーの有無に関わらず作成
		results = append(results, fileResult{Path: path, Errors: validationErrors, CriticalError: processError})

		return nil
	})

	return results, walkErr
}
