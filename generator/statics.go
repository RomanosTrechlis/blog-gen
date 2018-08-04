package generator

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/RomanosTrechlis/blog-generator/config"
	"github.com/RomanosTrechlis/blog-generator/util/fs"
)

// staticsGenerator object
type staticsGenerator struct {
	fileToDestination map[string]string
	templateToFile    map[string]string
	template          *template.Template
	siteInfo          *config.SiteInformation
}

// Generate creates the static pages
func (g *staticsGenerator) Generate() (err error) {
	fmt.Println("\tCopying Statics...")
	fileToDestination := g.fileToDestination
	templateToFile := g.templateToFile
	t := g.template
	if len(fileToDestination) > 0 {
		for k, v := range fileToDestination {
			err := createFolderIfNotExist(getFolder(v))
			if err != nil {
				return err
			}
			err = fs.CopyFile(k, v)
			if err != nil {
				return err
			}
		}
	}
	if len(templateToFile) > 0 {
		for k, v := range templateToFile {
			err := createFolderIfNotExist(getFolder(v))
			if err != nil {
				return err
			}
			content, err := ioutil.ReadFile(k)
			if err != nil {
				return fmt.Errorf("error reading file %s: %v", k, err)
			}

			c := htmlConfig{
				path: getFolder(v),
				pageTitle: getTitle(k),
				pageNum: 0,
				maxPageNum: 0,
				isPost: false,
				temp: t,
				content: template.HTML(content),
				siteInfo: g.siteInfo,
			}
			err = c.writeHTML()
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("\tFinished copying statics...")
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
