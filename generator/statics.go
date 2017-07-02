package generator

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
	"github.com/RomanosTrechlis/blog-generator/util/fs"
)

// StaticsGenerator object
type StaticsGenerator struct {
	Config *StaticsConfig
}

// StaticsConfig holds the data for the static sites
type StaticsConfig struct {
	FileToDestination map[string]string
	TemplateToFile    map[string]string
	Template          *template.Template
}

// Generate creates the static pages
func (g *StaticsGenerator) Generate() error {
	fmt.Println("\tCopying Statics...")
	fileToDestination := g.Config.FileToDestination
	templateToFile := g.Config.TemplateToFile
	t := g.Config.Template
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
	for k, v := range templateToFile {
		err := createFolderIfNotExist(getFolder(v))
		if err != nil {
			return err
		}
		content, err := ioutil.ReadFile(k)
		if err != nil {
			return fmt.Errorf("error reading file %s: %v", k, err)
		}
		err = writeIndexHTML(getFolder(v), getTitle(k), template.HTML(content), t)
		if err != nil {
			return err
		}
	}
	fmt.Println("\tFinished copying statics...")
	return nil
}

func createFolderIfNotExist(path string) error {
	_, err := os.Stat(path)
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

func getTitle(path string) string {
	fileName := path[strings.LastIndex(path, "/")+1 : strings.LastIndex(path, ".")]
	return fmt.Sprintf("%s%s", strings.ToUpper(string(fileName[0])), fileName[1:])
}
