package validator

import (
	"strings"
	"testing"

	"siguma0013/reskk-dictionary/internal/model"
)

// testEntry : 未実装のメソッドは埋め込みのEntry(nil)に委譲されるため、テストで使うメソッドだけ実装する
type testEntry struct {
	model.Entry
	kanji     string
	indexYomi string
	indexKey  string
}

func (e testEntry) GetKanji() string {
	return e.kanji
}

func (e testEntry) GetIndexYomi() string {
	return e.indexYomi
}

func (e testEntry) GetIndexKey() string {
	return e.indexKey
}

func TestGetIndexKeys(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     []string
		wantErr  bool
	}{
		{
			name:     "あ行のファイル名から「あ」の索引文字を取得できる",
			filePath: "/dictionary/あ/01-あ.yaml",
			want:     []string{"あ"},
		},
		{
			name:     "濁音を含むファイル名から複数の索引文字を取得できる",
			filePath: "/dictionary/か/11-か.yaml",
			want:     []string{"か", "が"},
		},
		{
			name:     "管理対象外のファイル名ではエラーになる",
			filePath: "/dictionary/未知/99-unknown.yaml",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getIndexKeys(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Fatalf("getIndexKeys() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if strings.Join(got, ",") != strings.Join(tt.want, ",") {
				t.Errorf("getIndexKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetKanjiLength(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     int
		wantErr  bool
	}{
		{
			name:     "常用漢字ディレクトリでは文字列長1を取得できる",
			filePath: "/dictionary/常用漢字/01-あ.yaml",
			want:     1,
		},
		{
			name:     "二字熟語ディレクトリでは文字列長2を取得できる",
			filePath: "/dictionary/二字熟語/11-か.yaml",
			want:     2,
		},
		{
			name:     "三字熟語ディレクトリでは文字列長3を取得できる",
			filePath: "/dictionary/三字熟語/21-さ.yaml",
			want:     3,
		},
		{
			name:     "四字熟語ディレクトリでは文字列長4を取得できる",
			filePath: "/dictionary/四字熟語/31-た.yaml",
			want:     4,
		},
		{
			name:     "管理対象外のディレクトリではエラーになる",
			filePath: "/dictionary/五字熟語/41-な.yaml",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getKanjiLength(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Fatalf("getKanjiLength() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("getKanjiLength() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestValidateKanji(t *testing.T) {
	t.Run("指定した文字列長の漢字は重複していなければ確認を通る", func(t *testing.T) {
		seen := make(map[string]struct{})
		if err := validateKanji("市場", 2, seen); err != nil {
			t.Errorf("validateKanji() error = %v", err)
		}
		if _, exists := seen["市場"]; !exists {
			t.Error("確認済みの漢字が重複確認用マップに登録されていない")
		}
	})

	t.Run("指定した文字列長と異なる漢字はエラーになる", func(t *testing.T) {
		err := validateKanji("市場に", 2, make(map[string]struct{}))
		if err == nil || err.Error() != "漢字の文字列長が不正です" {
			t.Errorf("validateKanji() error = %v, want 漢字の文字列長が不正です", err)
		}
	})

	t.Run("同じ漢字が重複している場合はエラーになる", func(t *testing.T) {
		seen := map[string]struct{}{"市場": {}}
		err := validateKanji("市場", 2, seen)
		if err == nil || err.Error() != "漢字が重複しています" {
			t.Errorf("validateKanji() error = %v, want 漢字が重複しています", err)
		}
	})
}

func TestIsValidString(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		allowStrings string
		want         bool
	}{
		{
			name:         "許可文字だけで構成された文字列は確認を通る",
			value:        "あいうえお",
			allowStrings: allowOnYomiString,
			want:         true,
		},
		{
			name:         "許可文字に含まれない文字がある場合は確認を通らない",
			value:        "あxう",
			allowStrings: allowOnYomiString,
			want:         false,
		},
		{
			name:         "訓読みと熟語の読みで使用できるアルファベットは確認を通る",
			value:        "iuekstngzb",
			allowStrings: allowYomiString,
			want:         true,
		},
		{
			name:         "空文字列は許可文字列で構成されていると判定されない",
			value:        "",
			allowStrings: allowOnYomiString,
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidString(tt.value, tt.allowStrings); got != tt.want {
				t.Errorf("isValidString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateYomi(t *testing.T) {
	t.Run("許可された文字種の読みは確認を通る", func(t *testing.T) {
		if err := validateYomi([]string{"いちば", "しじょう"}, allowYomiString); err != nil {
			t.Errorf("validateYomi() error = %v", err)
		}
	})

	t.Run("不正な文字種を含む読みはエラーになる", func(t *testing.T) {
		err := validateYomi([]string{"市場@"}, allowYomiString)
		if err == nil || !strings.Contains(err.Error(), "市場@") {
			t.Errorf("validateYomi() error = %v, want 不正な読みを含むエラー", err)
		}
	})

	t.Run("複数の不正な読みはエラーにすべて列挙される", func(t *testing.T) {
		err := validateYomi([]string{"abc", "123"}, allowYomiString)
		if err == nil || !strings.Contains(err.Error(), "abc") || !strings.Contains(err.Error(), "123") {
			t.Errorf("validateYomi() error = %v, want abcと123を含むエラー", err)
		}
	})
}

func TestValidateIndex(t *testing.T) {
	t.Run("代表読みと索引文字がプロファイルに一致する場合は確認を通る", func(t *testing.T) {
		entry := testEntry{kanji: "市場", indexYomi: "いちば", indexKey: "い"}
		if err := validateIndex(entry, []string{"あ", "い", "う"}); err != nil {
			t.Errorf("validateIndex() error = %v", err)
		}
	})

	t.Run("代表読みが設定されていない場合はエラーになる", func(t *testing.T) {
		entry := testEntry{kanji: "市場", indexYomi: "", indexKey: ""}
		err := validateIndex(entry, []string{"い"})
		if err == nil || err.Error() != "読みは1つ以上必要です" {
			t.Errorf("validateIndex() error = %v, want 読みは1つ以上必要です", err)
		}
	})

	t.Run("索引文字がプロファイルに含まれない場合はエラーになる", func(t *testing.T) {
		entry := testEntry{kanji: "市場", indexYomi: "いちば", indexKey: "い"}
		err := validateIndex(entry, []string{"あ", "か"})
		if err == nil || err.Error() != "索引が正しくありません" {
			t.Errorf("validateIndex() error = %v, want 索引が正しくありません", err)
		}
	})
}
