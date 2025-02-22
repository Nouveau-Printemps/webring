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
	Name              string     `toml:"name"`
	URL               string     `toml:"url"`
	Description       []string   `toml:"description"`
	JoinTheRingPath   string     `toml:"join_the_ring_path"`
	LegalMentionsPath string     `toml:"legal_mentions_path"`
	Websites          []*Website `toml:"websites"`
}

type Website struct {
	Name string
	URL  string
}

// THIS WILL BE REFACTORED WITH GOLATT 0.3.0

func (c *Config) ModuloEq(i int, mod int, eq int) bool {
	return i%mod == eq
}

func (c *Config) GetJoinTheRing() template.HTML {
	return c.get(c.JoinTheRingPath, &joinTheRing)
}

func (c *Config) GetLegalMentions() template.HTML {
	return c.get(c.LegalMentionsPath, &legalMentions)
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

// THIS WILL BE REFACTORED WITH GOLATT 0.3.0

func (w *Website) ModuloEq(i int, mod int, eq int) bool {
	return i%mod == eq
}
