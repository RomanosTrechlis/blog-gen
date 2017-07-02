package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/RomanosTrechlis/blog-generator/config"
)

// ListingData holds the data for the listing page
type ListingData struct {
	Title      string
	Date       string
	Short      string
	Link       string
	TimeToRead string
	Tags       []*Tag
}

// ListingGenerator Object
type ListingGenerator struct {
	Config *ListingConfig
}

// ListingConfig holds the configuration for the listing page
type ListingConfig struct {
	Posts                  []*Post
	Template               *template.Template
	Destination, PageTitle string
	PageNum                int
	MaxPageNum             int
}

var shortTemplatePath string

// Generate starts the listing generation
func (g *ListingGenerator) Generate() error {
	shortTemplatePath = config.SiteInfo.ThemePath + "short.html"
	posts := g.Config.Posts
	t := g.Config.Template
	destination := g.Config.Destination
	pageTitle := g.Config.PageTitle
	short, err := getTemplate(shortTemplatePath)
	if err != nil {
		return err
	}
	var postBlocks []string
	for _, post := range posts {
		meta := post.Meta
		link := fmt.Sprintf("%s/", post.Name)
		ld := ListingData{
			Title:      meta.Title,
			Date:       meta.Date,
			Short:      meta.Short,
			Link:       link,
			Tags:       createTags(meta.Tags),
			TimeToRead: calculateTimeToRead(string(post.HTML)),
		}
		block := bytes.Buffer{}
		err := short.Execute(&block, ld)
		if err != nil {
			return fmt.Errorf("error executing template %s: %v", shortTemplatePath, err)
		}
		postBlocks = append(postBlocks, block.String())
	}
	htmlBlocks := template.HTML(strings.Join(postBlocks, "<br />"))
	if g.Config.PageNum > 1 {
		err := os.Mkdir(destination, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory at %s: %v", destination, err)
		}
	}
	err = writeIndexHTMLPlus(destination, pageTitle, htmlBlocks, t, false, g.Config.PageNum, g.Config.MaxPageNum)
	if err != nil {
		return err
	}
	return nil
}

func createTags(tags []string) []*Tag {
	var result []*Tag
	for _, tag := range tags {
		result = append(result, &Tag{Name: tag, Link: getTagLink(tag)})
	}
	return result
}

func calculateTimeToRead(input string) string {
	// an average human reads about 200 wpm
	var secondsPerWord = 60.0 / 200.0
	// multiply with the amount of words
	words := secondsPerWord * float64(len(strings.Split(input, " ")))
	// add 12 seconds for each image
	images := 12.0 * strings.Count(input, "<img")
	result := (words + float64(images)) / 60.0
	if result < 1.0 {
		result = 1.0
	}
	return fmt.Sprintf("%.0fm", result)
}
