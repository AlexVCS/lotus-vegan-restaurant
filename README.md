# Lotus Vegan Restaurant

A modern one-page website for Lotus Vegan Restaurant in Pinellas Park, Florida. Built in Go with [GoFastr](https://github.com/DonaldMurillo/gofastr), using the restaurant's own menu content and photography.

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

The standalone menu is available at <http://localhost:8080/menu>.
