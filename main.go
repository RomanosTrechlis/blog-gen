package main

import (
	"flag"
	"github.com/RomanosTrechlis/blog-generator/cli"
	"github.com/RomanosTrechlis/blog-generator/config"
	"log"
	"net/http"
	"strconv"
)

var server bool
var generate bool
var download bool
var upload bool
var port int

func init() {
	flag.StringVar(&config.ConfigFile, "configfile", "config.json", "is the file containing site's information")
	flag.BoolVar(&server, "run", false, "runs a simple server")
	flag.IntVar(&port, "port", 3000, "port of server")
	flag.BoolVar(&generate, "generate", false, "generates site content")
	flag.BoolVar(&download, "fetch", false, "fetches site content")
	flag.BoolVar(&upload, "upload", false, "uploads site content")
}

func main() {
	flag.Parse()
	cli.ReadConfig(config.ConfigFile)

	if download {
		cli.Download()
	}

	if generate {
		cli.Generate()
	}

	if upload {
		cli.Upload()
	}

	if server {
		fs := http.FileServer(http.Dir(config.SiteInfo.DestFolder))
		http.Handle("/", fs)

		serverPort := strconv.Itoa(port)
		log.Println("Listening @ localhost:" + serverPort + "/")
		http.ListenAndServe(":"+serverPort, nil)
	}
}
