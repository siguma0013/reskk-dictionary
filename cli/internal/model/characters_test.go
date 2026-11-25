package model

import (
	"reflect"
	"testing"
)

func TestCharacterEntryAccessors(t *testing.T) {
	tests := []struct {
		name          string
		entry         CharacterEntry
		wantKanji     string
		wantOnYomi    []string
		wantKunYomi   []string
		wantIndexYomi string
		wantIndexKey  string
	}{
		{
			name: "音読みと訓読みがある",
			entry: CharacterEntry{
				Kanji:   "見",
				OnYomi:  []string{"けん"},
				KunYomi: []string{"みe", "みr"},
			},
			wantKanji:     "見",
			wantOnYomi:    []string{"けん"},
			wantKunYomi:   []string{"みe", "みr"},
			wantIndexYomi: "けん",
			wantIndexKey:  "け",
		},
		{
			name: "音読みが空",
			entry: CharacterEntry{
				Kanji:   "嵐",
				KunYomi: []string{"あらし"},
			},
			wantKanji:     "嵐",
			wantKunYomi:   []string{"あらし"},
			wantIndexYomi: "あらし",
			wantIndexKey:  "あ",
		},
		{
			name:          "異常系：両方の読みが空",
			entry:         CharacterEntry{},
			wantIndexYomi: "",
			wantIndexKey:  "",
		},
		{
			name: "異常系：先頭の音読みが空文字",
			entry: CharacterEntry{
				OnYomi:  []string{"", "おんよみ"},
				KunYomi: []string{"くんよみ"},
			},
			wantOnYomi:    []string{"", "おんよみ"},
			wantKunYomi:   []string{"くんよみ"},
			wantIndexYomi: "",
			wantIndexKey:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.entry.Kanji; got != tt.wantKanji {
				t.Errorf("Kanji = %q, want %q", got, tt.wantKanji)
			}

			if got := tt.entry.OnYomi; !reflect.DeepEqual(got, tt.wantOnYomi) {
				t.Errorf("OnYomi = %v, want %v", got, tt.wantOnYomi)
			}

			if got := tt.entry.KunYomi; !reflect.DeepEqual(got, tt.wantKunYomi) {
				t.Errorf("KunYomi = %v, want %v", got, tt.wantKunYomi)
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
