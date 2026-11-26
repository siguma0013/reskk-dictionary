package loader

import (
	"os"
	"path/filepath"
	"testing"

	"siguma0013/reskk-dictionary/internal/model"
)

func writeYamlTestFile(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "dictionary.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("テスト用YAMLの作成に失敗: %v", err)
	}
	return path
}

func TestLoadYamlWithCharacterModel(t *testing.T) {
	path := writeYamlTestFile(t, `entries:
  - 漢字: 亜
    音読み: [あ]
    訓読み: []
`)

	got, err := LoadYaml[model.Characters](path)
	if err != nil {
		t.Fatalf("LoadYaml() error = %v", err)
	}
	if len(got.Entries) != 1 {
		t.Fatalf("Entries length = %d, want 1", len(got.Entries))
	}

	entry := got.Entries[0]
	if entry.GetKanji() != "亜" {
		t.Errorf("GetKanji() = %q, want %q", entry.GetKanji(), "亜")
	}
	if got := entry.OnYomi; len(got) != 1 || got[0] != "あ" {
		t.Errorf("GetOnYomi() = %v, want [あ]", got)
	}
	if got := entry.KunYomi; len(got) != 0 {
		t.Errorf("GetKunYomi() = %v, want empty", got)
	}
}

func TestLoadYamlWithWordModel(t *testing.T) {
	path := writeYamlTestFile(t, `entries:
  - 漢字: 市場
    読み: [いちば, しじょう]
`)

	got, err := LoadYaml[model.Words](path)
	if err != nil {
		t.Fatalf("LoadYaml() error = %v", err)
	}
	if len(got.Entries) != 1 {
		t.Fatalf("Entries length = %d, want 1", len(got.Entries))
	}

	entry := got.Entries[0]
	if entry.GetKanji() != "市場" {
		t.Errorf("GetKanji() = %q, want %q", entry.GetKanji(), "市場")
	}
	if got := entry.Yomi; len(got) != 2 || got[0] != "いちば" || got[1] != "しじょう" {
		t.Errorf("GetYomi() = %v, want [いちば しじょう]", got)
	}
}

func TestLoadYamlErrors(t *testing.T) {
	t.Run("ファイルが存在しない", func(t *testing.T) {
		_, err := LoadYaml[model.Characters](filepath.Join(t.TempDir(), "missing.yaml"))
		if err == nil {
			t.Fatal("LoadYaml() error = nil, want error")
		}
	})

	t.Run("YAML構文が不正", func(t *testing.T) {
		path := writeYamlTestFile(t, "entries:\n  - 漢字: [\n")
		_, err := LoadYaml[model.Characters](path)
		if err == nil {
			t.Fatal("LoadYaml() error = nil, want error")
		}
	})

	t.Run("未知フィールド", func(t *testing.T) {
		path := writeYamlTestFile(t, `entries:
  - 漢字: 市場
    読み: [いちば]
    未知: true
`)
		_, err := LoadYaml[model.Characters](path)
		if err == nil {
			t.Fatal("LoadYaml() error = nil, want error")
		}
	})
}

func TestSaveYamlRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.yaml")
	want := model.Words{Entries: []model.WordEntry{
		{Kanji: "市場", Yomi: []string{"いちば", "しじょう"}},
	}}

	if err := SaveYaml(path, want); err != nil {
		t.Fatalf("SaveYaml() error = %v", err)
	}

	got, err := LoadYaml[model.Words](path)
	if err != nil {
		t.Fatalf("LoadYaml() error = %v", err)
	}
	if len(got.Entries) != 1 || got.Entries[0].Kanji != "市場" || len(got.Entries[0].Yomi) != 2 {
		t.Errorf("round trip = %+v, want %+v", got, want)
	}
}
