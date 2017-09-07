package generator

import (
	"bufio"
	"fmt"
	"html/template"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/RomanosTrechlis/blog-generator/config"
)

// siteGenerator object
type siteGenerator struct {
	config *SiteConfig
}

// SiteConfig holds the sources and destination folder
type SiteConfig struct {
	Sources  []string
	SiteInfo *config.SiteInformation
}

// New creates a new SiteGenerator
func NewSiteGenerator(config *SiteConfig) *siteGenerator {
	return &siteGenerator{config: config}
}

var templatePath string

// Generate starts the static blog generation
func (g *siteGenerator) Generate() (err error) {
	templatePath = g.config.SiteInfo.ThemeFolder + "template.html"
	fmt.Println("Generating Site...")
	sources := g.config.Sources
	destination := g.config.SiteInfo.DestFolder
	err = clearAndCreateDestination(destination)
	if err != nil {
		return err
	}
	err = clearAndCreateDestination(fmt.Sprintf("%s/archive", destination))
	if err != nil {
		return err
	}
	t, err := getTemplate(templatePath)
	if err != nil {
		return err
	}
	var posts []*post
	for _, path := range sources {
		post, err := newPost(path, g.config.SiteInfo.DateFormat)
		if err != nil {
			return err
		}
		posts = append(posts, post)
	}
	sort.Sort(byDateDesc(posts))
	err = runTasks(posts, t, g.config.SiteInfo)
	if err != nil {
		return err
	}
	fmt.Println("Finished generating Site...")
	return nil
}

func runTasks(posts []*post, t *template.Template, siteInfo *config.SiteInformation) (err error) {
	var wg sync.WaitGroup
	finished := make(chan bool, 1)
	errors := make(chan error, 1)
	pool := make(chan struct{}, 50)
	generators := []Generator{}
	destination := siteInfo.DestFolder

	//posts
	for _, post := range posts {
		pg := postGenerator{newPostConfig(post, destination, t, siteInfo)}
		generators = append(generators, &pg)
	}
	tagPostsMap := createTagPostsMap(posts)

	// frontpage
	paging := siteInfo.NumPostsFrontPage
	numOfPages := getNumberOfPages(posts, paging)
	for i := 0; i < numOfPages; i++ {
		to := destination
		if i != 0 {
			to = fmt.Sprintf("%s/%d", destination, i+1)
		}
		toP := (i + 1) * paging
		if (i + 1) == numOfPages {
			toP = len(posts)
		}
		generators = append(generators, &listingGenerator{newListingConfig(posts[i*paging : toP], t, siteInfo, to, "", i + 1, numOfPages)})
	}

	// archive
	ag := listingGenerator{newListingConfig(posts, t, siteInfo, fmt.Sprintf("%s/archive", destination), "Archive", 0, 0)}
	// tags
	tg := tagsGenerator{&tagsConfig{
		tagPostsMap: tagPostsMap,
		template:    t,
		siteInfo:    siteInfo,
	}}
	// categories
	catPostsMap := createCatPostsMap(posts)
	ct := categoriesGenerator{&categoriesConfig{
		catPostsMap: catPostsMap,
		template:    t,
		destination: destination,
		siteInfo:    siteInfo,
	}}

	// sitemap
	sg := sitemapGenerator{&sitemapConfig{
		posts:            posts,
		tagPostsMap:      tagPostsMap,
		categoryPostsMap: catPostsMap,
		destination:      destination,
		blogURL:          siteInfo.BlogURL,
	}}
	// rss
	rg := rssGenerator{&rssConfig{
		posts:       posts,
		destination: destination,
		siteInfo:    siteInfo,
	}}
	// statics
	fileToDestination := make(map[string]string)
	templateToFile := make(map[string]string)
	for _, row := range siteInfo.StaticPages {
		if row.IsTemplate {
			templateToFile[siteInfo.ThemeFolder+row.File] = fmt.Sprintf("%s/%s", destination, row.To)
			continue
		}
		fileToDestination[siteInfo.ThemeFolder+row.File] = fmt.Sprintf("%s/%s", destination, row.To)
	}
	statg := staticsGenerator{&staticsConfig{
		fileToDestination: fileToDestination,
		templateToFile:    templateToFile,
		template:          t,
		siteInfo:          siteInfo,
	}}
	generators = append(generators, &ag, &tg, &ct, &sg, &rg, &statg)

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

func writeIndexHTML(path, pageTitle string, content template.HTML, t *template.Template, siteInfo *config.SiteInformation) (err error) {
	return writeIndexHTMLPlus(path, pageTitle, content, t, siteInfo, false, 0, 0)
}

func writeIndexHTMLPost(path, pageTitle string, content template.HTML, t *template.Template, siteInfo *config.SiteInformation,
	isPost bool) (err error) {
	return writeIndexHTMLPlus(path, pageTitle, content, t, siteInfo, isPost, 0, 0)
}

func writeIndexHTMLPlus(path, pageTitle string, content template.HTML, t *template.Template, siteInfo *config.SiteInformation,
	isPost bool, page, maxPage int) (err error) {
	filePath := fmt.Sprintf("%s/index.html", path)
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", filePath, err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	next := page + 1
	prev := page - 1
	if page == maxPage {
		next = 0
	}
	td := IndexData{
		Name:          siteInfo.Author,
		Year:          time.Now().Year(),
		HTMLTitle:     getHTMLTitle(pageTitle, siteInfo.BlogTitle),
		PageTitle:     pageTitle,
		Content:       content,
		CanonicalLink: buildCanonicalLink(path, siteInfo.BlogURL),
		PageNum:       page,
		NextPageNum:   next,
		PrevPageNum:   prev,
		URL:           buildCanonicalLink(path, siteInfo.BlogURL),
		IsPost:        isPost,
	}

	err = t.Execute(w, td)
	if err != nil {
		return fmt.Errorf("error executing template %s: %v", templatePath, err)
	}
	err = w.Flush()
	if err != nil {
		return fmt.Errorf("error writing file %s: %v", filePath, err)
	}
	return nil
}

func copyAdditionalArtifacts(path, postName, tempFolder string) (err error) {
	src := tempFolder + postName + "/artifacts/"
	return copyDir(src, path)
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
