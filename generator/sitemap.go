package generator

import (
	"fmt"
	"github.com/beevik/etree"
	"os"

	"github.com/RomanosTrechlis/blog-generator/config"
)

// SitemapGenerator object
type SitemapGenerator struct {
	Config *SitemapConfig
}

// SitemapConfig holds the config for the sitemap
type SitemapConfig struct {
	Posts            []*Post
	TagPostsMap      map[string][]*Post
	CategoryPostsMap map[string][]*Post
	Destination      string
}

// Generate creates the sitemap
func (g *SitemapGenerator) Generate() error {
	fmt.Println("\tGenerating Sitemap...")
	posts := g.Config.Posts
	tagPostsMap := g.Config.TagPostsMap
	catPostsMap := g.Config.CategoryPostsMap
	destination := g.Config.Destination
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	urlSet := doc.CreateElement("urlset")
	urlSet.CreateAttr("xmlns", "http://www.sitemaps.org/schemas/sitemap/0.9")
	urlSet.CreateAttr("xmlns:image", "http://www.google.com/schemas/sitemap-image/1.1")

	url := urlSet.CreateElement("url")
	loc := url.CreateElement("loc")
	loc.SetText(config.SiteInfo.BlogURL)

	addURL(urlSet, "about", nil)
	addURL(urlSet, "archive", nil)
	addURL(urlSet, "tags", nil)
	addURL(urlSet, "categories", nil)

	for tag := range tagPostsMap {
		addURL(urlSet, tag, nil)
	}

	for cat := range catPostsMap {
		addURL(urlSet, cat, nil)
	}

	for _, post := range posts {
		addURL(urlSet, post.Name[1:], post.Images)
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

func addURL(element *etree.Element, location string, images []string) {
	url := element.CreateElement("url")
	loc := url.CreateElement("loc")
	loc.SetText(fmt.Sprintf("%s/%s/", config.SiteInfo.BlogURL, location))

	if len(images) > 0 {
		for _, image := range images {
			img := url.CreateElement("image:image")
			imgLoc := img.CreateElement("image:loc")
			imgLoc.SetText(fmt.Sprintf("%s/%s/images/%s", config.SiteInfo.BlogURL, location, image))
		}
	}
}
