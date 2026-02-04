package cmd

import (
	"bytes"
	"testing"
)

func TestCheckInitials_AllValid(t *testing.T) {
	data := `{"key":"あいう","value":["v1"]}
{"key":"えお","value":["v2"]}`
	errs := checkInitials(bytes.NewBufferString(data), []string{"あ", "い", "う", "え", "お"})
	if len(errs) != 0 {
		t.Fatalf("expected 0 errors, got %d: %v", len(errs), errs)
	}
}

func TestCheckInitials_ParseError(t *testing.T) {
	data := `not json
{"key":"あい","value":["v"]}`
	errs := checkInitials(bytes.NewBufferString(data), []string{"あ", "い"})
	if len(errs) != 1 || errs[0].Error() != "parse error: 1" {
		t.Fatalf("expected parse error on line 1, got %v", errs)
	}
}

func TestCheckInitials_InitialError(t *testing.T) {
	data := `{"key":"かい","value":["v"]}
{"key":"あい","value":["v"]}`
	errs := checkInitials(bytes.NewBufferString(data), []string{"あ", "い"})
	if len(errs) != 1 || errs[0].Error() != "initial error: 1" {
		t.Fatalf("expected initial error on line 1, got %v", errs)
	}
}
