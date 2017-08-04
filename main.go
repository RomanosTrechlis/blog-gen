package main

import (
	"flag"
	"github.com/RomanosTrechlis/blog-generator/cli"
	"log"
	"net/http"
	"strconv"
)

var (
	generate bool
	server bool
	download bool
	upload bool
	port int
	configFile string
)

func init() {
	flag.StringVar(&configFile, "configfile", "config.json", "is the file containing site's information")
	flag.BoolVar(&server, "run", false, "runs a simple server")
	flag.IntVar(&port, "port", 3000, "port of server")
	flag.BoolVar(&generate, "generate", false, "generates site content")
	flag.BoolVar(&download, "fetch", false, "fetches site content")
	flag.BoolVar(&upload, "upload", false, "uploads site content")
}

func main() {
	flag.Parse()
	siteInfo := cli.ReadConfig(configFile)

	if download {
		cli.Download(siteInfo.DataSource.Type, siteInfo.DataSource.Repository, siteInfo.TempFolder)
		cli.Download(siteInfo.Theme.Type, siteInfo.Theme.Repository, siteInfo.ThemeFolder)
	}

	if generate {
		cli.Generate(&siteInfo)
	}

	if upload {
		cli.Upload(&siteInfo)
	}

	if server {
		fs := http.FileServer(http.Dir(siteInfo.DestFolder))
		http.Handle("/", fs)

		serverPort := strconv.Itoa(port)
		log.Println("Listening @ localhost:" + serverPort + "/")
		http.ListenAndServe(":"+serverPort, nil)
	}
}
