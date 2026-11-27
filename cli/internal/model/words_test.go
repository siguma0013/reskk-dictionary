package model

import (
	"reflect"
	"testing"
)

func TestWordEntryAccessors(t *testing.T) {
	tests := []struct {
		name          string
		entry         WordEntry
		wantKanji     string
		wantYomi      []string
		wantIndexYomi string
		wantIndexKey  string
	}{
		{
			name: "読みが複数ある",
			entry: WordEntry{
				Kanji: "市場",
				Yomi:  []string{"いちば", "しじょう"},
			},
			wantKanji:     "市場",
			wantYomi:      []string{"いちば", "しじょう"},
			wantIndexYomi: "いちば",
			wantIndexKey:  "い",
		},
		{
			name:          "異常系：両方の読みが空",
			entry:         WordEntry{},
			wantIndexYomi: "",
			wantIndexKey:  "",
		},
		{
			name: "異常系：先頭の読みが空文字",
			entry: WordEntry{
				Yomi: []string{"", "よみ"},
			},
			wantYomi:      []string{"", "よみ"},
			wantIndexYomi: "",
			wantIndexKey:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.entry.Kanji; got != tt.wantKanji {
				t.Errorf("Kanji = %q, want %q", got, tt.wantKanji)
			}

			if got := tt.entry.GetKanji(); got != tt.wantKanji {
				t.Errorf("GetKanji() = %q, want %q", got, tt.wantKanji)
			}

			if got := tt.entry.Yomi; !reflect.DeepEqual(got, tt.wantYomi) {
				t.Errorf("Yomi = %v, want %v", got, tt.wantYomi)
			}

			if got := tt.entry.GetIndexYomi(); got != tt.wantIndexYomi {
				t.Errorf("GetIndexYomi() = %q, want %q", got, tt.wantIndexYomi)
			}

			if got := tt.entry.GetIndexKey(); got != tt.wantIndexKey {
				t.Errorf("GetIndexKey() = %q, want %q", got, tt.wantIndexKey)
			}
		})
	}
}
