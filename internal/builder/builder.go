// Coordinates reading files, parsing content, rendering templates, and writing output.
package builder

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"html/template"
)

type Page struct {
	Title   string
	Content template.HTML
}

func Build() error {
	templatePath := "templates/base.html"
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("template parse error: %w", err)
	}

	err = filepath.WalkDir("content", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		input, readErr := os.ReadFile(path)
		if readErr != nil {
			return fmt.Errorf("read file error: %w", readErr)
		}

		opts := html.RendererOptions{Flags: html.CommonFlags}
		renderer := html.NewRenderer(opts)
		htmlContent := markdown.ToHTML(input, nil, renderer)

		// Create output path
		relPath, _ := filepath.Rel("content", path)
		outputPath := filepath.Join("output", strings.ReplaceAll(relPath, ".md", ".html"))

		if osErr := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); osErr != nil {
			return fmt.Errorf("mkdir error: %w", osErr)
		}

		// Render template with content
		page := Page{
			Title:   strings.TrimSuffix(filepath.Base(path), ".md"),
			Content: template.HTML(htmlContent),
		}

		f, createErr := os.Create(outputPath)

		if createErr != nil {
			return fmt.Errorf("create output error: %w", createErr)
		}

		defer f.Close()

		return tmpl.Execute(f, page)
	})

	return err
}
