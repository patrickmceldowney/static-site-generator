// Coordinates reading files, parsing content, rendering templates, and writing output.
package builder

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"html/template"

	"github.com/patrickmceldowney/static-site-generator/internal/parser"
)

type Page struct {
	Title   string
	Date    string
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

		page, parseErr := parser.ParseMarkdownWithFrontMatter(input)
		if parseErr != nil {
			return fmt.Errorf("parse markdown error in %s: %w", path, parseErr)
		}

		// Get relative path from content/ → e.g., "blog/post.md"
		relPath, pathErr := filepath.Rel("content", path)
		if pathErr != nil {
			return fmt.Errorf("path error: %w", pathErr)
		}

		// Replace .md with .html and place into /output
		outputPath := filepath.Join("output", strings.ReplaceAll(relPath, ".md", ".html"))

		if osErr := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); osErr != nil {
			return fmt.Errorf("mkdir error: %w", osErr)
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
