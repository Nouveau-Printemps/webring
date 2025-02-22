package main

import (
	"embed"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/anhgelus/golatt"
	"os"
)

//go:embed templates
var templates embed.FS

var g *golatt.Golatt

var (
	dev            bool
	generateConfig bool
	configPath     string
)

func init() {
	flag.BoolVar(&dev, "dev", false, "Run in development mode")
	flag.BoolVar(&generateConfig, "generate-config", false, "Generate default config file")
	flag.StringVar(&configPath, "config", "config.toml", "Webring's config file")
}

func main() {
	flag.Parse()
	if generateConfig {
		genConfigToStdOut()
		return
	}

	b, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	var cfg Config
	err = toml.Unmarshal(b, &cfg)
	if err != nil {
		panic(err)
	}

	g = golatt.New(templates)
	g.DefaultSeoData = &golatt.SeoData{
		Image:       "",
		Description: "",
		Domain:      cfg.URL,
	}
	g.FormatTitle = func(t string) string {
		return t + " - " + cfg.Name
	}
	g.Templates = append(g.Templates,
		"templates/base/*.gohtml",
	)
	g.NewTemplate("index",
		"/",
		"Home",
		"",
		cfg.Description[0],
		cfg).
		Handle()
	if dev {
		g.StartServer(":8000")
	} else {
		g.StartServer(":80")
	}
}

func genConfigToStdOut() {
	cfg := Config{
		Name:              "My Webring",
		URL:               "ring.example.org",
		JoinTheRingPath:   "/app/join_ring.html",
		LegalMentionsPath: "/app/legal_mentions.html",
		Description:       []string{"Welcome to my fantastic webring!", "It has all my friends' websites and mine!"},
		Websites: []*Website{
			{
				Name: "Example",
				URL:  "https://example.org/",
			},
		},
	}
	err := toml.NewEncoder(os.Stdout).Encode(cfg)
	if err != nil {
		panic(err)
	}
}
