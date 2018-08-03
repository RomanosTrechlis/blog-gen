package cli

import (
	"log"

	"github.com/RomanosTrechlis/blog-generator/config"
	"github.com/RomanosTrechlis/blog-generator/datasource"
	"github.com/RomanosTrechlis/blog-generator/endpoint"
	"github.com/RomanosTrechlis/blog-generator/generator"
	"github.com/RomanosTrechlis/blog-generator/util/fs"
)

// ReadConfig creates object holding site information
func ReadConfig(configFile string) (siteInfo config.SiteInformation) {
	s, _ := config.New(configFile)
	return s
}

// Download fetches content from the data source
func Download(dsType, dsRepository, tempFolder string) {
	// handle blog posts repository
	ds, err := datasource.New(dsType)
	if err != nil {
		log.Fatal(err)
	}

	_, err = ds.Fetch(dsRepository, tempFolder)
	if err != nil {
		log.Fatal(err)
	}
}

// Generate creates site's content
func Generate(siteInfo *config.SiteInformation) {
	dirs, err := fs.GetContentFolders(siteInfo.TempFolder)
	if err != nil {
		log.Fatal(err)
	}
	g := generator.NewSiteGenerator(&generator.SiteConfig{
		Sources:  dirs,
		SiteInfo: siteInfo,
	})

	err = g.Generate()
	if err != nil {
		log.Fatal(err)
	}
}

// Upload uploads content to endpoint
func Upload(siteInfo *config.SiteInformation) {
	e, err := endpoint.New(siteInfo.Upload.Type)
	if err != nil {
		log.Fatal(err)
	}
	err = e.Upload(siteInfo.DestFolder, siteInfo.Upload.Username, siteInfo.Upload.Password, siteInfo.Upload.URL)
	if err != nil {
		log.Fatal(err)
	}
}
