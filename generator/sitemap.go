package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/beevik/etree"
)

// sitemapGenerator object
type sitemapGenerator struct {
	posts            []*post
	tagPostsMap      map[string][]*post
	categoryPostsMap map[string][]*post
	destination      string
	blogURL          string
}

// Generate creates the sitemap
func (g *sitemapGenerator) Generate() (err error) {
	fmt.Println("\tGenerating Sitemap...")
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	urlSet := doc.CreateElement("urlset")
	urlSet.CreateAttr("xmlns", "http://www.sitemaps.org/schemas/sitemap/0.9")
	urlSet.CreateAttr("xmlns:image", "http://www.google.com/schemas/sitemap-image/1.1")

	url := urlSet.CreateElement("url")
	loc := url.CreateElement("loc")
	loc.SetText(g.blogURL)

	g.addURL(urlSet, "about", nil)
	g.addURL(urlSet, "archive", nil)
	g.addURL(urlSet, "tags", nil)
	g.addURL(urlSet, "categories", nil)

	for tag := range g.tagPostsMap {
		g.addURL(urlSet, tag, nil)
	}

	for cat := range g.categoryPostsMap {
		g.addURL(urlSet, cat, nil)
	}

	for _, post := range g.posts {
		g.addURL(urlSet, post.name[1:], post.images)
	}

	filePath := filepath.Join(g.destination, "sitemap.xml")
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

func (g *sitemapGenerator) addURL(element *etree.Element, location string, images []string) {
	url := element.CreateElement("url")
	loc := url.CreateElement("loc")
	loc.SetText(fmt.Sprintf("%s/%s/", g.blogURL, location))

	if len(images) > 0 {
		for _, image := range images {
			img := url.CreateElement("image:image")
			imgLoc := img.CreateElement("image:loc")
			imgLoc.SetText(fmt.Sprintf("%s/%s/images/%s", g.blogURL, location, image))
		}
	}
}
