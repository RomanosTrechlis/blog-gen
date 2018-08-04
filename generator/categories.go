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

// categoriesGenerator struct
type categoriesGenerator struct {
	catPostsMap map[string][]*post
	template    *template.Template
	destination string
	siteInfo    *config.SiteInformation
}

// Generate creates the categories page
func (g *categoriesGenerator) Generate() (err error) {
	fmt.Println("\tGenerating Categories...")
	catPostsMap := g.catPostsMap
	destination := g.destination
	catsPath := fmt.Sprintf("%s/categories", destination)
	err = clearAndCreateDestination(catsPath)
	if err != nil {
		return err
	}
	err = g.generateCatIndex()
	if err != nil {
		return err
	}
	for cat, catPosts := range catPostsMap {
		catPagePath := fmt.Sprintf("%s/%s", catsPath, cat)
		err = g.generateCatPage(cat, catPosts, catPagePath)
		if err != nil {
			return err
		}
	}
	fmt.Println("\tFinished generating Categories...")
	return nil
}

func (g *categoriesGenerator) generateCatIndex() (err error) {
	catTemplatePath := g.siteInfo.ThemeFolder + "categories.html"
	tmpl, err := getTemplate(catTemplatePath)
	if err != nil {
		return err
	}
	categories := []*Category{}
	for cat, posts := range g.catPostsMap {
		categories = append(categories, &Category{Name: cat, Link: getCatLink(cat), Count: len(posts)})
	}
	sort.Sort(categoryByCountDesc(categories))
	buf := bytes.Buffer{}
	err = tmpl.Execute(&buf, categories)
	if err != nil {
		return fmt.Errorf("error executing template %s: %v", catTemplatePath, err)
	}

	c := htmlConfig{
		path:       fmt.Sprintf("%s/categories", g.destination),
		pageTitle:  "Categories",
		pageNum:    0,
		maxPageNum: 0,
		isPost:     false,
		temp:       g.template,
		content:    template.HTML(buf.String()),
		siteInfo:   g.siteInfo,
	}
	err = c.writeHTML()
	if err != nil {
		return err
	}
	return nil
}

func (g *categoriesGenerator) generateCatPage(cat string, posts []*post, path string) (err error) {
	err = clearAndCreateDestination(path)
	if err != nil {
		return err
	}
	lg := listingGenerator{
		posts:       posts,
		template:    g.template,
		destination: path,
		pageTitle:   cat,
		siteInfo:    g.siteInfo,
	}
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

// categoryByCountDesc sorts the cats
type categoryByCountDesc []*Category

func (t categoryByCountDesc) Len() (l int) {
	return len(t)
}

func (t categoryByCountDesc) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t categoryByCountDesc) Less(i, j int) (l bool) {
	return t[i].Count > t[j].Count
}
