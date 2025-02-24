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
	Name                 string     `toml:"name"`
	URL                  string     `toml:"url"`
	Description          []string   `toml:"description"`
	JoinTheRingPath      string     `toml:"join_the_ring_path"`
	LegalInformationPath string     `toml:"legal_information_path"`
	Websites             []*Website `toml:"websites"`
}

type Website struct {
	Name string
	URL  string
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
