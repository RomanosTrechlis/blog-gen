package cli

import (
	"log"
	"github.com/RomanosTrechlis/blog-generator/config"
	"github.com/RomanosTrechlis/blog-generator/datasource"
	"github.com/RomanosTrechlis/blog-generator/util/fs"
	"github.com/RomanosTrechlis/blog-generator/generator"
	"github.com/RomanosTrechlis/blog-generator/endpoint"
)

// ReadConfig creates object holding site information
func ReadConfig() {
	config.SiteInfo = config.NewSiteInformation()
}

// Download fetches content from the data source
func Download() {
	if config.ConfigFile == "" {
		log.Fatal("please provide a configuration file -configfile flag")
	}

	// handle blog posts repository
	var err error
	switch config.SiteInfo.DataSource.Type {
	case "git":
		ds := datasource.NewGitDataSource()
		_, err = ds.Fetch(config.SiteInfo.DataSource.Repository,
			config.SiteInfo.TempFolder)
	case "local":
		ds := datasource.NewLocalDataSource()
		_, err = ds.Fetch(config.SiteInfo.DataSource.Repository,
			config.SiteInfo.TempFolder)
	case "":
		log.Fatal("please provide a datasource in the configuration file")
	}
	if err != nil {
		log.Fatal(err)
	}

	// handle theme repository
	switch config.SiteInfo.Theme.Type {
	case "git":
		ds := datasource.NewGitDataSource()
		_, err = ds.Fetch(config.SiteInfo.Theme.Repository,
			config.SiteInfo.ThemeFolder)
	case "local":
		ds := datasource.NewLocalDataSource()
		_, err = ds.Fetch(config.SiteInfo.Theme.Repository,
			config.SiteInfo.ThemeFolder)
	case "":
		log.Fatal("please provide a datasource in the configuration file")
	}
	if err != nil {
		log.Fatal(err)
	}


}

// Generate creates site's content
func Generate() {
	dirs, err := fs.GetContentFolders(config.SiteInfo.TempFolder)
	if err != nil {
		log.Fatal(err)
	}
	g := generator.NewSiteGenerator(&generator.SiteConfig{
		Sources:     dirs,
		Destination: config.SiteInfo.DestFolder,
	})

	err = g.Generate()
	if err != nil {
		log.Fatal(err)
	}
}

// Upload uploads content to endpoint
func Upload() {
	e := endpoint.NewGitEndpoint()
	err := e.Upload(config.SiteInfo.Upload.URL)
	if err != nil {
		log.Fatal(err)
	}
}
