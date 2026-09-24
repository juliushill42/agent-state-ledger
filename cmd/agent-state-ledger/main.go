package main

import (
	"encoding/json"
	"fmt"
	"github.com/titanu/agent-state-ledger/internal/ledger"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		usage()
	}
	path := os.Args[2]
	switch os.Args[1] {
	case "append":
		if len(os.Args) != 6 {
			usage()
		}
		x, e := ledger.Append(path, os.Args[3], os.Args[4], json.RawMessage(os.Args[5]))
		if e != nil {
			die(e)
		}
		b, _ := json.MarshalIndent(x, "", "  ")
		fmt.Println(string(b))
	case "verify":
		es, e := ledger.Load(path)
		if e != nil {
			die(e)
		}
		if e = ledger.Verify(es); e != nil {
			die(e)
		}
		fmt.Printf("VALID entries=%d watermark=%s\n", len(es), ledger.Watermark)
	case "list":
		es, e := ledger.Load(path)
		if e != nil {
			die(e)
		}
		for _, x := range es {
			b, _ := json.Marshal(x)
			fmt.Println(string(b))
		}
	default:
		usage()
	}
}
func usage() {
	fmt.Fprintln(os.Stderr, "usage: agent-state-ledger append FILE AGENT EVENT JSON | verify FILE | list FILE")
	os.Exit(2)
}
func die(e error) { fmt.Fprintln(os.Stderr, "ERROR:", e); os.Exit(1) }
