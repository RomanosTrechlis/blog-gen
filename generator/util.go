package generator

import (
	"fmt"
	"html/template"
	"os"
	"strings"
)

func clearAndCreateDestination(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("error removing folder at destination %s: %v ", path, err)
		}
	}
	return os.Mkdir(path, os.ModePerm)
}

func getTemplate(path string) (*template.Template, error) {
	t, err := template.ParseFiles(path)
	if err != nil {
		return nil, fmt.Errorf("error reading template %s: %v", path, err)
	}
	return t, nil
}

func buildCanonicalLink(path, baseURL string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 2 {
		return fmt.Sprintf("%s/%s/index.html", baseURL, strings.Join(parts[2:], "/"))
	}
	return "/"
}

func getTagLink(tag string) string {
	return fmt.Sprintf("/tags/%s/", strings.ToLower(tag))
}

func getFolder(path string) string {
	return path[:strings.LastIndex(path, "/")]
}
