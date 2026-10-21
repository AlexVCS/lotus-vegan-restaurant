# Lotus Vegan Restaurant

A modern one-page website for Lotus Vegan Restaurant in Pinellas Park, Florida. Built in Go with [GoFastr](https://github.com/DonaldMurillo/gofastr) v0.90 and its UI system, using the restaurant's own menu content and photography.

Live site: [lotus-vegan.onrender.com](https://lotus-vegan.onrender.com/)

## Run locally

Requires **Go 1.27** or newer (`go.mod` pins toolchain `go1.27.2`; with the default `GOTOOLCHAIN=auto`, an older Go downloads it automatically). GoFastr v0.90.0 requires Go 1.27.

This command requires the [GoFastr CLI](https://github.com/DonaldMurillo/gofastr#quickstart) to be installed.

```bash
gofastr dev
```

The development server automatically reloads when Go, JavaScript, CSS, or HTML files change.

Alternatively, run the project without the GoFastr CLI:

```bash
go run .
```

Open <http://localhost:8080>.

The standalone menu is available at <http://localhost:8080/menu> and the contact page at <http://localhost:8080/contact>.

## How it's built

- Pages are GoFastr screens rendered by the UI host (`internal/site`), composed from `framework/ui` components: `Hero`, `Section`, `Card`/`Grid`, `FactBox`, `LinkButton`, `Collapsible`, `StatusBadge`, `Callout`, `ThemeToggle` and `BackToTop`.
- The header and footer are owned packages (`internal/siteheader`, `internal/sitefooter`) generated with `gofastr generate package` and extended for Lotus.
- The Lotus brand is a GoFastr theme (`internal/theme`): light and dark palettes, DM Sans, the lime/gold/forest page bands, and brand tokens in `brand.tokens.css`. Dark mode uses GoFastr's `data-color-scheme` attribute and the `gofastr.colorScheme` saved preference.
- Owned style sheets (`*.style.css`) cover the few things the kit doesn't draw. After editing one, or a `*.tokens.css`, run `gofastr gen styles`.
- The menu lives in `menu.json`. After editing it, run `go generate ./...` to regenerate `internal/menu/data.gen.go`.
- Photos, fonts and favicons are served from `public/`.
