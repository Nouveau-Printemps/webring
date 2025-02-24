package main

import (
	"embed"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/anhgelus/golatt"
	"html/template"
	"math/rand"
	"net/http"
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
		Description: cfg.Description[0],
		Domain:      cfg.URL,
	}
	g.NotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		g.Render(w, "not_found", &golatt.TemplateData{
			Title: "404",
			SEO:   &golatt.SeoData{},
			Data:  &cfg,
		})
	}
	g.FormatTitle = func(t string) string {
		return t + " - " + cfg.Name
	}
	g.Templates = append(g.Templates,
		"templates/base/*.gohtml",
	)
	g.TemplateFuncMap = template.FuncMap{
		"moduloEq": func(i int, mod int, eq int) bool {
			return i%mod == eq
		},
	}

	g.NewTemplate("index",
		"/",
		"Home",
		"",
		cfg.Description[0],
		&cfg).
		Handle()
	g.NewTemplate("join",
		"/join",
		"Join the Ring",
		"",
		"Read the instructions to join the ring!",
		&cfg).
		Handle()
	g.NewTemplate("legal",
		"/legal",
		"Legal information",
		"",
		"Legal information about the ring.",
		&cfg).
		Handle()
	g.HandleFunc("/random", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, cfg.Websites[rand.Intn(len(cfg.Websites))].URL, http.StatusFound)
	})

	if dev {
		g.StartServer(":8000")
	} else {
		g.StartServer(":80")
	}
}

func genConfigToStdOut() {
	cfg := Config{
		Name:                 "My Webring",
		URL:                  "ring.example.org",
		JoinTheRingPath:      "join_ring.html",
		LegalInformationPath: "legal_information.html",
		Description:          []string{"Welcome to my fantastic webring!", "It has all my friends' websites and mine!"},
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
