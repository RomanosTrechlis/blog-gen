package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"
	"strings"

	"github.com/RomanosTrechlis/blog-generator/config"
)

// Category holds the data for a category
type Category struct {
	Name  string
	Link  string
	Count int
}

// categoryByCountDesc sorts the cats
type categoryByCountDesc []*Category

// categoriesGenerator object
type categoriesGenerator struct {
	config *categoriesConfig
}

// CategoriesConfig holds the category's config
type categoriesConfig struct {
	catPostsMap map[string][]*post
	template    *template.Template
	destination string
	siteInfo    *config.SiteInformation
}

var catTemplatePath string

// Generate creates the categories page
func (g *categoriesGenerator) Generate() (err error) {
	fmt.Println("\tGenerating Categories...")
	siteInfo := g.config.siteInfo
	catTemplatePath = siteInfo.ThemeFolder + "categories.html"
	catPostsMap := g.config.catPostsMap
	t := g.config.template
	destination := g.config.destination
	catsPath := fmt.Sprintf("%s/categories", destination)
	err = clearAndCreateDestination(catsPath)
	if err != nil {
		return err
	}
	err = generateCatIndex(catPostsMap, t, catsPath, siteInfo)
	if err != nil {
		return err
	}
	for cat, catPosts := range catPostsMap {
		catPagePath := fmt.Sprintf("%s/%s", catsPath, cat)
		err = generateCatPage(cat, catPosts, t, catPagePath, siteInfo)
		if err != nil {
			return err
		}
	}
	fmt.Println("\tFinished generating Categories...")
	return nil
}

func generateCatIndex(catPostsMap map[string][]*post, t *template.Template,
	destination string, siteInfo *config.SiteInformation) (err error) {
	tmpl, err := getTemplate(catTemplatePath)
	if err != nil {
		return err
	}
	categories := []*Category{}
	for cat, posts := range catPostsMap {
		categories = append(categories, &Category{Name: cat, Link: getCatLink(cat), Count: len(posts)})
	}
	sort.Sort(categoryByCountDesc(categories))
	buf := bytes.Buffer{}
	err = tmpl.Execute(&buf, categories)
	if err != nil {
		return fmt.Errorf("error executing template %s: %v", catTemplatePath, err)
	}
	err = writeIndexHTML(destination, "Categories", template.HTML(buf.String()), t, siteInfo)
	if err != nil {
		return err
	}
	return nil
}

func generateCatPage(cat string, posts []*post, t *template.Template,
	destination string, siteInfo *config.SiteInformation) (err error) {
	err = clearAndCreateDestination(destination)
	if err != nil {
		return err
	}
	lg := listingGenerator{&listingConfig{
		posts:       posts,
		template:    t,
		destination: destination,
		pageTitle:   cat,
		siteInfo:    siteInfo,
	}}
	err = lg.Generate()
	if err != nil {
		return err
	}
	return nil
}

func getCatLink(cat string) (link string) {
	link = fmt.Sprintf("/categories/%s/", strings.ToLower(cat))
	return link
}

func (t categoryByCountDesc) Len() (l int) {
	return len(t)
}

func (t categoryByCountDesc) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t categoryByCountDesc) Less(i, j int) (l bool) {
	return t[i].Count > t[j].Count
}
