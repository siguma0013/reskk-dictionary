package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestWalkJSONL_FindsJSONLFiles(t *testing.T) {
	d := t.TempDir()
	if err := os.WriteFile(filepath.Join(d, "a.jsonl"), []byte("{}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(d, "sub"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(d, "sub", "c.jsonl"), []byte("{}\n"), 0644); err != nil {
		t.Fatal(err)
	}

	results, err := WalkJSONL(d, func(path string, r io.Reader) ([]error, error) {
		// consume reader to ensure it's valid
		_, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 jsonl files, got %d", len(results))
	}

	m := map[string]bool{}
	for _, res := range results {
		m[filepath.Base(res.Path)] = true
	}
	if !m["a.jsonl"] {
		t.Fatalf("a.jsonl not found in results")
	}
	if !m["c.jsonl"] {
		t.Fatalf("c.jsonl not found in results")
	}
}

func TestWalkJSONL_ProcessErrorIsRecorded(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "err.jsonl")
	if err := os.WriteFile(p, []byte("{}\n"), 0644); err != nil {
		t.Fatal(err)
	}

	results, err := WalkJSONL(d, func(path string, r io.Reader) ([]error, error) {
		if filepath.Base(path) == "err.jsonl" {
			return nil, fmt.Errorf("process error")
		}
		return nil, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].CriticalError == nil {
		t.Fatalf("expected CriticalError to be set when process returns error")
	}
}

func TestWalkJSONL_ErrorOnNonJSONLFile(t *testing.T) {
	d := t.TempDir()
	if err := os.WriteFile(filepath.Join(d, "a.jsonl"), []byte("{}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(d, "b.txt"), []byte("ignored"), 0644); err != nil {
		t.Fatal(err)
	}

	results, _ := WalkJSONL(d, func(path string, r io.Reader) ([]error, error) { return nil, nil })
	if results[1].CriticalError == nil {
		t.Fatalf("expected WalkJSONL to return an error when non-jsonl files exist")
	}
}

func TestWalkJSONL_NonExistentRootReturnsError(t *testing.T) {
	_, err := WalkJSONL("/path/does/not/exist", func(path string, r io.Reader) ([]error, error) { return nil, nil })
	if err == nil {
		t.Fatalf("expected WalkJSONL to return an error for non-existent root")
	}
}
