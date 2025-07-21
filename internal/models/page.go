package models

import (
	"time"
	"html/template"
)

type Page struct {
	Title        string
	Slug         string
	Date         time.Time
	Tags         []string
	HTMLContent  template.HTML
	Content      string
	OutputPath   string
	Layout       string
}
