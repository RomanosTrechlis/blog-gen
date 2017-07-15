package cli

import (
	"github.com/RomanosTrechlis/blog-generator/config"
	"github.com/RomanosTrechlis/blog-generator/datasource"
	"github.com/RomanosTrechlis/blog-generator/endpoint"
	"github.com/RomanosTrechlis/blog-generator/generator"
	"github.com/RomanosTrechlis/blog-generator/util/fs"
	"log"
)

// ReadConfig creates object holding site information
func ReadConfig(configFile string) (siteInfo *config.SiteInformation) {
	config.SiteInfo = config.NewSiteInformation(configFile)
	return &config.SiteInfo
}

// Download fetches content from the data source
func Download(dsType, dsRepository, tempFolder string) {
	// handle blog posts repository
	var err error
	switch dsType {
	case "git":
		ds := datasource.NewGitDataSource()
		_, err = ds.Fetch(dsRepository, tempFolder)
	case "local":
		ds := datasource.NewLocalDataSource()
		_, err = ds.Fetch(dsRepository, tempFolder)
	case "":
		log.Fatal("please provide a datasource in the configuration file")
	}
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
	e := endpoint.NewGitEndpoint()
	err := e.Upload(siteInfo.Upload.URL)
	if err != nil {
		log.Fatal(err)
	}
}
