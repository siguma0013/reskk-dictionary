package loader

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

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

// SaveYaml : データモデルを一時ファイルに書き出してから置き換える
func SaveYaml[T any](path string, data T) (result error) {
	// YAMLをメモリ上で完成させ、書き込み途中の内容が元ファイルへ反映されないようにする
	var buf bytes.Buffer

	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode YAML: %w", err)
	}

	if err := encoder.Close(); err != nil {
		return fmt.Errorf("failed to encode YAML: %w", err)
	}

	// 置換後も既存ファイルの権限を維持する
	mode := os.FileMode(0o644)
	if info, err := os.Stat(path); err == nil {
		mode = info.Mode().Perm()
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	// Renameが同一ファイルシステム内で行われるよう、保存先と同じディレクトリに一時ファイルを作る
	tempFile, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	tempPath := tempFile.Name()
	// 置換前に失敗した場合も、一時ファイルを残さない
	defer func() {
		if err := os.Remove(tempPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			result = errors.Join(result, fmt.Errorf("failed to remove temporary file: %w", err))
		}
	}()
	// 書き込み中にエラーが起きた場合は、開いた一時ファイルを閉じる
	defer func() {
		if tempFile != nil {
			if err := tempFile.Close(); err != nil {
				result = errors.Join(result, fmt.Errorf("failed to close temporary file: %w", err))
			}
		}
	}()

	// 一時ファイルにも元ファイルと同じ権限を設定する
	if err := tempFile.Chmod(mode); err != nil {
		return fmt.Errorf("failed to set temporary file permissions: %w", err)
	}
	// 完成済みのYAMLを一時ファイルへ書き込む
	if _, err := tempFile.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}
	// 置換前に内容をストレージへ同期する
	if err := tempFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync temporary file: %w", err)
	}
	// 書き込みを完了してから置換する
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %w", err)
	}
	tempFile = nil

	// 完成した一時ファイルで元ファイルを置き換える
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("failed to replace file: %w", err)
	}
	return nil
}
