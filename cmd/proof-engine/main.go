package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/yigitcan-ozturk/pamilanga-proof-engine/internal/proof"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "verify" {
		fmt.Fprintln(os.Stderr, "Usage: proof-engine verify --manifest proof.json")
		os.Exit(2)
	}
	fs := flag.NewFlagSet("verify", flag.ContinueOnError)
	manifest := fs.String("manifest", "proof.json", "path to evidence manifest")
	if err := fs.Parse(os.Args[2:]); err != nil {
		os.Exit(2)
	}
	result := proof.Verify(*manifest)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(result)
	if result.Status != "PASS" {
		os.Exit(1)
	}
}
