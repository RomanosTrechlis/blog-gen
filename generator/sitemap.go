package generator

import (
	"fmt"
	"os"

	"github.com/beevik/etree"
)

// sitemapGenerator object
type sitemapGenerator struct {
	config *sitemapConfig
}

// sitemapConfig holds the config for the sitemap
type sitemapConfig struct {
	posts            []*post
	tagPostsMap      map[string][]*post
	categoryPostsMap map[string][]*post
	destination      string
	blogURL          string
}

// Generate creates the sitemap
func (g *sitemapGenerator) Generate() (err error) {
	fmt.Println("\tGenerating Sitemap...")
	posts := g.config.posts
	tagPostsMap := g.config.tagPostsMap
	catPostsMap := g.config.categoryPostsMap
	destination := g.config.destination
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	urlSet := doc.CreateElement("urlset")
	urlSet.CreateAttr("xmlns", "http://www.sitemaps.org/schemas/sitemap/0.9")
	urlSet.CreateAttr("xmlns:image", "http://www.google.com/schemas/sitemap-image/1.1")

	url := urlSet.CreateElement("url")
	loc := url.CreateElement("loc")
	loc.SetText(g.config.blogURL)

	addURL(urlSet, "about", g.config.blogURL, nil)
	addURL(urlSet, "archive", g.config.blogURL, nil)
	addURL(urlSet, "tags", g.config.blogURL, nil)
	addURL(urlSet, "categories", g.config.blogURL, nil)

	for tag := range tagPostsMap {
		addURL(urlSet, tag, g.config.blogURL, nil)
	}

	for cat := range catPostsMap {
		addURL(urlSet, cat, g.config.blogURL, nil)
	}

	for _, post := range posts {
		addURL(urlSet, post.name[1:], g.config.blogURL, post.images)
	}

	filePath := fmt.Sprintf("%s/sitemap.xml", destination)
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", filePath, err)
	}
	f.Close()
	err = doc.WriteToFile(filePath)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %v", filePath, err)
	}
	fmt.Println("\tFinished generating Sitemap...")
	return nil
}

func addURL(element *etree.Element, location, blogUrl string, images []string) {
	url := element.CreateElement("url")
	loc := url.CreateElement("loc")
	loc.SetText(fmt.Sprintf("%s/%s/", blogUrl, location))

	if len(images) > 0 {
		for _, image := range images {
			img := url.CreateElement("image:image")
			imgLoc := img.CreateElement("image:loc")
			imgLoc.SetText(fmt.Sprintf("%s/%s/images/%s", blogUrl, location, image))
		}
	}
}
