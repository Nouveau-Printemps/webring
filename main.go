package main

import (
	"embed"
	"flag"
	"github.com/anhgelus/golatt"
)

//go:embed templates
var templates embed.FS

var g *golatt.Golatt

var dev bool

func init() {
	flag.BoolVar(&dev, "dev", false, "Run in development mode")
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
		struct {
			Websites []struct {
				Name string
				URL  string
			}
		}{
			Websites: []struct {
				Name string
				URL  string
			}{
				{"Anhgelus Morhtuuzh", "https://now.anhgelus.world"},
			},
		}).
		Handle()
	if dev {
		g.StartServer(":8000")
	} else {
		g.StartServer(":80")
	}
}
