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
	if err := os.WriteFile(path, []byte("old content"), 0o600); err != nil {
		t.Fatalf("テスト用ファイルの作成に失敗: %v", err)
	}
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

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("保存後のファイル情報取得に失敗: %v", err)
	}
	if got, want := info.Mode().Perm(), os.FileMode(0o600); got != want {
		t.Errorf("保存後の権限 = %04o, want %04o", got, want)
	}
}

func TestSaveYamlReplaceFailureKeepsOriginal(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "target.yaml")
	if err := os.Mkdir(path, 0o700); err != nil {
		t.Fatalf("テスト用ディレクトリの作成に失敗: %v", err)
	}

	err := SaveYaml(path, model.Words{})
	if err == nil {
		t.Fatal("SaveYaml() error = nil, want error")
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ディレクトリ読み込みに失敗: %v", err)
	}
	if len(entries) != 1 || entries[0].Name() != "target.yaml" || !entries[0].IsDir() {
		t.Errorf("置換失敗後のディレクトリ内容 = %v, want 元のディレクトリのみ", entries)
	}
}
