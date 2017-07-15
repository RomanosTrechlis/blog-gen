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

// DownloadPosts fetches content from the data source
func DownloadPosts(siteInfo *config.SiteInformation) {
	// handle blog posts repository
	var err error
	switch siteInfo.DataSource.Type {
	case "git":
		ds := datasource.NewGitDataSource()
		_, err = ds.Fetch(siteInfo.DataSource.Repository,
			siteInfo.TempFolder)
	case "local":
		ds := datasource.NewLocalDataSource()
		_, err = ds.Fetch(siteInfo.DataSource.Repository,
			siteInfo.TempFolder)
	case "":
		log.Fatal("please provide a datasource in the configuration file")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func DownloadTheme(siteInfo *config.SiteInformation) {
	var err error
	// handle theme repository
	switch siteInfo.Theme.Type {
	case "git":
		ds := datasource.NewGitDataSource()
		_, err = ds.Fetch(siteInfo.Theme.Repository,
			siteInfo.ThemeFolder)
	case "local":
		ds := datasource.NewLocalDataSource()
		_, err = ds.Fetch(siteInfo.Theme.Repository,
			siteInfo.ThemeFolder)
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
