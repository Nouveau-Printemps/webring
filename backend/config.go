package backend

import (
	"html/template"
	"log/slog"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Name                 string       `toml:"name"`
	Domain               string       `toml:"domain"`
	Description          []string     `toml:"description"`
	JoinTheRingPath      string       `toml:"join_the_ring_path"`
	LegalInformationPath string       `toml:"legal_information_path"`
	FaviconPath          string       `toml:"favicon_path"`
	Translation          *Translation `toml:"translation"`
	Websites             []*Website   `toml:"website"`
	PublicFolder         string       `toml:"public_folder"`

	joinTheRing   string `toml:"-"`
	legalMentions string `toml:"-"`
}

type Website struct {
	Name string `toml:"name"`
	URL  string `toml:"url"`
}

type Translation struct {
	RandomWebsite        string `toml:"random_website"`
	InternalPages        string `toml:"internal_pages"`
	HomePage             string `toml:"home_page"`
	JoinTheRingPage      string `toml:"join_the_ring_page"`
	LegalInformationPage string `toml:"legal_information_page"`
}

func (c *Config) GetJoinTheRing() template.HTML {
	return c.get(c.JoinTheRingPath, &c.joinTheRing)
}

func (c *Config) GetLegalInformation() template.HTML {
	return c.get(c.LegalInformationPath, &c.legalMentions)
}

func (c *Config) get(path string, v *string) template.HTML {
	if *v != "" {
		return template.HTML(*v)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	*v = string(b)
	return template.HTML(*v)
}

func (c *Config) DefaultValues() {
	*c = Config{
		Name:                 "My Webring",
		Domain:               "ring.example.org",
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
}

func LoadConfig(path string) (*Config, bool) {
	b, err := os.ReadFile(path)
	var cfg Config
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
		slog.Warn("config file not found", "path", path)
		slog.Info("creating a new config file", "path", path)
		cfg.DefaultValues()
		b, err := toml.Marshal(&cfg)
		if err != nil {
			panic(err)
		}
		if err = os.WriteFile(path, b, 0644); err != nil {
			panic(err)
		}
		return nil, false
	}
	err = toml.Unmarshal(b, &cfg)
	if err != nil {
		panic(err)
	}
	return &cfg, true
}
