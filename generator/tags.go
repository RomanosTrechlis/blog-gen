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

// ByCountDesc sorts the tags
type ByCountDesc []*Tag

// TagsGenerator object
type TagsGenerator struct {
	Config *TagsConfig
}

// TagsConfig holds the tag's config
type TagsConfig struct {
	TagPostsMap     map[string][]*Post
	Template        *template.Template
	SiteInfo        *config.SiteInformation
}

var tagsTemplatePath string

// Generate creates the tags page
func (g *TagsGenerator) Generate() (err error) {
	fmt.Println("\tGenerating Tags...")
	siteInfo := g.Config.SiteInfo

	tagsTemplatePath = siteInfo.ThemeFolder + "tags.html"
	tagPostsMap := g.Config.TagPostsMap
	t := g.Config.Template
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

func generateTagIndex(tagPostsMap map[string][]*Post, t *template.Template, tagsPath string, siteInfo *config.SiteInformation) (err error) {
	tmpl, err := getTemplate(tagsTemplatePath)
	if err != nil {
		return err
	}
	tags := []*Tag{}
	for tag, posts := range tagPostsMap {
		tags = append(tags, &Tag{Name: tag, Link: getTagLink(tag), Count: len(posts)})
	}
	sort.Sort(ByCountDesc(tags))
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

func generateTagPage(destination, tag string, posts []*Post, t *template.Template, siteInfo *config.SiteInformation) (err error) {
	err = clearAndCreateDestination(destination)
	if err != nil {
		return err
	}
	lg := ListingGenerator{&ListingConfig{
		Posts:       posts,
		Template:    t,
		PageTitle:   tag,
		SiteInfo:    siteInfo,
		Destination: destination,
	}}

	err = lg.Generate()
	if err != nil {
		return err
	}
	return nil
}

func (t ByCountDesc) Len() int {
	return len(t)
}

func (t ByCountDesc) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t ByCountDesc) Less(i, j int) bool {
	return t[i].Count > t[j].Count
}
