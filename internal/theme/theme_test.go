package theme

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/DonaldMurillo/gofastr/core-ui/style"
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
		"--color-dietary-sf: #f5ab92",
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
// and surface colours, in both schemes, for the page theme and every
// band. MustValidate checks the filled button pairs; this covers the
// reading text.
func TestTextContrast(t *testing.T) {
	for name, th := range map[string]style.Theme{
		"page": Lotus(), "lime band": LimeBand, "gold band": GoldBand, "forest band": ForestBand,
	} {
		light := map[string]string{
			"background":  th.Colors.Background.Value,
			"surface":     th.Colors.Surface.Value,
			"text":        th.Colors.Text.Value,
			"text-muted":  th.Colors.TextMuted.Value,
			"text-subtle": th.Colors.TextSubtle.Value,
			"accent":      th.Colors.Accent.Value,
			"danger":      th.Colors.Danger.Value,
		}
		dark := map[string]string{}
		for k, v := range light {
			dark[k] = v
		}
		for k, v := range th.DarkColors {
			dark[k] = v
		}
		for scheme, pal := range map[string]map[string]string{"light": light, "dark": dark} {
			for _, bg := range []string{"background", "surface"} {
				for _, fg := range []string{"text", "text-muted", "text-subtle", "accent", "danger"} {
					if r := contrast(t, pal[fg], pal[bg]); r < 4.5 {
						t.Errorf("%s %s: %s on %s is %.2f:1, want >= 4.5", name, scheme, fg, bg, r)
					}
				}
			}
		}
	}
}

// TestDietaryBadgesReadOnTheMenuBand checks the dietary badges, which
// are soft-tone status badges (their colour as text on a 15% tint of
// itself) and only ever sit on the forest menu band.
func TestDietaryBadgesReadOnTheMenuBand(t *testing.T) {
	th := ForestBand
	for scheme, bg := range map[string]string{"light": th.Colors.Background.Value, "dark": th.DarkColors["background"]} {
		for name, ink := range map[string]string{
			"gf": pick(scheme, "brand-lime", th.DarkColors),
			"sf": pick(scheme, "dietary-sf", th.DarkColors),
		} {
			chip := mix(t, ink, bg, 0.15)
			if r := contrast(t, ink, chip); r < 4.5 {
				t.Errorf("%s %s badge: %s on its tint %s is %.2f:1", scheme, name, ink, chip, r)
			}
		}
	}
}

// pick reads a brand token for a scheme: the dark map in dark mode,
// the generated light value otherwise.
func pick(scheme, key string, dark map[string]string) string {
	if scheme == "dark" {
		return dark[key]
	}
	switch key {
	case "brand-lime":
		return Tokens.Colors.BrandLime.Value
	case "dietary-sf":
		return Tokens.Colors.DietarySf.Value
	}
	return ""
}

func mix(t *testing.T, fg, bg string, alpha float64) string {
	t.Helper()
	f, b := rgb(t, fg), rgb(t, bg)
	out := "#"
	for i := range 3 {
		v := alpha*f[i] + (1-alpha)*b[i]
		out += strconv.FormatInt(int64(math.Round(v))+0x100, 16)[1:]
	}
	return out
}

func rgb(t *testing.T, hex string) [3]float64 {
	t.Helper()
	h := strings.TrimPrefix(hex, "#")
	var c [3]float64
	for i := range 3 {
		v, err := strconv.ParseUint(h[i*2:i*2+2], 16, 8)
		if err != nil {
			t.Fatalf("bad colour %q: %v", hex, err)
		}
		c[i] = float64(v)
	}
	return c
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
