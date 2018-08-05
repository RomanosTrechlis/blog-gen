package generator

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/RomanosTrechlis/blog-gen/config"
	"github.com/RomanosTrechlis/blog-gen/util/fs"
	"github.com/RomanosTrechlis/blog-gen/util/url"
	"gopkg.in/yaml.v2"
)

// siteGenerator object
type siteGenerator struct {
	Sources  []string
	SiteInfo *config.SiteInformation
}

// New creates a new SiteGenerator
func NewSiteGenerator(sources []string, siteInfo *config.SiteInformation) *siteGenerator {
	return &siteGenerator{sources, siteInfo}
}

var templatePath string

// Generate starts the static blog generation
func (g *siteGenerator) Generate() (err error) {
	templatePath = filepath.Join(g.SiteInfo.ThemeFolder, "template.html")
	fmt.Println("Generating Site...")
	err = clearAndCreateDestination(g.SiteInfo.DestFolder)
	if err != nil {
		return err
	}

	err = clearAndCreateDestination(filepath.Join(g.SiteInfo.DestFolder, "archive"))
	if err != nil {
		return err
	}

	t, err := getTemplate(templatePath)
	if err != nil {
		return err
	}

	posts := make([]*post, 0)
	for _, path := range g.Sources {
		post, err := g.newPost(path)
		if err != nil {
			return err
		}
		posts = append(posts, post)
	}
	sort.Sort(byDateDesc(posts))

	generators := g.createTasks(posts, t)
	err = g.runTasks(generators)
	if err != nil {
		return err
	}
	fmt.Println("Finished generating Site...")
	return nil
}

func (g *siteGenerator) newPost(path string) (p *post, err error) {
	meta, err := g.getPostMeta(path)
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
	name := path[strings.LastIndex(path, fs.GetSeparator())+1:]
	p = &post{name: name, meta: meta, html: html, imagesDir: imagesDir, images: images}
	return p, nil
}

func (g *siteGenerator) getPostMeta(path string) (*Meta, error) {
	filePath := filepath.Join(path, "meta.yml")
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error while reading file %s: %v", filePath, err)
	}
	meta := Meta{}
	err = yaml.Unmarshal(b, &meta)
	if err != nil {
		return nil, fmt.Errorf("error reading yml in %s: %v", filePath, err)
	}
	parsedDate, err := time.Parse(g.SiteInfo.DateFormat, meta.Date)
	if err != nil {
		return nil, fmt.Errorf("error parsing date in %s: %v", filePath, err)
	}
	meta.ParsedDate = parsedDate
	return &meta, nil
}

func (g *siteGenerator) createTasks(posts []*post, t *template.Template) []Generator {
	generators := make([]Generator, 0)
	destination := g.SiteInfo.DestFolder

	//posts
	for _, post := range posts {
		pg := postGenerator{post, g.SiteInfo, t, destination}
		generators = append(generators, &pg)
	}
	tagPostsMap := createTagPostsMap(posts)

	// frontpage
	paging := g.SiteInfo.NumPostsFrontPage
	numOfPages := getNumberOfPages(posts, paging)
	for i := 0; i < numOfPages; i++ {
		to := destination
		if i != 0 {
			to = filepath.Join(destination, fmt.Sprintf("%d", i+1))
		}
		toP := (i + 1) * paging
		if (i + 1) == numOfPages {
			toP = len(posts)
		}
		lg := &listingGenerator{posts[i*paging : toP], t, g.SiteInfo, to, "", i + 1, numOfPages}
		generators = append(generators, lg)
	}

	// archive
	ag := listingGenerator{posts, t, g.SiteInfo, filepath.Join(destination, "archive"), "Archive", 0, 0}
	// tags
	tg := tagsGenerator{
		tagPostsMap: tagPostsMap,
		template:    t,
		siteInfo:    g.SiteInfo,
	}
	// categories
	catPostsMap := createCatPostsMap(posts)
	ct := categoriesGenerator{
		catPostsMap: catPostsMap,
		template:    t,
		destination: destination,
		siteInfo:    g.SiteInfo,
	}
	// sitemap
	sg := sitemapGenerator{
		posts:            posts,
		tagPostsMap:      tagPostsMap,
		categoryPostsMap: catPostsMap,
		destination:      destination,
		blogURL:          g.SiteInfo.BlogURL,
	}
	// rss
	rg := rssGenerator{
		posts:       posts,
		destination: destination,
		siteInfo:    g.SiteInfo,
	}
	// statics
	fileToDestination := make(map[string]string)
	templateToFile := make(map[string]string)
	for _, row := range g.SiteInfo.StaticPages {
		if row.IsTemplate {
			templateToFile[filepath.Join(g.SiteInfo.ThemeFolder, row.File)] = filepath.Join(destination, row.To)
			continue
		}
		fileToDestination[filepath.Join(g.SiteInfo.ThemeFolder, row.File)] = filepath.Join(destination, row.To)
	}
	statg := staticsGenerator{
		fileToDestination: fileToDestination,
		templateToFile:    templateToFile,
		template:          t,
		siteInfo:          g.SiteInfo,
	}
	generators = append(generators, &ag, &tg, &ct, &sg, &rg, &statg)
	return generators
}

