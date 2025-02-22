package main

import (
	"embed"
	"flag"
	"github.com/anhgelus/golatt"
)

type Data struct {
	Websites []*Website
}

type Website struct {
	Name string
	URL  string
}

//go:embed templates
var templates embed.FS

var g *golatt.Golatt

var dev bool

func init() {
	flag.BoolVar(&dev, "dev", false, "Run in development mode")
}

// THIS WILL BE REFACTORED WITH GOLATT 0.3.0

func (d *Data) ModuloEq(i int, mod int, eq int) bool {
	return i%mod == eq
}

// THIS WILL BE REFACTORED WITH GOLATT 0.3.0

func (d *Website) ModuloEq(i int, mod int, eq int) bool {
	return i%mod == eq
}

func main() {
	flag.Parse()
	g = golatt.New(templates)
	g.DefaultSeoData = &golatt.SeoData{
		Image:       "",
		Description: "",
		Domain:      "ring.nouveauprintemps.org",
	}
	g.FormatTitle = func(t string) string {
		return t + " - Webring - Nouveau Printemps"
	}
	g.Templates = append(g.Templates,
		"templates/base/*.gohtml",
	)
	g.NewTemplate("index",
		"/",
		"Home",
		"",
		"Nouveau Printemps' Webring",
		Data{
			Websites: []*Website{
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world/"},
			},
		}).
		Handle()
	if dev {
		g.StartServer(":8000")
	} else {
		g.StartServer(":80")
	}
}
