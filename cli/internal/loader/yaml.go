package loader

import (
	"bytes"
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

// LoadYaml : YAMLファイルを読み込んでデータモデルを返す
func LoadYaml[T any](path string) (T, error) {
	// 最終成果物
	var result T

	// ファイルオープン
	data, err := os.ReadFile(path)

	if err != nil {
		// ファイルオープン失敗
		return result, fmt.Errorf("failed to read file: %w", err)
	}

	decoer := yaml.NewDecoder(bytes.NewReader(data))
	decoer.KnownFields(true)

	if err := decoer.Decode(&result); err != nil {
		// デコード失敗
		return result, fmt.Errorf("failed to decode YAML: %w", err)
	}

	return result, nil
}
