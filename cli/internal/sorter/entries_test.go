package sorter

import (
	"reflect"
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

func TestSortEntries(t *testing.T) {
	t.Run("代表読み順にソートしたスライスを返す", func(t *testing.T) {
		entries := []testEntry{
			{kanji: "三", indexYomi: "さん"},
			{kanji: "愛", indexYomi: "あい"},
			{kanji: "二", indexYomi: "に"},
		}

		got := SortEntries(entries)
		want := []testEntry{
			{kanji: "愛", indexYomi: "あい"},
			{kanji: "三", indexYomi: "さん"},
			{kanji: "二", indexYomi: "に"},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("SortEntries() = %v, want %v", got, want)
		}
	})

	t.Run("元のスライスを変更しない", func(t *testing.T) {
		entries := []testEntry{
			{kanji: "三", indexYomi: "さん"},
			{kanji: "愛", indexYomi: "あい"},
		}
		original := append([]testEntry(nil), entries...)

		SortEntries(entries)

		if !reflect.DeepEqual(entries, original) {
			t.Errorf("SortEntries() changed entries = %v, want %v", entries, original)
		}
	})

	t.Run("同じ代表読みのエントリは元の順序を維持する", func(t *testing.T) {
		entries := []testEntry{
			{kanji: "市場", indexYomi: "いちば"},
			{kanji: "位置", indexYomi: "いちば"},
			{kanji: "一番", indexYomi: "いちば"},
		}

		got := SortEntries(entries)
		for i, entry := range entries {
			if got[i] != entry {
				t.Errorf("SortEntries() stable order changed at index %d: got %v, want %v", i, got[i], entry)
			}
		}
	})

	t.Run("空のスライスをソートできる", func(t *testing.T) {
		got := SortEntries([]testEntry{})
		if got == nil || len(got) != 0 {
			t.Errorf("SortEntries() = %v, want non-nil empty slice", got)
		}
	})
}

func TestIsSorted(t *testing.T) {
	t.Run("代表読み順ならソート済みと判定する", func(t *testing.T) {
		entries := []testEntry{
			{kanji: "愛", indexYomi: "あい"},
			{kanji: "二", indexYomi: "に"},
		}

		if errs := IsSorted(entries); len(errs) != 0 {
			t.Errorf("IsSorted() errors = %v, want none", errs)
		}
	})

	t.Run("代表読み順でない場合はソート済みと判定しない", func(t *testing.T) {
		entries := []testEntry{
			{kanji: "二", indexYomi: "に"},
			{kanji: "愛", indexYomi: "あい"},
		}

		got := IsSorted(entries)
		if len(got) != 2 {
			t.Fatalf("IsSorted() errors length = %d, want 2", len(got))
		}
		if got[0].Error() != "  - [0] 二（に） -> 愛（あい）" {
			t.Errorf("IsSorted()[0] = %q, want unsorted difference", got[0])
		}
		if got[1].Error() != "  - [1] 愛（あい） -> 二（に）" {
			t.Errorf("IsSorted()[1] = %q, want unsorted difference", got[1])
		}
	})
}
