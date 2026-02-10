package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMergeCommand_Minimal(t *testing.T) {
	// Setup temp dir
	d := t.TempDir()
	oldwd, _ := os.Getwd()
	defer os.Chdir(oldwd)
	if err := os.Chdir(d); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// create files
	if err := os.WriteFile("a.jsonl", []byte(strings.Join([]string{
		`{"key":"k1","value":["A"]}`,
		`{"key":"k2","value":["A2"]}`,
	}, "\n")+"\n"), 0o644); err != nil {
		t.Fatalf("write a.jsonl: %v", err)
	}

	// b.jsonl contains duplicate key k2 which should be ignored because a.jsonl appears first
	if err := os.WriteFile("b.jsonl", []byte(strings.Join([]string{
		`{"key":"k2","value":["B2"]}`,
		`{"value":["no_key"]}`,
	}, "\n")+"\n"), 0o644); err != nil {
		t.Fatalf("write b.jsonl: %v", err)
	}

	// create merge_order.yml
	if err := os.WriteFile("merge_order.yml", []byte(strings.Join([]string{
		"files:",
		"  - \"a.jsonl\"",
		"  - \"b.jsonl\"",
	}, "\n")+"\n"), 0o644); err != nil {
		t.Fatalf("write merge_order.yml: %v", err)
	}

	// Run command
	if err := mergeCmd.RunE(nil, nil); err != nil {
		t.Fatalf("merge command failed: %v", err)
	}

	// check output file
	out, err := os.ReadFile(filepath.Join(d, "merged.jsonl"))
	if err != nil {
		t.Fatalf("read merged.jsonl: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	// expect keyed entries sorted by key (k1,k2) then the unkeyed record
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines in merged output, got %d: %v", len(lines), lines)
	}

	if !strings.Contains(lines[0], `"k1"`) {
		t.Fatalf("expected first line to contain k1, got %s", lines[0])
	}
	if !strings.Contains(lines[1], `"k2"`) {
		t.Fatalf("expected second line to contain k2, got %s", lines[1])
	}
	if !strings.Contains(lines[2], `"no_key"`) {
		t.Fatalf("expected third line to contain unkeyed record, got %s", lines[2])
	}
}
