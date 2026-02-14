package cmd

import (
	"siguma0013/reskk-dictionary/internal/dictionary"
	"strings"
	"testing"
)

func TestCheckSorted_Valid(t *testing.T) {
	reader := strings.NewReader(strings.Join([]string{
		`{"key":"あ","value":["a"]}`,
		`{"key":"い","value":["i"]}`,
		`{"key":"う","value":["u"]}`,
	}, "\n"))

	errs := checkSorted(reader)

	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestCheckSorted_Invalid(t *testing.T) {
	reader := strings.NewReader(strings.Join([]string{
		`{"key":"い","value":["i"]}`,
		`{"key":"あ","value":["a"]}`,
	}, "\n"))

	errs := checkSorted(reader)

	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestCompareKeys(t *testing.T) {
	order := dictionary.SortOrder()

	if compareKeys("か", "が", order) < 0 {
		// "か" should come before "が"
	} else {
		t.Fatalf("expected \"か\" < \"が\"")
	}

	if compareKeys("き", "きゃ", order) < 0 {
		// "き" shorter should come before "きゃ"
	} else {
		t.Fatalf("expected \"き\" < \"きゃ\"")
	}
}

func TestSortData(t *testing.T) {
	order := dictionary.SortOrder()

	reader := strings.NewReader(strings.Join([]string{
		`{"key":"きのう","value":["機能"]}`,
		`{"key":"あき","value":["秋"]}`,
		`{"key":"あい","value":["愛"]}`,
	}, "\n"))

	sorted, _ := sortData(reader, order)

	if len(sorted) != 3 {
		t.Fatalf("expected 3 records, got %d", len(sorted))
	}

	expectedKeys := []string{"あい", "あき", "きのう"}

	for i, e := range sorted {
		if e.Key != expectedKeys[i] {
			t.Fatalf("expected key %q at index %d, got %q", expectedKeys[i], i, e.Key)
		}
	}

}
