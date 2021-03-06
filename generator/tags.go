package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"sort"

	"github.com/RomanosTrechlis/blog-gen/config"
)

// Tag holds the data for a Tag
type Tag struct {
	Name  string
	Link  string
	Count int
}

// byCountDesc sorts the tags
type byCountDesc []*Tag

// tagsGenerator object
type tagsGenerator struct {
	tagPostsMap map[string][]*post
	template    *template.Template
	siteInfo    *config.SiteInformation
}

// Generate creates the tags page
func (g *tagsGenerator) Generate() (err error) {
	fmt.Println("\tGenerating Tags...")
	siteInfo := g.siteInfo

	tagPostsMap := g.tagPostsMap
	tagsPath := filepath.Join(siteInfo.DestFolder, "tags")
	err = clearAndCreateDestination(tagsPath)
	if err != nil {
		return err
	}
	err = g.generateTagIndex()
	if err != nil {
		return err
	}
	for tag, tagPosts := range tagPostsMap {
		err := g.generateTagPage(tag, tagPosts)
		if err != nil {
			return err
		}
	}
	fmt.Println("\tFinished generating Tags...")
	return nil
}

func (g *tagsGenerator) generateTagIndex() (err error) {
	tagsPath := filepath.Join(g.siteInfo.DestFolder, "tags")
	tagsTemplatePath := filepath.Join(g.siteInfo.ThemeFolder, "tags.html")
	tmpl, err := getTemplate(tagsTemplatePath)
	if err != nil {
		return err
	}
	tags := make([]*Tag, 0)
	for tag, posts := range g.tagPostsMap {
		tags = append(tags, &Tag{Name: tag, Link: getTagLink(tag), Count: len(posts)})
	}
	sort.Sort(byCountDesc(tags))
	buf := bytes.Buffer{}
	err = tmpl.Execute(&buf, tags)
	if err != nil {
		return fmt.Errorf("error executing template %s: %v", tagsTemplatePath, err)
	}

	c := htmlConfig{
		path:       tagsPath,
		pageTitle:  "Tags",
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

func (g *tagsGenerator) generateTagPage(tag string, posts []*post) (err error) {
	tagPagePath := filepath.Join(g.siteInfo.DestFolder, "tags", tag)
	err = clearAndCreateDestination(tagPagePath)
	if err != nil {
		return err
	}
	lg := listingGenerator{
		posts:       posts,
		template:    g.template,
		pageTitle:   tag,
		siteInfo:    g.siteInfo,
		destination: tagPagePath,
	}

	err = lg.Generate()
	if err != nil {
		return err
	}
	return nil
}

func (t byCountDesc) Len() int {
	return len(t)
}

func (t byCountDesc) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t byCountDesc) Less(i, j int) bool {
	return t[i].Count > t[j].Count
}
