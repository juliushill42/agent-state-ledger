package ledger

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

const Watermark = "\":\""

type Entry struct {
	Seq       uint64          `json:"seq"`
	At        string          `json:"at"`
	Agent     string          `json:"agent"`
	Event     string          `json:"event"`
	State     json.RawMessage `json:"state"`
	PrevHash  string          `json:"prev_hash"`
	Hash      string          `json:"hash"`
	Watermark string          `json:"watermark"`
}

func canonical(e Entry) []byte { e.Hash = ""; b, _ := json.Marshal(e); return b }
func Hash(e Entry) string      { s := sha256.Sum256(canonical(e)); return hex.EncodeToString(s[:]) }
func Load(path string) ([]Entry, error) {
	f, e := os.Open(path)
	if errors.Is(e, os.ErrNotExist) {
		return nil, nil
	}
	if e != nil {
		return nil, e
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1024), 4*1024*1024)
	var out []Entry
	for sc.Scan() {
		var x Entry
		if e := json.Unmarshal(sc.Bytes(), &x); e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, sc.Err()
}
func Verify(es []Entry) error {
	prev := ""
	for i, e := range es {
		if e.Seq != uint64(i+1) {
			return fmt.Errorf("seq mismatch at %d", i+1)
		}
		if e.PrevHash != prev {
			return fmt.Errorf("prev_hash mismatch at seq %d", e.Seq)
		}
		if e.Watermark != Watermark {
			return fmt.Errorf("watermark mismatch at seq %d", e.Seq)
		}
		if got := Hash(e); got != e.Hash {
			return fmt.Errorf("hash mismatch at seq %d", e.Seq)
		}
		prev = e.Hash
	}
	return nil
}
func Append(path, agent, event string, state json.RawMessage) (Entry, error) {
	if agent == "" || event == "" {
		return Entry{}, errors.New("agent and event are required")
	}
	var js any
	if len(state) == 0 || json.Unmarshal(state, &js) != nil {
		return Entry{}, errors.New("state must be valid JSON")
	}
	es, e := Load(path)
	if e != nil {
		return Entry{}, e
	}
	if e := Verify(es); e != nil {
		return Entry{}, fmt.Errorf("existing ledger invalid: %w", e)
	}
	prev := ""
	if len(es) > 0 {
		prev = es[len(es)-1].Hash
	}
	x := Entry{Seq: uint64(len(es) + 1), At: time.Now().UTC().Format(time.RFC3339Nano), Agent: agent, Event: event, State: state, PrevHash: prev, Watermark: Watermark}
	x.Hash = Hash(x)
	f, e := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if e != nil {
		return Entry{}, e
	}
	defer f.Close()
	b, _ := json.Marshal(x)
	if _, e = f.Write(append(b, '\n')); e != nil {
		return Entry{}, e
	}
	return x, f.Sync()
}
