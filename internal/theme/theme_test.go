package theme

import (
	"math"
	"strconv"
	"strings"
	"testing"
)

func TestLotusValidates(t *testing.T) {
	th := Lotus() // panics through MustValidate on a bad token
	if err := th.Validate(); err != nil {
		t.Fatal(err)
	}
	css := th.CSSCustomProperties()
	for _, want := range []string{
		"--color-primary: #173c2d",
		"--color-background: #f5f0e5",
		"--color-brand-lime: #cae75d",
		`:root[data-color-scheme="dark"]`,
		"--color-primary: #d5ed74",
		"DM Sans",
	} {
		if !strings.Contains(css, want) {
			t.Errorf("token CSS missing %q", want)
		}
	}
}

// TestTextContrast holds every text tier at WCAG AA (4.5:1) on the page
// and surface colours, in both schemes. MustValidate checks the filled
// button pairs; this covers the reading text.
func TestTextContrast(t *testing.T) {
	th := Lotus()
	light := map[string]string{
		"background":   th.Colors.Background.Value,
		"surface":      th.Colors.Surface.Value,
		"surface-soft": th.Colors.SurfaceSoft.Value,
		"text":         th.Colors.Text.Value,
		"text-muted":   th.Colors.TextMuted.Value,
		"text-subtle":  th.Colors.TextSubtle.Value,
		"accent":       th.Colors.Accent.Value,
		"danger":       th.Colors.Danger.Value,
	}
	dark := map[string]string{}
	for k, v := range light {
		dark[k] = v
	}
	for k, v := range th.DarkColors {
		dark[k] = v
	}
	for name, pal := range map[string]map[string]string{"light": light, "dark": dark} {
		for _, bg := range []string{"background", "surface"} {
			for _, fg := range []string{"text", "text-muted", "text-subtle", "accent", "danger"} {
				if r := contrast(t, pal[fg], pal[bg]); r < 4.5 {
					t.Errorf("%s: %s on %s is %.2f:1, want >= 4.5", name, fg, bg, r)
				}
			}
		}
	}
	// The dietary badges carry their ink in both schemes.
	for _, fill := range []string{th.DarkColors["dietary-gf"], th.DarkColors["dietary-sf"]} {
		if r := contrast(t, th.DarkColors["dietary-ink"], fill); r < 4.5 {
			t.Errorf("dark dietary badge ink on %s is %.2f:1", fill, r)
		}
	}
}

func contrast(t *testing.T, a, b string) float64 {
	t.Helper()
	la, lb := luminance(t, a), luminance(t, b)
	if la < lb {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}

func luminance(t *testing.T, hex string) float64 {
	t.Helper()
	h := strings.TrimPrefix(hex, "#")
	if len(h) != 6 {
		t.Fatalf("not a #rrggbb colour: %q", hex)
	}
	var c [3]float64
	for i := range 3 {
		v, err := strconv.ParseUint(h[i*2:i*2+2], 16, 8)
		if err != nil {
			t.Fatalf("bad colour %q: %v", hex, err)
		}
		s := float64(v) / 255
		if s <= 0.03928 {
			c[i] = s / 12.92
		} else {
			c[i] = math.Pow((s+0.055)/1.055, 2.4)
		}
	}
	return 0.2126*c[0] + 0.7152*c[1] + 0.0722*c[2]
}
