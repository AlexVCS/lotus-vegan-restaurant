package menu

import (
	"os"
	"path/filepath"
	"strconv"
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
	rendered := Render(m)

	targets := []string{
		filepath.Join(root, "public", "index.html"),
		filepath.Join(root, "public", "menu", "index.html"),
	}
	for _, path := range targets {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		extracted, err := ExtractBetweenMarkers(string(data))
		if err != nil {
			t.Fatalf("%s: %v", path, err)
		}
		if extracted != rendered {
			off := firstDivergence(extracted, rendered)
			t.Fatalf("%s: stale — run `go generate ./...` (first divergence at offset %d)", path, off)
		}
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

func TestTooltipNumbering(t *testing.T) {
	root, err := RepoRoot()
	if err != nil {
		t.Fatal(err)
	}
	m, err := Load(filepath.Join(root, "menu.json"))
	if err != nil {
		t.Fatal(err)
	}
	rendered := Render(m)

	total := 0
	for _, g := range m.Groups {
		for _, item := range g.Items {
			total += len(item.Dietary)
		}
	}

	for n := 1; n <= total; n++ {
		s := strconv.Itoa(n)
		id := "id=\"dietary-tooltip-" + s + "\""
		aria := "aria-describedby=\"dietary-tooltip-" + s + "\""
		idCount := strings.Count(rendered, id)
		ariaCount := strings.Count(rendered, aria)
		if idCount != 1 {
			t.Fatalf("dietary-tooltip-%d: expected id= once, got %d", n, idCount)
		}
		if ariaCount != 1 {
			t.Fatalf("dietary-tooltip-%d: expected aria-describedby= once, got %d", n, ariaCount)
		}
	}

	next := "id=\"dietary-tooltip-" + strconv.Itoa(total+1) + "\""
	if strings.Count(rendered, next) != 0 {
		t.Fatalf("dietary-tooltip-%d: unexpected id beyond derived total %d", total+1, total)
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
