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

// SaveYaml : データモデルをYAMLファイルへ書き出す
func SaveYaml[T any](path string, data T) error {
	var buf bytes.Buffer

	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode YAML: %w", err)
	}

	if err := encoder.Close(); err != nil {
		return fmt.Errorf("failed to encode YAML: %w", err)
	}

	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
