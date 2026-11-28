package model

import (
	"bytes"
	"testing"
)

func TestDictionaryEncode(t *testing.T) {
	t.Run("読みごとにまとめ、キー昇順で出力する", func(t *testing.T) {
		d := Dictionary{
			"はし": {"橋"},
			"あお": {"青"},
			"あい": {"愛", "藍"},
		}

		var buf bytes.Buffer
		if err := d.Encode(&buf); err != nil {
			t.Fatal(err)
		}

		want := "{\"あい\":[\"愛\",\"藍\"],\"あお\":[\"青\"],\"はし\":[\"橋\"]}\n"
		if buf.String() != want {
			t.Errorf("Encode() = %q, want %q", buf.String(), want)
		}
	})

	t.Run("空の場合は空オブジェクトを出力する", func(t *testing.T) {
		var buf bytes.Buffer
		if err := (Dictionary(nil)).Encode(&buf); err != nil {
			t.Fatal(err)
		}
		if buf.String() != "{}\n" {
			t.Errorf("Encode() = %q, want {}", buf.String())
		}
	})
}

func TestDictionaryMerge(t *testing.T) {
	a := Dictionary{"あい": {"愛"}}
	a.Merge(Dictionary{"あい": {"藍"}, "あお": {"青"}})

	if got := a["あい"]; len(got) != 2 || got[1] != "藍" {
		t.Errorf("Merge() あい = %v", got)
	}
	if len(a.Yomis()) != 2 {
		t.Errorf("Yomis() = %v", a.Yomis())
	}
}
