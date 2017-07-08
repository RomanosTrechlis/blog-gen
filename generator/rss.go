package generator

import (
	"fmt"
	"github.com/beevik/etree"
	"os"
	"time"
	"github.com/RomanosTrechlis/blog-generator/config"
)

// RSSGenerator object
type RSSGenerator struct {
	Config *RSSConfig
}

// RSSConfig holds the configuration for an RSS feed
type RSSConfig struct {
	Posts           []*Post
	Destination     string
	SiteInfo				*config.SiteInformation
}

const rssDateFormat string = "02 Jan 2006 15:04 -0700"

// Generate creates an RSS feed
func (g *RSSGenerator) Generate() (err error) {
	fmt.Println("\tGenerating RSS...")
	posts := g.Config.Posts
	destination := g.Config.Destination
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	rss := doc.CreateElement("rss")
	rss.CreateAttr("xmlns:atom", "http://www.w3.org/2005/Atom")
	rss.CreateAttr("version", "2.0")
	channel := rss.CreateElement("channel")
	siteInfo := g.Config.SiteInfo

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
		err := addItem(channel, post, fmt.Sprintf("%s/%s/", siteInfo.BlogURL, post.Name[1:]), siteInfo.DateFormat)
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

func addItem(element *etree.Element, post *Post, path, dateFormat string) (err error) {
	meta := post.Meta
	item := element.CreateElement("item")
	item.CreateElement("title").SetText(meta.Title)
	item.CreateElement("link").SetText(path)
	item.CreateElement("guid").SetText(path)
	pubDate, err := time.Parse(dateFormat, meta.Date)
	if err != nil {
		return fmt.Errorf("error parsing date %s: %v", meta.Date, err)
	}
	item.CreateElement("pubDate").SetText(pubDate.Format(rssDateFormat))
	item.CreateElement("description").SetText(string(post.HTML))
	return nil
}
