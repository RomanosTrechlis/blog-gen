package generator

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/RomanosTrechlis/blog-gen/config"
	"github.com/RomanosTrechlis/blog-gen/util/fs"
)

// staticsGenerator object
type staticsGenerator struct {
	fileToDestination map[string]string
	templateToFile    map[string]string
	template          *template.Template
	siteInfo          *config.SiteInformation
}

// Generate creates the static pages
func (g *staticsGenerator) Generate() error {
	fmt.Println("\tCopying Statics...")

	err := g.resolveFileToDestination()
	if err != nil {
		return err
	}

	err = g.resolveTemplateToFile()
	if err != nil {
		return err
	}

	fmt.Println("\tFinished copying statics...")
	return nil
}

func (g *staticsGenerator) resolveFileToDestination() error {
	if len(g.fileToDestination) == 0 {
		return nil
	}

	for k, v := range g.fileToDestination {
		folder := fs.GetFolderNameFrom(v)
		if folder != "" {
			err := fs.CreateFolderIfNotExist(folder)
			if err != nil {
				return err
			}
		}

		err := fs.CopyFile(k, fs.GetFolderNameFrom(v))
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *staticsGenerator) resolveTemplateToFile() error {
	if len(g.templateToFile) == 0 {
		return nil
	}
	for k, v := range g.templateToFile {
		folder := fs.GetFolderNameFrom(v)
		if folder != "" {
			err := fs.CreateFolderIfNotExist(fs.GetFolderNameFrom(v))
			if err != nil {
				return err
			}
		}

		content, err := ioutil.ReadFile(k)
		if err != nil {
			return fmt.Errorf("error reading file %s: %v", k, err)
		}

		c := htmlConfig{
			path:       fs.GetFolderNameFrom(v),
			pageTitle:  getTitle(k),
			pageNum:    0,
			maxPageNum: 0,
			isPost:     false,
			temp:       g.template,
			content:    template.HTML(content),
			siteInfo:   g.siteInfo,
		}
		err = c.writeHTML()
		if err != nil {
			return err
		}
	}
	return nil
}

func createFolderIfNotExist(path string) (err error) {
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(path, os.ModePerm)
			if err != nil {
				return fmt.Errorf("error creating directory %s: %v", path, err)
			}
		} else {
			return fmt.Errorf("error accessing directory %s: %v", path, err)
		}
	}
	return nil
}

func getTitle(path string) (title string) {
	fileName := path[strings.LastIndex(path, "/")+1 : strings.LastIndex(path, ".")]
	title = fmt.Sprintf("%s%s", strings.ToUpper(string(fileName[0])), fileName[1:])
	return title
}
