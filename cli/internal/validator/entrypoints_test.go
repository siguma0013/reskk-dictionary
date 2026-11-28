package validator

import (
	"strings"
	"testing"

	"siguma0013/reskk-dictionary/internal/model"
)

func TestValidateCharacters(t *testing.T) {
	t.Run("エントリがない単漢字辞書は確認を通る", func(t *testing.T) {
		errs := ValidateCharacters(model.Characters{}, "/dictionary/常用漢字/01-あ.yaml")
		if len(errs) != 0 {
			t.Errorf("ValidateCharacters() errors = %v, want no errors", errs)
		}
	})

	t.Run("管理対象外のパスの単漢字辞書はエラーになる", func(t *testing.T) {
		errs := ValidateCharacters(model.Characters{}, "/dictionary/管理外/99-unknown.yaml")
		if len(errs) != 2 {
			t.Fatalf("ValidateCharacters() errors = %v, want 2 errors", errs)
		}
		if !strings.Contains(errs[0].Error(), "管理対象外のファイル名") {
			t.Errorf("ValidateCharacters() first error = %v, want file name error", errs[0])
		}
		if !strings.Contains(errs[1].Error(), "管理対象外のディレクトリ名") {
			t.Errorf("ValidateCharacters() second error = %v, want directory name error", errs[1])
		}
	})
}

func TestValidateWords(t *testing.T) {
	t.Run("エントリがない熟語辞書は確認を通る", func(t *testing.T) {
		errs := ValidateWords(model.Words{}, "/dictionary/二字熟語/11-か.yaml")
		if len(errs) != 0 {
			t.Errorf("ValidateWords() errors = %v, want no errors", errs)
		}
	})

	t.Run("管理対象外のパスの熟語辞書はエラーになる", func(t *testing.T) {
		errs := ValidateWords(model.Words{}, "/dictionary/管理外/99-unknown.yaml")
		if len(errs) != 2 {
			t.Fatalf("ValidateWords() errors = %v, want 2 errors", errs)
		}
		if !strings.Contains(errs[0].Error(), "管理対象外のファイル名") {
			t.Errorf("ValidateWords() first error = %v, want file name error", errs[0])
		}
		if !strings.Contains(errs[1].Error(), "管理対象外のディレクトリ名") {
			t.Errorf("ValidateWords() second error = %v, want directory name error", errs[1])
		}
	})
}

func TestValidatedEntriesProduceNonNilDictionaryCandidates(t *testing.T) {
	t.Run("単漢字", func(t *testing.T) {
		entries := model.Characters{Entries: []model.CharacterEntry{
			{Kanji: "愛", OnYomi: []string{"あい"}},
		}}
		if errs := ValidateCharacters(entries, "/dictionary/常用漢字/01-あ.yaml"); len(errs) > 0 {
			t.Fatalf("ValidateCharacters() errors = %v", errs)
		}

		dict := model.Dictionary{}
		for _, entry := range entries.Entries {
			dict.AddEntry(entry)
		}
		for yomi, candidates := range dict {
			if candidates == nil {
				t.Errorf("nil candidates for reading %q", yomi)
			}
		}
		if got := dict["あい"]; len(got) != 1 || got[0] != "愛" {
			t.Errorf("dictionary candidates for あい = %v, want [愛]", got)
		}
	})

	t.Run("熟語", func(t *testing.T) {
		entries := model.Words{Entries: []model.WordEntry{
			{Kanji: "会話", Yomi: []string{"かいわ"}},
		}}
		if errs := ValidateWords(entries, "/dictionary/二字熟語/11-か.yaml"); len(errs) > 0 {
			t.Fatalf("ValidateWords() errors = %v", errs)
		}

		dict := model.Dictionary{}
		for _, entry := range entries.Entries {
			dict.AddEntry(entry)
		}
		for yomi, candidates := range dict {
			if candidates == nil {
				t.Errorf("nil candidates for reading %q", yomi)
			}
		}
		if got := dict["かいわ"]; len(got) != 1 || got[0] != "会話" {
			t.Errorf("dictionary candidates for かいわ = %v, want [会話]", got)
		}
	})
}

func TestValidateEmptyReadingsAreRejected(t *testing.T) {
	cases := []struct {
		name string
		errs func() []error
		want string
	}{
		{"単漢字: 読みなし", func() []error {
			return ValidateCharacters(model.Characters{Entries: []model.CharacterEntry{{Kanji: "愛"}}}, "/dictionary/常用漢字/01-あ.yaml")
		}, "読みは1つ以上必要です"},
		{"単漢字: 空文字の読み", func() []error {
			return ValidateCharacters(model.Characters{Entries: []model.CharacterEntry{{Kanji: "愛", OnYomi: []string{""}}}}, "/dictionary/常用漢字/01-あ.yaml")
		}, "文字種が不正です"},
		{"熟語: 読みなし", func() []error {
			return ValidateWords(model.Words{Entries: []model.WordEntry{{Kanji: "会話"}}}, "/dictionary/二字熟語/11-か.yaml")
		}, "読みは1つ以上必要です"},
		{"熟語: 空文字の読み", func() []error {
			return ValidateWords(model.Words{Entries: []model.WordEntry{{Kanji: "会話", Yomi: []string{""}}}}, "/dictionary/二字熟語/11-か.yaml")
		}, "文字種が不正です"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := tc.errs()
			if len(errs) != 1 || !strings.Contains(errs[0].Error(), tc.want) {
				t.Errorf("errors = %v, want one error containing %q", errs, tc.want)
			}
		})
	}
}
