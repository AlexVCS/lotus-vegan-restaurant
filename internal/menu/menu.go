package menu

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	BeginMarker = "      <!-- BEGIN GENERATED MENU (edit menu.json, run go generate) -->"
	EndMarker   = "      <!-- END GENERATED MENU -->"
)

var (
	slugPattern    = regexp.MustCompile(`^[a-z0-9-]+$`)
	validDietary   = map[string]bool{"gf": true, "sf": true}
	htmlEsc        = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")
	dietaryLabels  = map[string]struct{ label, text string }{
		"gf": {"GF", "Gluten-free"},
		"sf": {"SF", "Soy-free"},
	}
)

type Menu struct {
	AllergyNote AllergyNote `json:"allergyNote"`
	Groups      []Group     `json:"groups"`
}

type AllergyNote struct {
	Label string `json:"label"`
	Text  string `json:"text"`
}

type Group struct {
	Slug     string `json:"slug"`
	Title    string `json:"title"`
	NavLabel string `json:"navLabel"`
	Note     string `json:"note,omitempty"`
	Items    []Item   `json:"items"`
}

type Item struct {
	Title         string   `json:"title"`
	Price         string   `json:"price,omitempty"`
	Description   string   `json:"description,omitempty"`
	Dietary       []string `json:"dietary,omitempty"`
	CopyPreserved bool     `json:"copyPreserved,omitempty"`
}

func esc(s string) string {
	return htmlEsc.Replace(s)
}

func Load(path string) (*Menu, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()

	var m Menu
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *Menu) Validate() error {
	seen := make(map[string]bool, len(m.Groups))
	for _, g := range m.Groups {
		if g.Slug == "" {
			return fmt.Errorf("group slug is required")
		}
		if !slugPattern.MatchString(g.Slug) {
			return fmt.Errorf("invalid slug %q", g.Slug)
		}
		if seen[g.Slug] {
			return fmt.Errorf("duplicate slug %q", g.Slug)
		}
		seen[g.Slug] = true
		if strings.TrimSpace(g.Title) == "" {
			return fmt.Errorf("group %q: title is required", g.Slug)
		}
		if strings.TrimSpace(g.NavLabel) == "" {
			return fmt.Errorf("group %q: navLabel is required", g.Slug)
		}
		for _, item := range g.Items {
			if strings.TrimSpace(item.Title) == "" {
				return fmt.Errorf("group %q: item title is required", g.Slug)
			}
			for _, d := range item.Dietary {
				if !validDietary[d] {
					return fmt.Errorf("group %q item %q: invalid dietary %q", g.Slug, item.Title, d)
				}
			}
		}
	}
	if strings.TrimSpace(m.AllergyNote.Label) == "" || strings.TrimSpace(m.AllergyNote.Text) == "" {
		return fmt.Errorf("allergyNote label and text are required")
	}
	return nil
}

func RepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}

func ExtractBetweenMarkers(content string) (string, error) {
	lines := strings.Split(content, "\n")
	begin, end, err := findMarkerLines(lines)
	if err != nil {
		return "", err
	}
	return strings.Join(lines[begin+1:end], "\n"), nil
}

func findMarkerLines(lines []string) (begin, end int, err error) {
	begin = -1
	end = -1
	for i, line := range lines {
		switch line {
		case BeginMarker:
			if begin >= 0 {
				return -1, -1, fmt.Errorf("duplicate begin marker at line %d", i+1)
			}
			begin = i
		case EndMarker:
			if end >= 0 {
				return -1, -1, fmt.Errorf("duplicate end marker at line %d", i+1)
			}
			end = i
		}
	}
	if begin < 0 {
		return -1, -1, fmt.Errorf("begin marker not found")
	}
	if end < 0 {
		return -1, -1, fmt.Errorf("end marker not found")
	}
	if end <= begin {
		return -1, -1, fmt.Errorf("end marker before begin marker")
	}
	return begin, end, nil
}

func SpliceFile(path, generated string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	content := string(data)
	lines := strings.Split(content, "\n")
	begin, end, err := findMarkerLines(lines)
	if err != nil {
		return false, fmt.Errorf("%s: %w", path, err)
	}

	newMiddle := generated
	if newMiddle != "" && !strings.HasSuffix(newMiddle, "\n") {
		newMiddle += "\n"
	}

	before := strings.Join(lines[:begin+1], "\n")
	after := strings.Join(lines[end:], "\n")
	newContent := before + "\n" + newMiddle + after

	if newContent == content {
		return false, nil
	}
	return true, os.WriteFile(path, []byte(newContent), 0o644)
}

