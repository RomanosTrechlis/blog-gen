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
	Sources  []string
	SiteInfo config.SiteInformation
}

// New creates a new SiteGenerator
func NewSiteGenerator(config *SiteConfig) *SiteGenerator {
	return &SiteGenerator{Config: config}
}

var templatePath string

// Generate starts the static blog generation
func (g *SiteGenerator) Generate() error {
	templatePath = g.Config.SiteInfo.ThemeFolder + "template.html"
	fmt.Println("Generating Site...")
	sources := g.Config.Sources
	destination := g.Config.SiteInfo.DestFolder
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
		post, err := newPost(path, g.Config.SiteInfo.DateFormat)
		if err != nil {
			return err
		}
		posts = append(posts, post)
	}
	sort.Sort(ByDateDesc(posts))
	err = runTasks(posts, t, g.Config.SiteInfo)
	if err != nil {
		return err
	}
	fmt.Println("Finished generating Site...")
	return nil
}

func runTasks(posts []*Post, t *template.Template, siteInfo config.SiteInformation) error {
	var wg sync.WaitGroup
	finished := make(chan bool, 1)
	errors := make(chan error, 1)
	pool := make(chan struct{}, 50)
	generators := []Generator{}
	destination := siteInfo.DestFolder

	//posts
	for _, post := range posts {
		pg := PostGenerator{&PostConfig{
			Post:        post,
			Destination: destination,
			Template:    t,
			DateFormat:  siteInfo.DateFormat,
			TempFolder:  siteInfo.TempFolder,
			BlogTitle:   siteInfo.BlogTitle,
			Author:      siteInfo.Author,
			BlogURL:     siteInfo.BlogURL,
		}}
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
		generators = append(generators, &ListingGenerator{&ListingConfig{
			Posts:       posts[i*paging : toP],
			Template:    t,
			Destination: to,
			PageTitle:   "",
			PageNum:     i + 1,
			MaxPageNum:  numOfPages,
			ThemeFolder: siteInfo.ThemeFolder,
		}})
	}

	// archive
	ag := ListingGenerator{&ListingConfig{
		Posts:       posts,
		Template:    t,
		Destination: fmt.Sprintf("%s/archive", destination),
		PageTitle:   "Archive",
		ThemeFolder: siteInfo.ThemeFolder,
		BlogTitle:   siteInfo.BlogTitle,
		Author:      siteInfo.Author,
		BlogURL:     siteInfo.BlogURL,
	}}
	// tags
	tg := TagsGenerator{&TagsConfig{
		TagPostsMap: tagPostsMap,
		Template:    t,
		Destination: destination,
		BlogTitle:   siteInfo.BlogTitle,
		Author:      siteInfo.Author,
		BlogURL:     siteInfo.BlogURL,
		ThemeFolder: siteInfo.ThemeFolder,
	}}
	// categories
	catPostsMap := createCatPostsMap(posts)
	ct := CategoriesGenerator{&CategoriesConfig{
		CatPostsMap: catPostsMap,
		Template:    t,
		Destination: destination,
		BlogTitle:   siteInfo.BlogTitle,
		Author:      siteInfo.Author,
		BlogURL:     siteInfo.BlogURL,
		ThemeFolder: siteInfo.ThemeFolder,
	}}

	// sitemap
	sg := SitemapGenerator{&SitemapConfig{
		Posts:            posts,
		TagPostsMap:      tagPostsMap,
		CategoryPostsMap: catPostsMap,
		Destination:      destination,
		BlogURL:          siteInfo.BlogURL,
	}}
	// rss
	rg := RSSGenerator{&RSSConfig{
		Posts:           posts,
		Destination:     destination,
		BlogTitle:       siteInfo.BlogTitle,
		BlogURL:         siteInfo.BlogURL,
		BlogLanguage:    siteInfo.BlogLanguage,
		BlogDescription: siteInfo.BlogDescription,
		DateFormat:      siteInfo.DateFormat,
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

func writeIndexHTML(path, pageTitle, author, blogURL, blogTitle string, content template.HTML, t *template.Template) error {
	return writeIndexHTMLPlus(path, pageTitle, author, blogURL, blogTitle, content, t, false, 0, 0)
}

func writeIndexHTMLPost(path, pageTitle, author, blogURL, blogTitle string, content template.HTML, t *template.Template, isPost bool) error {
	return writeIndexHTMLPlus(path, pageTitle, author, blogURL, blogTitle, content, t, isPost, 0, 0)
}

func writeIndexHTMLPlus(path, pageTitle, author, blogURL, blogTitle string, content template.HTML, t *template.Template,
	isPost bool, page, maxPage int) error {
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
		Name:          author,
		Year:          time.Now().Year(),
		HTMLTitle:     getHTMLTitle(pageTitle, blogTitle),
		PageTitle:     pageTitle,
		Content:       content,
		CanonicalLink: buildCanonicalLink(path, blogURL),
		PageNum:       page,
		NextPageNum:   next,
		PrevPageNum:   prev,
		URL:           buildCanonicalLink(path, blogURL),
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

func copyAdditionalArtifacts(path, postName, tempFolder string) error {
	src := tempFolder + postName + "/artifacts/"
	return copyDir(src, path)
}

func getHTMLTitle(pageTitle, blogTitle string) string {
	if pageTitle == "" {
		return blogTitle
	}
	return fmt.Sprintf("%s - %s", pageTitle, blogTitle)
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

func getNumOfPagesOnFrontpage(posts []*Post, postsPerPage int) int {
	if len(posts) < postsPerPage {
		return len(posts)
	}
	return postsPerPage
}

func getNumberOfPages(posts []*Post, postsPerPage int) int {
	res := float64(len(posts)) / float64(postsPerPage)
	r, _ := strconv.Atoi(fmt.Sprintf("%.0f", math.Ceil(res)))
	return r
}
