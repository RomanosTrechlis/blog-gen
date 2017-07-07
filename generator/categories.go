package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"
	"strings"
)

// Category holds the data for a category
type Category struct {
	Name  string
	Link  string
	Count int
}

// CategoryByCountDesc sorts the cats
type CategoryByCountDesc []*Category

// CategoriesGenerator object
type CategoriesGenerator struct {
	Config *CategoriesConfig
}

// CategoriesConfig holds the category's config
type CategoriesConfig struct {
	CatPostsMap     map[string][]*Post
	Template        *template.Template
	Destination     string
	ThemeFolder     string
	BlogTitle       string
	Author, BlogURL string
}

var catTemplatePath string

// Generate creates the categories page
func (g *CategoriesGenerator) Generate() error {
	fmt.Println("\tGenerating Categories...")
	catTemplatePath = g.Config.ThemeFolder + "categories.html"
	catPostsMap := g.Config.CatPostsMap
	t := g.Config.Template
	destination := g.Config.Destination
	catsPath := fmt.Sprintf("%s/categories", destination)
	err := clearAndCreateDestination(catsPath)
	if err != nil {
		return err
	}
	err = generateCatIndex(catPostsMap, t, catsPath, g.Config.Author, g.Config.BlogURL, g.Config.BlogTitle)
	if err != nil {
		return err
	}
	for cat, catPosts := range catPostsMap {
		catPagePath := fmt.Sprintf("%s/%s", catsPath, cat)
		err = generateCatPage(cat, catPosts, t, catPagePath)
		if err != nil {
			return err
		}
	}
	fmt.Println("\tFinished generating Categories...")
	return nil
}

func generateCatIndex(catPostsMap map[string][]*Post, t *template.Template, destination, author, blogURL, blogTitle string) error {
	tmpl, err := getTemplate(catTemplatePath)
	if err != nil {
		return err
	}
	categories := []*Category{}
	for cat, posts := range catPostsMap {
		categories = append(categories, &Category{Name: cat, Link: getCatLink(cat), Count: len(posts)})
	}
	sort.Sort(CategoryByCountDesc(categories))
	buf := bytes.Buffer{}
	err = tmpl.Execute(&buf, categories)
	if err != nil {
		return fmt.Errorf("error executing template %s: %v", catTemplatePath, err)
	}
	err = writeIndexHTML(destination, "Categories", author, blogURL, blogTitle, template.HTML(buf.String()), t)
	if err != nil {
		return err
	}
	return nil
}

func generateCatPage(cat string, posts []*Post, t *template.Template, destination string) error {
	err := clearAndCreateDestination(destination)
	if err != nil {
		return err
	}
	lg := ListingGenerator{&ListingConfig{
		Posts:       posts,
		Template:    t,
		Destination: destination,
		PageTitle:   cat,
	}}
	err = lg.Generate()
	if err != nil {
		return err
	}
	return nil
}

func getCatLink(cat string) string {
	return fmt.Sprintf("/categories/%s/", strings.ToLower(cat))
}

func (t CategoryByCountDesc) Len() int {
	return len(t)
}

func (t CategoryByCountDesc) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t CategoryByCountDesc) Less(i, j int) bool {
	return t[i].Count > t[j].Count
}
