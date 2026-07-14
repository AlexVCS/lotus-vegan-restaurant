package main

import (
	"embed"
	"io/fs"
	"log"
	"os"

	staticfs "github.com/DonaldMurillo/gofastr/core/static"
	"github.com/DonaldMurillo/gofastr/framework"
)

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
