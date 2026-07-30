package main

import (
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

	generated := menu.Render(m)
	targets := []string{
		filepath.Join(root, "public", "index.html"),
		filepath.Join(root, "public", "menu", "index.html"),
	}

	changed := false
	for _, path := range targets {
		didChange, err := menu.SpliceFile(path, generated)
		if err != nil {
			log.Fatal(err)
		}
		if didChange {
			changed = true
			fmt.Fprintf(os.Stderr, "updated %s\n", path)
		}
	}
	if !changed {
		fmt.Fprintln(os.Stderr, "menu HTML unchanged")
	}
}
