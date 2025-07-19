// Coordinates reading files, parsing content, rendering templates, and writing output.
package builder

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
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
	tmpl, err := template.ParseFiles("templates/base.html")
	if err != nil {
		return fmt.Errorf("template parse error: %w", err)
	}

	var pages []parser.Page

	err = filepath.WalkDir("content", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return walkErr
		}

		input, readErr := os.ReadFile(path)
		if readErr != nil {
			return fmt.Errorf("read file error: %w", readErr)
		}

		page, parseErr := parser.ParseMarkdownWithFrontMatter(input, path)
		if parseErr != nil {
			return fmt.Errorf("parse markdown error in %s: %w", path, parseErr)
		}

		pages = append(pages, page)

		// write individual pages
		outputPath := page.Path
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

	if err != nil {
		return err
	}

	// sort by date descending
	sort.SliceStable(pages, func(i, j int) bool {
		return pages[i].Time.After(pages[j].Time)
	})

	// render blog index
	indexTemplate, parseErr := template.ParseFiles("templates/blog_index.html")
	if parseErr != nil {
		return fmt.Errorf("error parsing blog index template: %w", parseErr)
	}

	indexPath := "output/blog/index.html"
	if err := os.MkdirAll(filepath.Dir(indexPath), os.ModePerm); err != nil {
		return err
	}

	out, err := os.Create(indexPath)
	if err != nil {
		return err
	}
	defer out.Close()

	return indexTemplate.Execute(out, pages)
}
