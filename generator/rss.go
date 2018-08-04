package generator

import (
	"fmt"
	"os"
	"time"

	"github.com/RomanosTrechlis/blog-generator/config"
	"github.com/beevik/etree"
)

// rssGenerator object
type rssGenerator struct {
	posts       []*post
	destination string
	siteInfo    *config.SiteInformation
}

const rssDateFormat = "02 Jan 2006 15:04 -0700"

// Generate creates an RSS feed
func (g *rssGenerator) Generate() (err error) {
	fmt.Println("\tGenerating RSS...")
	posts := g.posts
	destination := g.destination
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	rss := doc.CreateElement("rss")
	rss.CreateAttr("xmlns:atom", "http://www.w3.org/2005/Atom")
	rss.CreateAttr("version", "2.0")
	channel := rss.CreateElement("channel")
	siteInfo := g.siteInfo

	channel.CreateElement("title").SetText(siteInfo.BlogTitle)
	channel.CreateElement("link").SetText(siteInfo.BlogURL)
	channel.CreateElement("language").SetText(siteInfo.BlogLanguage)
	channel.CreateElement("description").SetText(siteInfo.BlogDescription)
	channel.CreateElement("lastBuildDate").SetText(time.Now().Format(rssDateFormat))

	atomLink := channel.CreateElement("atom:link")
	atomLink.CreateAttr("href", fmt.Sprintf("%s/index.xml", siteInfo.BlogURL))
	atomLink.CreateAttr("rel", "self")
	atomLink.CreateAttr("type", "application/rss+xml")

	for _, post := range posts {
		err := g.addItem(channel, post)
		if err != nil {
			return err
		}
	}

	filePath := fmt.Sprintf("%s/index.xml", destination)
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", filePath, err)
	}
	f.Close()
	err = doc.WriteToFile(filePath)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %v", filePath, err)
	}
	fmt.Println("\tFinished generating RSS...")
	return nil
}

func (g *rssGenerator) addItem(element *etree.Element, post *post) (err error) {
	path := fmt.Sprintf("%s/%s/", g.siteInfo.BlogURL, post.name[1:])
	meta := post.meta
	item := element.CreateElement("item")
	item.CreateElement("title").SetText(meta.Title)
	item.CreateElement("link").SetText(path)
	item.CreateElement("guid").SetText(path)
	pubDate, err := time.Parse(g.siteInfo.DateFormat, meta.Date)
	if err != nil {
		return fmt.Errorf("error parsing date %s: %v", meta.Date, err)
	}
	item.CreateElement("pubDate").SetText(pubDate.Format(rssDateFormat))
	item.CreateElement("description").SetText(string(post.html))
	return nil
}
