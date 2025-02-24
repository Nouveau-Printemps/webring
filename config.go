package main

import (
	"html/template"
	"os"
)

var (
	joinTheRing   = ""
	legalMentions = ""
)

type Config struct {
	Name                 string       `toml:"name"`
	URL                  string       `toml:"url"`
	Description          []string     `toml:"description"`
	JoinTheRingPath      string       `toml:"join_the_ring_path"`
	LegalInformationPath string       `toml:"legal_information_path"`
	FaviconPath          string       `toml:"favicon_path"`
	Translation          *Translation `toml:"translation"`
	Websites             []*Website   `toml:"website"`
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
	return c.get(c.JoinTheRingPath, &joinTheRing)
}

func (c *Config) GetLegalInformation() template.HTML {
	return c.get(c.LegalInformationPath, &legalMentions)
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
