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
	config *tagsConfig
}

// tagsConfig holds the tag's config
type tagsConfig struct {
	tagPostsMap map[string][]*post
	template    *template.Template
	siteInfo    *config.SiteInformation
}

var tagsTemplatePath string

// Generate creates the tags page
func (g *tagsGenerator) Generate() (err error) {
	fmt.Println("\tGenerating Tags...")
	siteInfo := g.config.siteInfo

	tagsTemplatePath = siteInfo.ThemeFolder + "tags.html"
	tagPostsMap := g.config.tagPostsMap
	t := g.config.template
	destination := siteInfo.DestFolder
	tagsPath := fmt.Sprintf("%s/tags", destination)
	err = clearAndCreateDestination(tagsPath)
	if err != nil {
		return err
	}
	err = generateTagIndex(tagPostsMap, t, tagsPath, siteInfo)
	if err != nil {
		return err
	}
	for tag, tagPosts := range tagPostsMap {
		tagPagePath := fmt.Sprintf("%s/%s", tagsPath, tag)
		err := generateTagPage(tagPagePath, tag, tagPosts, t, siteInfo)
		if err != nil {
			return err
		}
	}
	fmt.Println("\tFinished generating Tags...")
	return nil
}

func generateTagIndex(tagPostsMap map[string][]*post, t *template.Template, tagsPath string, siteInfo *config.SiteInformation) (err error) {
	tmpl, err := getTemplate(tagsTemplatePath)
	if err != nil {
		return err
	}
	tags := []*Tag{}
	for tag, posts := range tagPostsMap {
		tags = append(tags, &Tag{Name: tag, Link: getTagLink(tag), Count: len(posts)})
	}
	sort.Sort(byCountDesc(tags))
	buf := bytes.Buffer{}
	err = tmpl.Execute(&buf, tags)
	if err != nil {
		return fmt.Errorf("error executing template %s: %v", tagsTemplatePath, err)
	}
	err = writeIndexHTML(tagsPath, "Tags", template.HTML(buf.String()), t, siteInfo)
	if err != nil {
		return err
	}
	return nil
}

func generateTagPage(destination, tag string, posts []*post, t *template.Template, siteInfo *config.SiteInformation) (err error) {
	err = clearAndCreateDestination(destination)
	if err != nil {
		return err
	}
	lg := listingGenerator{&listingConfig{
		posts:       posts,
		template:    t,
		pageTitle:   tag,
		siteInfo:    siteInfo,
		destination: destination,
	}}

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
