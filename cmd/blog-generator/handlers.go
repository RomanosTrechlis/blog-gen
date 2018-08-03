package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/RomanosTrechlis/blog-generator/config"
	"github.com/RomanosTrechlis/blog-generator/datasource"
	"github.com/RomanosTrechlis/blog-generator/generator"
	"github.com/RomanosTrechlis/blog-generator/util/fs"
	"github.com/RomanosTrechlis/go-icls/cli"
)

func getPostHandler(siteInfo config.SiteInformation) func(flags map[string]string) error {
	return func(flags map[string]string) error {
		ds, err := datasource.New(siteInfo.DataSource.Type)
		if err != nil {
			return fmt.Errorf("please provide a datasource in the configuration file: %v", err)
		}

		_, err = ds.Fetch(siteInfo.DataSource.Repository, siteInfo.TempFolder)
		if err != nil {
			return fmt.Errorf("failure to fetch posts: %v", err)
		}
		return nil
	}
}

func getThemeHandler(siteInfo config.SiteInformation) func(flags map[string]string) error {
	return func(flags map[string]string) error {
		ds, err := datasource.New(siteInfo.Theme.Type)
		if err != nil {
			return fmt.Errorf("please provide a datasource in the configuration file: %v", err)
		}

		_, err = ds.Fetch(siteInfo.Theme.Repository, siteInfo.ThemeFolder)
		if err != nil {
			return fmt.Errorf("failure to fetch theme: %v", err)
		}
		return nil
	}
}

func getGenerateHandler(siteInfo config.SiteInformation) func(flags map[string]string) error {
	return func(flags map[string]string) error {
		dirs, err := fs.GetContentFolders(siteInfo.TempFolder)
		if err != nil {
			return fmt.Errorf("failed to get contentes from %s: %v", siteInfo.TempFolder, err)
		}
		g := generator.NewSiteGenerator(&generator.SiteConfig{
			Sources:  dirs,
			SiteInfo: &siteInfo,
		})

		err = g.Generate()
		return fmt.Errorf("failed to generate blog: %v", err)
	}
}

func getExampleConfigHandler(siteInfo config.SiteInformation) func(flags map[string]string) error {
	return func(flags map[string]string) error {
		fmt.Fprint(os.Stdout, jsonExampleLongHelp)
		return nil
	}
}

func getServerHandler(c *cli.CLI, siteInfo config.SiteInformation) func(flags map[string]string) error {
	return func(flags map[string]string) error {
		fs := http.FileServer(http.Dir(siteInfo.DestFolder))
		http.Handle("/", fs)

		serverPort, err := c.IntValue("p", "server", flags)
		if err != nil {
			return fmt.Errorf("server port is not correct: %v", err)
		}

		fmt.Fprintf(os.Stdout, "Listening @ localhost: %d/\n", serverPort)
		http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil)
		return nil
	}
}

func getExecAllHandler(c *cli.CLI, siteInfo config.SiteInformation) func(flags map[string]string) error {
	return func(flags map[string]string) error {
		// download posts
		ds, err := datasource.New(siteInfo.DataSource.Type)
		if err != nil {
			return fmt.Errorf("please provide a datasource in the configuration file: %v", err)
		}

		_, err = ds.Fetch(siteInfo.DataSource.Repository, siteInfo.TempFolder)
		if err != nil {
			return fmt.Errorf("failure to fetch posts: %v", err)
		}

		// download theme
		ds, err = datasource.New(siteInfo.Theme.Type)
		if err != nil {
			return fmt.Errorf("please provide a datasource in the configuration file: %v", err)
		}

		_, err = ds.Fetch(siteInfo.Theme.Repository, siteInfo.ThemeFolder)
		if err != nil {
			return fmt.Errorf("failure to fetch posts: %v", err)
		}

		// generate blog
		dirs, err := fs.GetContentFolders(siteInfo.TempFolder)
		if err != nil {
			return fmt.Errorf("failed to get contentes from %s: %v", siteInfo.TempFolder, err)
		}
		g := generator.NewSiteGenerator(&generator.SiteConfig{
			Sources:  dirs,
			SiteInfo: &siteInfo,
		})
		err = g.Generate()
		if err != nil {
			return fmt.Errorf("failed to generate blog: %v", err)
		}

		// run web server
		blog := http.FileServer(http.Dir(siteInfo.DestFolder))
		http.Handle("/", blog)

		serverPort, err := c.IntValue("p", "all", flags)
		if err != nil {
			return fmt.Errorf("server port is not correct: %v", err)
		}

		fmt.Fprintf(os.Stdout, "Listening @ localhost: %d/\n", serverPort)
		http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil)
		return nil
	}
}
