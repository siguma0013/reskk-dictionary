package writer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"siguma0013/reskk-dictionary/internal/model"
)

// WriteJSON : 辞書データモデルをJSONファイルへ出力する
// 同名ファイルが存在する場合は上書きする
func WriteJSON(dict model.Dictionary, path string) error {
	// nilマップが null として出力されることを防ぐ
	if dict == nil {
		dict = model.Dictionary{}
	}

	// インデントなし、末尾改行付き(キーはencoding/jsonにより昇順)
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(dict); err != nil {
		return fmt.Errorf("failed to encode json: %w", err)
	}

	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
