package validator

import (
	"fmt"

	"siguma0013/reskk-dictionary/internal/model"
)

// バリデーションの基本処理
func validateEntries[T model.Entry](entries []T, filePath string, yomiValidator func(T) []error) []error {
	// 全体のバリデーション結果
	var errs []error

	// 索引を取得
	indexKeys, err := getIndexKeys(filePath)
	if err != nil {
		errs = append(errs, fmt.Errorf("  - %w", err))
	}

	// 漢字の文字列長を取得
	length, err := getKanjiLength(filePath)
	if err != nil {
		errs = append(errs, fmt.Errorf("  - %w", err))
	}

	// 配置に問題があるかチェック
	if len(errs) > 0 {
		// 問題がある場合は内容をチェックせずにreturn
		return errs
	}

	// 漢字の重複確認用ハッシュマップ
	seenKanji := make(map[string]struct{})

	// 辞書データの繰り返し処理
	for _, entry := range entries {
		var entryErrs []error

		// 漢字のバリデーション
		if err := validateKanji(entry.GetKanji(), length, seenKanji); err != nil {
			entryErrs = append(entryErrs, fmt.Errorf("  - %s: %w", entry.GetKanji(), err))
		}

		// 読みのバリデーション
		for _, err := range yomiValidator(entry) {
			entryErrs = append(entryErrs, fmt.Errorf("  - %s: %w", entry.GetKanji(), err))
		}

		if len(entryErrs) > 0 {
			errs = append(errs, entryErrs...)
			continue
		}

		// 代表読み・索引のバリデーション
		if err := validateIndex(entry, indexKeys); err != nil {
			entryErrs = append(entryErrs, fmt.Errorf("  - %s: %w", entry.GetKanji(), err))
		}

		errs = append(errs, entryErrs...)
	}

	return errs
}
