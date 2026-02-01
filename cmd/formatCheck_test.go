package cmd

import (
	"strings"
	"testing"
)

func TestValidateJSONL_Valid(t *testing.T) {
	reader := strings.NewReader(`{"key": "きのう", "value": ["機能", "昨日"]}`)
	validateError := validateJSONL(reader)
	if len(validateError) != 0 {
		t.Fatalf("expected no errors, got %v", validateError)
	}
}

func TestValidateJSONL_Invalid(t *testing.T) {
	tests := []struct {
		name string
		jsonl string
	}{
		{"empty line", "\n"},
		{"trailing space", " "},
		{"invalid key", `{"key": "きのう", "value": ["機能", "昨日"], "invalid": 1}`},
		{"key type error", `{"key": 1, "value": ["機能", "昨日"]}`},
		{"value type error", `{"key": "きのう", "value": [1, 2]}`},
		{"empty key", `{"value": ["機能", "昨日"]}`},
		{"empty value", `{"key": "きのう"}`},

	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := strings.NewReader(test.jsonl)
			validateError := validateJSONL(reader)
			if len(validateError) != 1 {
				t.Fatalf("expected 1 error, got %d: %v", len(validateError), validateError)
			}
		})
	}
}