func Render(m *Menu) string {
	var b strings.Builder
	tooltipN := 0

	b.WriteString("      <!-- The full menu is static HTML so it remains available without JavaScript. -->\n")
	b.WriteString("      <nav class=\"menu-jump\" id=\"menu-jump\" aria-label=\"Menu categories\">\n")
	for _, g := range m.Groups {
		b.WriteString("        <a href=\"#menu-")
		b.WriteString(esc(g.Slug))
		b.WriteString("\">")
		b.WriteString(esc(g.NavLabel))
		b.WriteString("</a>\n")
	}
	b.WriteString("      </nav>\n")
	b.WriteString("      <div class=\"menu-groups\" id=\"menu-groups\">\n")

	for gi, g := range m.Groups {
		b.WriteString("      <details class=\"menu-group\" id=\"menu-")
		b.WriteString(esc(g.Slug))
		b.WriteString("\"")
		if gi == 0 {
			b.WriteString(" open")
		}
		b.WriteString(">\n")
		b.WriteString("        <summary>\n")
		fmt.Fprintf(&b, "          <span class=\"menu-group-index\">%02d</span>\n", gi+1)
		b.WriteString("          <h3>")
		b.WriteString(esc(g.Title))
		b.WriteString("</h3>\n")
		b.WriteString("        </summary>\n")
		b.WriteString("        <div class=\"menu-items\">\n")

		if g.Note != "" {
			b.WriteString("          <p class=\"menu-category-note\">")
			b.WriteString(esc(g.Note))
			b.WriteString("</p>\n")
		}

		for _, item := range g.Items {
			renderItem(&b, item, &tooltipN)
		}

		b.WriteString("        </div>\n")
		b.WriteString("      </details>\n")
	}

	b.WriteString("      </div>\n")
	b.WriteString("      <div class=\"menu-notice\" id=\"allergy-note\">\n")
	b.WriteString("        <span>")
	b.WriteString(esc(m.AllergyNote.Label))
	b.WriteString("</span>\n")
	b.WriteString("        <p>")
	b.WriteString(esc(m.AllergyNote.Text))
	b.WriteString("</p>\n")
	b.WriteString("      </div>")

	return b.String()
}

func renderItem(b *strings.Builder, item Item, tooltipN *int) {
	b.WriteString("          <article class=\"menu-item")
	if item.CopyPreserved {
		b.WriteString(" menu-item-copy-preserved")
	}
	b.WriteString("\">\n")
	b.WriteString("            <h4 class=\"menu-item-name\">\n")
	b.WriteString("              <span class=\"menu-item-heading\">\n")
	b.WriteString("                <span class=\"menu-item-title\">")
	b.WriteString(esc(item.Title))
	b.WriteString("</span>")
	if len(item.Dietary) > 0 {
		b.WriteString("<span class=\"menu-item-dietary\">")
		for _, code := range item.Dietary {
			*tooltipN++
			info := dietaryLabels[code]
			fmt.Fprintf(b,
				`<span class="dietary-badge dietary-badge-tooltip dietary-badge-%s" tabindex="0" aria-label="%s" aria-describedby="dietary-tooltip-%d"><span aria-hidden="true">%s</span><span class="dietary-tooltip" id="dietary-tooltip-%d" role="tooltip">%s</span></span>`,
				code, info.label, *tooltipN, info.label, *tooltipN, info.text,
			)
		}
		b.WriteString("</span>")
	}
	b.WriteString("\n              </span>\n")
	if item.Price != "" {
		b.WriteString("            <i class=\"dots\" aria-hidden=\"true\"></i>\n")
		b.WriteString("            <span class=\"menu-price\">")
		b.WriteString(esc(item.Price))
		b.WriteString("</span>\n")
	}
	b.WriteString("            </h4>\n")
	if item.Description != "" {
		b.WriteString("          <p>")
		b.WriteString(esc(item.Description))
		b.WriteString("</p>\n")
	}
	b.WriteString("          </article>\n")
}
