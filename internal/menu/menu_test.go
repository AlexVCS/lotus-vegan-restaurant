package menu

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestGeneratedMenuIsFresh(t *testing.T) {
	root, err := RepoRoot()
	if err != nil {
		t.Fatal(err)
	}

	m, err := Load(filepath.Join(root, "menu.json"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := GoSource(m)
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(root, "internal", "menu", "data.gen.go")
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		off := firstDivergence(string(got), string(want))
		t.Fatalf("%s: stale — run `go generate ./...` (first divergence at offset %d)", path, off)
	}
}

// TestDataMatchesJSON checks the compiled-in menu is the one menu.json
// describes, item for item, prices and dietary codes included.
func TestDataMatchesJSON(t *testing.T) {
	root, err := RepoRoot()
	if err != nil {
		t.Fatal(err)
	}
	m, err := Load(filepath.Join(root, "menu.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := Data.Validate(); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(m, Data) {
		t.Fatal("menu.Data differs from menu.json — run `go generate ./...`")
	}
	priced := 0
	for _, g := range Data.Groups {
		for _, it := range g.Items {
			if it.Price == "" {
				continue
			}
			priced++
			if !strings.Contains(it.Price, "$") {
				t.Errorf("%s / %s: price %q has no $ prefix", g.Slug, it.Title, it.Price)
			}
		}
	}
	if priced == 0 {
		t.Fatal("no priced items in the menu")
	}
}

func TestValidate(t *testing.T) {
	valid := &Menu{
		AllergyNote: AllergyNote{Label: "Allergy note", Text: "text"},
		Groups: []Group{{
			Slug: "appetizers", Title: "Appetizers", NavLabel: "Appetizers",
			Items: []Item{{Title: "Rolls", Dietary: []string{"gf"}}},
		}},
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid menu: %v", err)
	}

	dup := &Menu{
		AllergyNote: AllergyNote{Label: "Allergy note", Text: "text"},
		Groups: []Group{
			{Slug: "dup", Title: "A", NavLabel: "A"},
			{Slug: "dup", Title: "B", NavLabel: "B"},
		},
	}
	if err := dup.Validate(); err == nil {
		t.Fatal("expected duplicate slug error")
	}

	badDietary := &Menu{
		AllergyNote: AllergyNote{Label: "Allergy note", Text: "text"},
		Groups: []Group{{
			Slug: "x", Title: "X", NavLabel: "X",
			Items: []Item{{Title: "Item", Dietary: []string{"vegan"}}},
		}},
	}
	if err := badDietary.Validate(); err == nil {
		t.Fatal("expected invalid dietary error")
	}
}

func firstDivergence(a, b string) int {
	min := len(a)
	if len(b) < min {
		min = len(b)
	}
	for i := 0; i < min; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	if len(a) != len(b) {
		return min
	}
	return -1
}
