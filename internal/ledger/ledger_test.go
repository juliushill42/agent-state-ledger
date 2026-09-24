package ledger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestAppendVerify(t *testing.T) {
	p := filepath.Join(t.TempDir(), "l.jsonl")
	if _, e := Append(p, "a", "start", json.RawMessage(`{"x":1}`)); e != nil {
		t.Fatal(e)
	}
	if _, e := Append(p, "a", "done", json.RawMessage(`{"x":2}`)); e != nil {
		t.Fatal(e)
	}
	es, e := Load(p)
	if e != nil || len(es) != 2 {
		t.Fatal(e, len(es))
	}
	if e = Verify(es); e != nil {
		t.Fatal(e)
	}
}
func TestTamper(t *testing.T) {
	p := filepath.Join(t.TempDir(), "l.jsonl")
	Append(p, "a", "start", json.RawMessage(`{"x":1}`))
	b, _ := os.ReadFile(p)
	for i := range b {
		if b[i] == '1' {
			b[i] = '9'
			break
		}
	}
	os.WriteFile(p, b, 0644)
	es, e := Load(p)
	if e != nil {
		t.Fatal(e)
	}
	if Verify(es) == nil {
		t.Fatal("tamper not detected")
	}
}