func (g *siteGenerator) runTasks(generators []Generator) (err error) {
	var wg sync.WaitGroup
	finished := make(chan bool, 1)
	errors := make(chan error, 1)
	pool := make(chan struct{}, 50)

	for _, generator := range generators {
		wg.Add(1)
		go func(g Generator) {
			defer wg.Done()
			pool <- struct{}{}
			defer func() { <-pool }()

			err := g.Generate()
			if err != nil {
				errors <- err
			}
		}(generator)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
		return nil
	case err := <-errors:
		if err != nil {
			return err
		}
	}
	return nil
}

type htmlConfig struct {
	path       string
	pageTitle  string
	pageNum    int
	maxPageNum int
	isPost     bool
	temp       *template.Template
	content    template.HTML
	siteInfo   *config.SiteInformation
}

func (h htmlConfig) writeHTML() error {
	filePath := filepath.Join(h.path, "index.html")
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", filePath, err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	next := h.pageNum + 1
	prev := h.pageNum - 1
	if h.pageNum == h.maxPageNum {
		next = 0
	}

	u := url.ChangePathToUrl(h.path)
	td := IndexData{
		Name:          h.siteInfo.Author,
		Year:          time.Now().Year(),
		HTMLTitle:     getHTMLTitle(h.pageTitle, h.siteInfo.BlogTitle),
		PageTitle:     h.pageTitle,
		Content:       h.content,
		CanonicalLink: buildCanonicalLink(u, h.siteInfo.BlogURL),
		PageNum:       h.pageNum,
		NextPageNum:   next,
		PrevPageNum:   prev,
		URL:           buildCanonicalLink(u, h.siteInfo.BlogURL),
		IsPost:        h.isPost,
	}

	err = h.temp.Execute(w, td)
	if err != nil {
		return fmt.Errorf("error executing template %s: %v", templatePath, err)
	}
	err = w.Flush()
	if err != nil {
		return fmt.Errorf("error writing file %s: %v", filePath, err)
	}
	return nil
}

func getHTMLTitle(pageTitle, blogTitle string) (title string) {
	if pageTitle == "" {
		return blogTitle
	}
	return fmt.Sprintf("%s - %s", pageTitle, blogTitle)
}

func createTagPostsMap(posts []*post) (result map[string][]*post) {
	result = make(map[string][]*post)
	for _, p := range posts {
		for _, tag := range p.meta.Tags {
			key := strings.ToLower(tag)
			if result[key] == nil {
				result[key] = []*post{p}
			} else {
				result[key] = append(result[key], p)
			}
		}
	}
	return result
}

func createCatPostsMap(posts []*post) (result map[string][]*post) {
	result = make(map[string][]*post)
	for _, p := range posts {
		for _, cat := range p.meta.Categories {
			key := strings.ToLower(cat)
			if result[key] == nil {
				result[key] = []*post{p}
			} else {
				result[key] = append(result[key], p)
			}
		}
	}
	return result
}

func getNumberOfPages(posts []*post, postsPerPage int) (n int) {
	res := float64(len(posts)) / float64(postsPerPage)
	n, _ = strconv.Atoi(fmt.Sprintf("%.0f", math.Ceil(res)))
	return n
}
