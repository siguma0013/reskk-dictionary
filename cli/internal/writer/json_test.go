package writer

import (
	"os"
	"path/filepath"
	"testing"

	"siguma0013/reskk-dictionary/internal/model"
)

func TestWriteJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.json")
	if err := os.WriteFile(path, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}

	dict := model.Dictionary{"はし": {"橋", "箸"}, "あい": {"愛"}}
	if err := WriteJSON(dict, path); err != nil {
		t.Fatal(err)
	}

	got, _ := os.ReadFile(path)
	want := "{\"あい\":[\"愛\"],\"はし\":[\"橋\",\"箸\"]}\n"
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestWriteJSONNil(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.json")
	if err := WriteJSON(nil, path); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != "{}\n" {
		t.Errorf("got %q", got)
	}
}
