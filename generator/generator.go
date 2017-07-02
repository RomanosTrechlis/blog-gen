package generator

import (
	"html/template"
	"time"
)

// Generator creates content.
type Generator interface {
	Generate() error
}

// Meta is a data container for Metadata
type Meta struct {
	Title      string
	Short      string
	Date       string
	Tags       []string
	Categories []string
	ParsedDate time.Time
}

// IndexData is a data container for the landing page
type IndexData struct {
	HTMLTitle     string
	PageTitle     string
	Content       template.HTML
	Year          int
	Name          string
	CanonicalLink string
	PageNum       int
	PrevPageNum   int
	NextPageNum   int
	URL           string
	IsPost        bool
}
