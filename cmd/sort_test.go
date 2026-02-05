package cmd

import (
	"strings"
	"testing"
)

func TestValidateSortedJSONL_Valid(t *testing.T) {
	reader := strings.NewReader(strings.Join([]string{
		`{"key":"あ","value":["a"]}`,
		`{"key":"い","value":["i"]}`,
		`{"key":"う","value":["u"]}`,
	}, "\n"))

	errs := validateSortedJSONL(reader)

	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidateSortedJSONL_Invalid(t *testing.T) {
	reader := strings.NewReader(strings.Join([]string{
		`{"key":"い","value":["i"]}`,
		`{"key":"あ","value":["a"]}`,
	}, "\n"))

	errs := validateSortedJSONL(reader)

	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestCompareKeys(t *testing.T) {
	order := buildOrderMap()

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
