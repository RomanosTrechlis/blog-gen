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
	Tags       []Tag
}

// listingGenerator Object
type listingGenerator struct {
	config *listingConfig
}

// listingConfig holds the configuration for the listing page
type listingConfig struct {
	posts                  []*post
	template               *template.Template
	siteInfo               *config.SiteInformation
	destination, pageTitle string
	pageNum                int
	maxPageNum             int
}

var shortTemplatePath string

// Generate starts the listing generation
func (g *listingGenerator) Generate() (err error) {
	shortTemplatePath = g.config.siteInfo.ThemeFolder + "short.html"
	posts := g.config.posts
	t := g.config.template
	destination := g.config.destination
	pageTitle := g.config.pageTitle
	short, err := getTemplate(shortTemplatePath)
	if err != nil {
		return err
	}
	var postBlocks []string
	for _, post := range posts {
		meta := post.meta
		link := fmt.Sprintf("%s/", post.name)
		ld := ListingData{
			Title:      meta.Title,
			Date:       meta.Date,
			Short:      meta.Short,
			Link:       link,
			Tags:       createTags(meta.Tags),
			TimeToRead: calculateTimeToRead(string(post.html)),
		}
		block := bytes.Buffer{}
		err := short.Execute(&block, ld)
		if err != nil {
			return fmt.Errorf("error executing template %s: %v", shortTemplatePath, err)
		}
		postBlocks = append(postBlocks, block.String())
	}
	htmlBlocks := template.HTML(strings.Join(postBlocks, "<br />"))
	if g.config.pageNum > 1 {
		err := os.Mkdir(destination, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory at %s: %v", destination, err)
		}
	}
	err = writeIndexHTMLPlus(destination, pageTitle, htmlBlocks, t, g.config.siteInfo, false, g.config.pageNum, g.config.maxPageNum)
	if err != nil {
		return err
	}
	return nil
}

func createTags(tags []string) (result []Tag) {
	for _, tag := range tags {
		result = append(result, Tag{Name: tag, Link: getTagLink(tag)})
	}
	return result
}

func calculateTimeToRead(input string) (time string) {
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
	time = fmt.Sprintf("%.0fm", result)
	return time
}
