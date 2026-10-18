package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/DonaldMurillo/gofastr/framework"
	"github.com/DonaldMurillo/gofastr/framework/uihost"

	"github.com/AlexVCS/lotus-vegan/internal/site"
)

//go:generate go run ./cmd/genmenu
//go:embed public
var siteFiles embed.FS

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app := buildApp()
	log.Printf("Lotus Vegan is serving on :%s", port)
	if err := app.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}

// buildApp wires the site into a framework app without binding a port,
// so tests can drive it in-process.
func buildApp() *framework.App {
	public, err := fs.Sub(siteFiles, "public")
	if err != nil {
		log.Fatal(err)
	}

	// The pages render through GoFastr's UI host; files under public/
	// (photos, fonts, favicons, lotus.js) are served by the host's static
	// fallback before page rendering, so /assets/* keeps its URLs.
	host := uihost.New(site.New(), site.HostOptions()...)
	host.SetStaticFS(public)

	app := framework.NewApp(
		framework.WithConfig(framework.AppConfig{Name: "lotus-vegan"}),
	)
	site.MountRoutes(app.Router())
	// The UI host answers /favicon.ico with 204 unless it generates the
	// icons itself (WithAppIcon); the site ships its own multi-size .ico,
	// so it gets an explicit route.
	app.Router().Get("/favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, public, "favicon.ico")
	}))
	app.Use(host.RouteMatchMiddleware())
	app.Mount(host)
	return app
}
