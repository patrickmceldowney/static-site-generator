// Uses a markdown library to convert .md to HTML and extract front matter.
package parser

import (
	"errors"
	"html/template"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"gopkg.in/yaml.v3"
)

type Page struct {
	Title   string
	Date    string
	Content template.HTML
}

func ParseMarkdownWithFrontMatter(content []byte) (Page, error) {
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

	// Parse YAML into page fields
	if err := yaml.Unmarshal([]byte(frontMatter), &page); err != nil {
		return page, err
	}

	renderer := html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags})
	htmlContent := markdown.ToHTML([]byte(body), nil, renderer)

	page.Content = template.HTML(htmlContent)
	return page, nil
}
