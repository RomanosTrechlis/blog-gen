package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/russross/blackfriday"
	"github.com/sourcegraph/syntaxhighlight"
	"gopkg.in/yaml.v2"

	"github.com/RomanosTrechlis/blog-generator/config"
	"github.com/RomanosTrechlis/blog-generator/util/fs"
)

// post holds data for a post
type post struct {
	name      string
	html      []byte
	meta      *Meta
	imagesDir string
	images    []string
}

// byDateDesc is the sorting object for posts
type byDateDesc []*post

// postGenerator object
type postGenerator struct {
	config *postConfig
}

// postConfig holds the post's configuration
type postConfig struct {
	post        *post
	siteInfo    *config.SiteInformation
	template    *template.Template
	destination string
}

// Generate generates a post
func (g *postGenerator) Generate() (err error) {
	post := g.config.post
	siteInfo := g.config.siteInfo
	destination := g.config.destination
	t := g.config.template
	fmt.Printf("\tGenerating Post: %s...\n", post.meta.Title)
	staticPath := fmt.Sprintf("%s%s", destination, post.name)
	err = os.Mkdir(staticPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directory at %s: %v", staticPath, err)
	}
	if post.imagesDir != "" {
		err := copyImagesDir(post.imagesDir, staticPath)
		if err != nil {
			return err
		}
	}

	err = writeIndexHTMLPost(staticPath, post.meta.Title, template.HTML(string(post.html)), t, siteInfo, true)
	if err != nil {
		return err
	}

	err = copyAdditionalArtifacts(staticPath, post.name, siteInfo.TempFolder)
	if err != nil {
		return err
	}
	fmt.Printf("\tFinished generating Post: %s...\n", post.meta.Title)
	return nil
}

func newPost(path, dateFormat string) (p *post, err error) {
	meta, err := getMeta(path, dateFormat)
	if err != nil {
		return nil, err
	}
	html, err := getHTML(path)
	if err != nil {
		return nil, err
	}
	imagesDir, images, err := getImages(path)
	if err != nil {
		return nil, err
	}
	name := path[strings.LastIndex(path, "/"):]
	p = &post{name: name, meta: meta, html: html, imagesDir: imagesDir, images: images}
	return p, nil
}

func copyDir(source, path string) (err error) {
	files, err := ioutil.ReadDir(source)
	if err != nil {
		return nil
	}
	for _, file := range files {
		src := fmt.Sprintf("%s/%s", source, file.Name())
		dst := fmt.Sprintf("%s/%s", path, file.Name())
		err := fs.CopyFile(src, dst)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyImagesDir(source, destination string) (err error) {
	path := fmt.Sprintf("%s/images", destination)
	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating images directory at %s: %v", path, err)
	}
	files, err := ioutil.ReadDir(source)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %v", path, err)
	}
	for _, file := range files {
		src := fmt.Sprintf("%s/%s", source, file.Name())
		dst := fmt.Sprintf("%s/%s", path, file.Name())
		err := fs.CopyFile(src, dst)
		if err != nil {
			return err
		}
	}
	return nil
}

func getMeta(path, dateFormat string) (metaP *Meta, err error) {
	filePath := fmt.Sprintf("%s/meta.yml", path)
	metaraw, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error while reading file %s: %v", filePath, err)
	}
	meta := Meta{}
	err = yaml.Unmarshal(metaraw, &meta)
	if err != nil {
		return nil, fmt.Errorf("error reading yml in %s: %v", filePath, err)
	}
	parsedDate, err := time.Parse(dateFormat, meta.Date)
	if err != nil {
		return nil, fmt.Errorf("error parsing date in %s: %v", filePath, err)
	}
	meta.ParsedDate = parsedDate
	return &meta, nil
}

func getHTML(path string) (html []byte, err error) {
	filePath := fmt.Sprintf("%s/post.md", path)
	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error while reading file %s: %v", filePath, err)
	}
	html = blackfriday.MarkdownCommon(input)
	replaced, err := replaceCodeParts(html)
	if err != nil {
		return nil, fmt.Errorf("error during syntax highlighting of %s: %v", filePath, err)
	}
	html = []byte(replaced)
	return html, nil

}

func getImages(path string) (dirPath string, images []string, err error) {
	dirPath = fmt.Sprintf("%s/images", path)
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil, nil
		}
		return "", nil, fmt.Errorf("error while reading folder %s: %v", dirPath, err)
	}
	images = []string{}
	for _, file := range files {
		images = append(images, file.Name())
	}
	return dirPath, images, nil
}

func replaceCodeParts(htmlFile []byte) (new string, err error) {
	byteReader := bytes.NewReader(htmlFile)
	doc, err := goquery.NewDocumentFromReader(byteReader)
	if err != nil {
		return "", fmt.Errorf("error while parsing html: %v", err)
	}
	// find code-parts via css selector and replace them with highlighted versions
	doc.Find("code[class*=\"language-\"]").Each(func(i int, s *goquery.Selection) {
		oldCode := s.Text()
		formatted, _ := syntaxhighlight.AsHTML([]byte(oldCode))
		s.SetHtml(string(formatted))
	})
	new, err = doc.Html()
	if err != nil {
		return "", fmt.Errorf("error while generating html: %v", err)
	}
	// replace unnecessarily added html tags
	new = strings.Replace(new, "<html><head></head><body>", "", 1)
	new = strings.Replace(new, "</body></html>", "", 1)
	return new, nil
}

func (p byDateDesc) Len() int {
	return len(p)
}

func (p byDateDesc) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p byDateDesc) Less(i, j int) bool {
	return p[i].meta.ParsedDate.After(p[j].meta.ParsedDate)
}
