package generator

import (
	"time"

	"github.com/goccy/go-yaml"
)

type Frontmatter struct {
	Title     string    `yaml:"title"`
	Excerpt   string    `yaml:"excerpt"`
	Timestamp time.Time `yaml:"timestamp"`
}

func ParseFrontmatter(v string) (*Frontmatter, error) {
	var fm Frontmatter
	if err := yaml.Unmarshal([]byte(v), &fm); err != nil {
		return nil, err
	}

	return &fm, nil
}
