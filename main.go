package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/anhgelus/golatt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
)

var (
	//go:embed templates
	templates embed.FS
	//go:embed dist
	assets embed.FS
)

var g *golatt.Golatt

var (
	dev            bool
	generateConfig bool
	configPath     string
	port           uint
)

func init() {
	flag.BoolVar(&dev, "dev", false, "Run in development mode")
	flag.BoolVar(&generateConfig, "generate-config", false, "Generate default config file")
	flag.StringVar(&configPath, "config", "config.toml", "Webring's config file")
	flag.UintVar(&port, "port", 80, "Port to use")
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

	if dev {
		g = golatt.New(templates, os.DirFS("public/"), os.DirFS("dist/"))
	} else {
		g = golatt.New(templates, os.DirFS("public/"), golatt.UsableEmbedFS("dist", assets))
	}
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
		g.StartServer(fmt.Sprintf(":%d", port))
	}
}

func genConfigToStdOut() {
	cfg := Config{
		Name:                 "My Webring",
		URL:                  "ring.example.org",
		JoinTheRingPath:      "join_ring.html",
		LegalInformationPath: "legal_information.html",
		FaviconPath:          "logo.png",
		Description:          []string{"Welcome to my fantastic webring!", "It has all my friends' websites and mine!"},
		Websites: []*Website{
			{
				Name: "Example",
				URL:  "https://example.org/",
			},
		},
		Translation: &Translation{
			RandomWebsite:        "Go to a random website",
			InternalPages:        "Internal pages:",
			HomePage:             "home",
			JoinTheRingPage:      "join the ring",
			LegalInformationPage: "legal information",
		},
	}
	err := toml.NewEncoder(os.Stdout).Encode(cfg)
	if err != nil {
		panic(err)
	}
}
