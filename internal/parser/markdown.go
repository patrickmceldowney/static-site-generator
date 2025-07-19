// Uses a markdown library to convert .md to HTML and extract front matter.
package parser

import (
	"errors"
	"html/template"
	"strings"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"gopkg.in/yaml.v3"
)

// yaml:"-" means we shouldn't parse it from YAML, we will generate it by code
type Page struct {
	Title    string        `yaml:"title"`
	Date     string        `yaml:"date"`
	Tags     []string      `yaml:"tags"`
	Slug     string        `yaml:"slug"`
	Template string        `yaml:"template"`
	Content  template.HTML `yaml:"-"`
	Path     string        `yaml:"-"`
}

func ParseMarkdownWithFrontMatter(content []byte, path string) (Page, error) {
	var page Page

	contentStr := string(content)
	if !strings.HasPrefix(contentStr, "---") {
		return page, errors.New("missing front matter")
	}

	parts := strings.SplitN(contentStr, "---", 3)
	if len(parts) < 3 {
		return page, errors.New("invalid front matter format")
	}

	frontMatter := parts[1]
	body := parts[2]

	// parse YAML into page fields
	if err := yaml.Unmarshal([]byte(frontMatter), &page); err != nil {
		return page, err
	}

	// compute fallback slug/path if not provided
	relPath := strings.TrimSuffix(path, ".md")
	if page.Slug == "" {
		page.Slug = filepath.Base(relPath)
	}

	page.Path = filepath.ToSlash(strings.Replace(relPath, "content", "output", 1) + ".html")

	// markdown to HTML
	renderer := html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags})
	htmlContent := markdown.ToHTML([]byte(body), nil, renderer)
	page.Content = template.HTML(htmlContent)

	return page, nil
}
