// Uses a markdown library to convert .md to HTML and extract front matter.
package parser

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
	"html/template"
	

	"github.com/gomarkdown/markdown"
	"github.com/patrickmceldowney/static-site-generator/internal/models"
	"gopkg.in/yaml.v3"
)

type frontMatter struct {
	Title  string   `yaml:"title"`
	Date   string   `yaml:"date"`
	Tags   []string `yaml:"tags"`
	Slug   string   `yaml:"slug,omitempty"`
	Layout string   `yaml:"layout"`
}

func ParseMarkdown(path, inputDir string) (models.Page, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return models.Page{}, fmt.Errorf("reading markdown: %w", err)
	}

	content := string(data)
	var fm frontMatter
	var body string

	if strings.HasPrefix(content, "---") {
		parts := strings.SplitN(content, "---", 3)
		if len(parts) == 3 {
			err := yaml.Unmarshal([]byte(parts[1]), &fm)
			if err != nil {
				return models.Page{}, fmt.Errorf("unmarshalling front matter: %w", err)
			}

			body = strings.TrimSpace(parts[2])
		} else {
			return models.Page{}, fmt.Errorf("invalid front matter format")
		}
	} else {
		// No front matter — fallback
		body = content
	}

	// Parse data
	var date time.Time
	if rawDate := fm.Date; rawDate != "" {
		date, err = time.Parse("2006-01-02", rawDate)
		if err != nil {
			return models.Page{}, fmt.Errorf("invalid date format in %s: %w", path, err)
		}
	}

	// set layout
	layout := fm.Layout
	if layout == "" {
		layout = "page"
	}

	// use slug from front matter or fallback to filename
	slug := fm.Slug
	if slug == "" {
		relPath, err := filepath.Rel(inputDir, path)
		if err != nil {
			return models.Page{}, fmt.Errorf("resolving relative path: %w", err)
		}

		relPath = strings.TrimSuffix(relPath, filepath.Ext(relPath))
		slug = filepath.ToSlash(relPath)
	}

	html := markdown.ToHTML([]byte(body), nil, nil)
	outputPath := filepath.Join(slug, "index.html")

	return models.Page{
		Title:       fm.Title,
		Date:        date,
		Tags:        fm.Tags,
		Slug:        slug,
		Content:     body,
		HTMLContent: template.HTML(html),
		OutputPath:  outputPath,
		Layout:			 layout,
	}, nil

}
