package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"

	"github.com/RomanosTrechlis/blog-generator/config"
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
	tagsPath := fmt.Sprintf("%s/tags", siteInfo.DestFolder)
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
	tagsPath := fmt.Sprintf("%s/tags", g.siteInfo.DestFolder)
	tagsTemplatePath := g.siteInfo.ThemeFolder + "tags.html"
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
	err = writeIndexHTML(tagsPath, "Tags", template.HTML(buf.String()), g.template, g.siteInfo)
	if err != nil {
		return err
	}
	return nil
}

func (g *tagsGenerator) generateTagPage(tag string, posts []*post) (err error) {
	tagPagePath := fmt.Sprintf("%s/%s/%s", g.siteInfo.DestFolder, "tags", tag)
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
