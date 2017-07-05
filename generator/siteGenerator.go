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

// SiteGenerator object
type SiteGenerator struct {
	Config *SiteConfig
}

// SiteConfig holds the sources and destination folder
type SiteConfig struct {
	Sources     []string
	Destination string
}

// New creates a new SiteGenerator
func NewSiteGenerator(config *SiteConfig) *SiteGenerator {
	return &SiteGenerator{Config: config}
}

var templatePath string

// Generate starts the static blog generation
func (g *SiteGenerator) Generate() error {
	templatePath = config.SiteInfo.ThemeFolder + "template.html"
	fmt.Println("Generating Site...")
	sources := g.Config.Sources
	destination := g.Config.Destination
	err := clearAndCreateDestination(destination)
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
	var posts []*Post
	for _, path := range sources {
		post, err := newPost(path)
		if err != nil {
			return err
		}
		posts = append(posts, post)
	}
	sort.Sort(ByDateDesc(posts))
	err = runTasks(posts, t, destination)
	if err != nil {
		return err
	}
	fmt.Println("Finished generating Site...")
	return nil
}

func runTasks(posts []*Post, t *template.Template, destination string) error {
	var wg sync.WaitGroup
	finished := make(chan bool, 1)
	errors := make(chan error, 1)
	pool := make(chan struct{}, 50)
	generators := []Generator{}

	//posts
	for _, post := range posts {
		pg := PostGenerator{&PostConfig{
			Post:        post,
			Destination: destination,
			Template:    t,
		}}
		generators = append(generators, &pg)
	}
	tagPostsMap := createTagPostsMap(posts)

	// frontpage
	paging := config.SiteInfo.NumPostsFrontPage
	numOfPages := getNumberOfPages(posts)
	for i := 0; i < numOfPages; i++ {
		to := destination
		if i != 0 {
			to = fmt.Sprintf("%s/%d", destination, i+1)
		}
		toP := (i + 1) * paging
		if (i + 1) == numOfPages {
			toP = len(posts)
		}
		generators = append(generators, &ListingGenerator{&ListingConfig{
			Posts:       posts[i*paging : toP],
			Template:    t,
			Destination: to,
			PageTitle:   "",
			PageNum:     i + 1,
			MaxPageNum:  numOfPages,
		}})
	}

	// archive
	ag := ListingGenerator{&ListingConfig{
		Posts:       posts,
		Template:    t,
		Destination: fmt.Sprintf("%s/archive", destination),
		PageTitle:   "Archive",
	}}
	// tags
	tg := TagsGenerator{&TagsConfig{
		TagPostsMap: tagPostsMap,
		Template:    t,
		Destination: destination,
	}}
	// categories
	catPostsMap := createCatPostsMap(posts)
	ct := CategoriesGenerator{&CategoriesConfig{
		CatPostsMap: catPostsMap,
		Template:    t,
		Destination: destination,
	}}

	// sitemap
	sg := SitemapGenerator{&SitemapConfig{
		Posts:            posts,
		TagPostsMap:      tagPostsMap,
		CategoryPostsMap: catPostsMap,
		Destination:      destination,
	}}
	// rss
	rg := RSSGenerator{&RSSConfig{
		Posts:       posts,
		Destination: destination,
	}}
	// statics
	fileToDestination := map[string]string{
		config.SiteInfo.ThemeFolder + "favicon.ico": fmt.Sprintf("%s/favicon.ico", destination),
		config.SiteInfo.ThemeFolder + "robots.txt":  fmt.Sprintf("%s/robots.txt", destination),
		config.SiteInfo.ThemeFolder + "about.png":   fmt.Sprintf("%s/about.png", destination),
	}
	templateToFile := map[string]string{
		config.SiteInfo.ThemeFolder + "about.html": fmt.Sprintf("%s/about/index.html", destination),
	}
	statg := StaticsGenerator{&StaticsConfig{
		FileToDestination: fileToDestination,
		TemplateToFile:    templateToFile,
		Template:          t,
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

func writeIndexHTML(path, pageTitle string, content template.HTML, t *template.Template) error {
	return writeIndexHTMLPlus(path, pageTitle, content, t, false, 0, 0)
}

func writeIndexHTMLPost(path, pageTitle string, content template.HTML, t *template.Template, isPost bool) error {
	return writeIndexHTMLPlus(path, pageTitle, content, t, isPost, 0, 0)
}

func writeIndexHTMLPlus(path, pageTitle string, content template.HTML, t *template.Template, isPost bool, page, maxPage int) error {
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
		Name:          config.SiteInfo.Author,
		Year:          time.Now().Year(),
		HTMLTitle:     getHTMLTitle(pageTitle),
		PageTitle:     pageTitle,
		Content:       content,
		CanonicalLink: buildCanonicalLink(path, config.SiteInfo.BlogURL),
		PageNum:       page,
		NextPageNum:   next,
		PrevPageNum:   prev,
		URL:           buildCanonicalLink(path, config.SiteInfo.BlogURL),
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

func copyAdditionalArtifacts(path, postName string) error {
	src := config.SiteInfo.TempFolder + postName + "/artifacts/"
	return copyDir(src, path)
}

func getHTMLTitle(pageTitle string) string {
	if pageTitle == "" {
		return config.SiteInfo.BlogTitle
	}
	return fmt.Sprintf("%s - %s", pageTitle, config.SiteInfo.BlogTitle)
}

func createTagPostsMap(posts []*Post) map[string][]*Post {
	result := make(map[string][]*Post)
	for _, post := range posts {
		for _, tag := range post.Meta.Tags {
			key := strings.ToLower(tag)
			if result[key] == nil {
				result[key] = []*Post{post}
			} else {
				result[key] = append(result[key], post)
			}
		}
	}
	return result
}

func createCatPostsMap(posts []*Post) map[string][]*Post {
	result := make(map[string][]*Post)
	for _, post := range posts {
		for _, cat := range post.Meta.Categories {
			key := strings.ToLower(cat)
			if result[key] == nil {
				result[key] = []*Post{post}
			} else {
				result[key] = append(result[key], post)
			}
		}
	}
	return result
}

func getNumOfPagesOnFrontpage(posts []*Post) int {
	if len(posts) < config.SiteInfo.NumPostsFrontPage {
		return len(posts)
	}
	return config.SiteInfo.NumPostsFrontPage
}

func getNumberOfPages(posts []*Post) int {
	res := float64(len(posts)) / float64(config.SiteInfo.NumPostsFrontPage)
	r, _ := strconv.Atoi(fmt.Sprintf("%.0f", math.Ceil(res)))
	return r
}
