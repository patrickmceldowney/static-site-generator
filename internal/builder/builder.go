// internal/builder/builder.go
package builder

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"

	"github.com/patrickmceldowney/static-site-generator/internal/models"
	"github.com/patrickmceldowney/static-site-generator/internal/parser"
)

// Build coordinates reading files, parsing, rendering, and writing output.
func Build(inputDir, outputDir, templateDir string) error {
	// Parse all .html files in the template directory.
	// When using template inheritance, ParseGlob makes all defined templates
	// available to each other.
	templates, err := template.ParseGlob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	var pages []models.Page
	err = filepath.Walk(inputDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		page, parserErr := parser.ParseMarkdown(path, inputDir)
		if parserErr != nil {
			return parserErr
		}
		pages = append(pages, page)
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking markdown files: %w", err)
	}

	// Sort by date descending.
	sort.SliceStable(pages, func(i, j int) bool {
		return pages[i].Date.After(pages[j].Date)
	})

	// Create output directory.
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Render individual pages.
	for _, page := range pages {
		outputPath := filepath.Join(outputDir, page.Slug, "index.html")
		if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
			return err
		}
		f, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer f.Close()

		fmt.Sprintf("%s", page.Layout)
		err = templates.ExecuteTemplate(f, page.Layout, page)
		if err != nil {
			return fmt.Errorf("error rendering %s (%s): %w", page.Layout, page.Slug, err)
		}
	}

	// Render the index page.
	indexPath := filepath.Join(outputDir, "index.html")
	indexFile, err := os.Create(indexPath)
	if err != nil {
		return err
	}
	defer indexFile.Close()

	// The data for the index page.
	indexData := struct {
		Title string
		Pages []models.Page
	}{
		Title: "Home",
		Pages: pages,
	}

	// Execute the "index" template.
	err = templates.ExecuteTemplate(indexFile, "index", indexData)
	if err != nil {
		return fmt.Errorf("error rendering index: %w", err)
	}

	return nil
}
