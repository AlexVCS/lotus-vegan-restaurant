package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	staticfs "github.com/DonaldMurillo/gofastr/core/static"
	"github.com/DonaldMurillo/gofastr/framework"

	lotustheme "github.com/AlexVCS/lotus-vegan/internal/theme"
)

//go:generate go run ./cmd/genmenu
//go:embed public
var siteFiles embed.FS

func main() {
	app := framework.NewApp(
		framework.WithConfig(framework.AppConfig{Name: "lotus-vegan"}),
	)

	public, err := fs.Sub(siteFiles, "public")
	if err != nil {
		log.Fatal(err)
	}

	// The GoFastr theme tokens (light :root block plus the
	// data-color-scheme="dark" re-declarations) that styles.css reads.
	// Registered before the static mount so the catch-all never sees it.
	tokenCSS := []byte(lotustheme.Lotus().CSSCustomProperties())
	app.Router().Get("/theme.css", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=300")
		_, _ = w.Write(tokenCSS)
	}))

	staticfs.Mount(app.Router(), staticfs.Config{FS: public, Prefix: ""})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Lotus Vegan is serving on :%s", port)
	if err := app.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
