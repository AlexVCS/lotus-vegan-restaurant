// Command genmenu validates menu.json and writes internal/menu/data.gen.go,
// the Go value the site renders the menu from. Run it with
// `go generate ./...` after editing menu.json.
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/AlexVCS/lotus-vegan/internal/menu"
)

func main() {
	root, err := menu.RepoRoot()
	if err != nil {
		log.Fatal(err)
	}

	m, err := menu.Load(filepath.Join(root, "menu.json"))
	if err != nil {
		log.Fatal(err)
	}

	src, err := menu.GoSource(m)
	if err != nil {
		log.Fatal(err)
	}

	target := filepath.Join(root, "internal", "menu", "data.gen.go")
	if old, err := os.ReadFile(target); err == nil && bytes.Equal(old, src) {
		fmt.Fprintln(os.Stderr, "menu data unchanged")
		return
	}
	if err := os.WriteFile(target, src, 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stderr, "updated %s\n", target)
}
