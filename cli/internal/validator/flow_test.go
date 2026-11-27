package validator

import (
	"errors"
	"strings"
	"testing"

	"siguma0013/reskk-dictionary/internal/model"
)

func TestValidateEntries(t *testing.T) {
	t.Run("正常なエントリは確認を通る", func(t *testing.T) {
		entries := []testEntry{
			{kanji: "市場", indexYomi: "かいば", indexKey: "か"},
		}

		errs := validateEntries(entries, "/dictionary/二字熟語/11-か.yaml", func(testEntry) []error {
			return nil
		})

		if len(errs) != 0 {
			t.Errorf("validateEntries() errors = %v, want no errors", errs)
		}
	})

	t.Run("ファイル名が管理対象外の場合はエントリ内容を確認せずエラーになる", func(t *testing.T) {
		entries := []testEntry{
			{kanji: "市場", indexYomi: "いちば", indexKey: "い"},
		}

		errs := validateEntries(entries, "/dictionary/管理外/99-unknown.yaml", func(testEntry) []error {
			t.Fatal("管理対象外のファイルでは読みバリデーションを呼び出してはいけない")
			return nil
		})

		if len(errs) != 2 {
			t.Fatalf("validateEntries() errors = %v, want 2 errors", errs)
		}
		if !strings.Contains(errs[0].Error(), "管理対象外のファイル名") {
			t.Errorf("validateEntries() first error = %v, want 管理対象外のファイル名", errs[0])
		}
		if !strings.Contains(errs[1].Error(), "管理対象外のディレクトリ名") {
			t.Errorf("validateEntries() second error = %v, want 管理対象外のディレクトリ名", errs[1])
		}
	})

	t.Run("漢字が重複しているエントリはエラーになる", func(t *testing.T) {
		entries := []testEntry{
			{kanji: "市場", indexYomi: "かいば", indexKey: "か"},
			{kanji: "市場", indexYomi: "がいじょう", indexKey: "が"},
		}

		errs := validateEntries(entries, "/dictionary/二字熟語/11-か.yaml", func(testEntry) []error {
			return nil
		})

		if len(errs) != 1 || !strings.Contains(errs[0].Error(), "漢字が重複しています") {
			t.Errorf("validateEntries() errors = %v, want duplicate Kanji error", errs)
		}
	})

	t.Run("読みバリデーションでエラーがある場合は索引確認を行わない", func(t *testing.T) {
		entries := []testEntry{
			{kanji: "市場", indexYomi: "", indexKey: ""},
		}

		errs := validateEntries(entries, "/dictionary/二字熟語/11-か.yaml", func(testEntry) []error {
			return []error{errors.New("文字種が不正です")}
		})

		if len(errs) != 1 || !strings.Contains(errs[0].Error(), "文字種が不正です") {
			t.Errorf("validateEntries() errors = %v, want reading error", errs)
		}
	})

	t.Run("代表読みの索引文字がファイルの索引プロフィールと一致しない場合はエラーになる", func(t *testing.T) {
		entries := []testEntry{
			{kanji: "市場", indexYomi: "しじょう", indexKey: "し"},
		}

		errs := validateEntries(entries, "/dictionary/二字熟語/11-か.yaml", func(testEntry) []error {
			return nil
		})

		if len(errs) != 1 || !strings.Contains(errs[0].Error(), "索引が正しくありません") {
			t.Errorf("validateEntries() errors = %v, want index error", errs)
		}
	})
}

func TestValidateCharacterReadings(t *testing.T) {
	t.Run("読みが空の単漢字エントリは文字種確認を通る", func(t *testing.T) {
		errs := validateCharacterReadings(model.CharacterEntry{})
		if len(errs) != 0 {
			t.Errorf("validateCharacterReadings() errors = %v, want no errors", errs)
		}
	})
}

func TestValidateWordReadings(t *testing.T) {
	t.Run("読みが空の熟語エントリは文字種確認を通る", func(t *testing.T) {
		errs := validateWordReadings(model.WordEntry{})
		if len(errs) != 0 {
			t.Errorf("validateWordReadings() errors = %v, want no errors", errs)
		}
	})
}
