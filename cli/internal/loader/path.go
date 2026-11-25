package loader

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileExists : パスが存在するファイルパスであることを検証
func FileExists(path string) error {
	info, err := os.Stat(path)

	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("%q is a directory", path)
	}

	return nil
}

// DictionaryType : 辞書の種類
type DictionaryType string

const (
	CharacterDictionary DictionaryType = "単漢字"
	WordDictionary      DictionaryType = "熟語"
)

// GetDictionaryType : ファイルの親ディレクトリ名から辞書の種類を判定
func GetDictionaryType(path string) (DictionaryType, error) {
	subcategoryDir := filepath.Dir(filepath.Clean(path))
	categoryDir := filepath.Dir(subcategoryDir)

	categoryName := filepath.Base(subcategoryDir)
	if categoryName != string(CharacterDictionary) && categoryName != string(WordDictionary) {
		categoryName = filepath.Base(categoryDir)
	}

	switch categoryName {
	case string(CharacterDictionary):
		return CharacterDictionary, nil
	case string(WordDictionary):
		return WordDictionary, nil
	default:
		return "", fmt.Errorf("管理対象外の辞書ディレクトリです: %q", path)
	}
}
