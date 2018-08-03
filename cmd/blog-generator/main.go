package main

import (
	"os"
	"github.com/RomanosTrechlis/go-icls/cli"
	"github.com/RomanosTrechlis/blog-generator/config"
	"fmt"
)

func createCommandTree(siteInfo config.SiteInformation) *cli.CLI {
	c := cli.New()
	c.New("posts", getPostsShortHelp, getPostsLongHelp, getPostHandler(siteInfo))
	c.New("theme", getThemeShortHelp, getThemeLongHelp, getThemeHandler(siteInfo))
	c.New("generate", generateShortHelp, generateLongHelp, getGenerateHandler(siteInfo))
	c.New("example", jsonExampleShortHelp, jsonExampleLongHelp, getExampleConfigHandler(siteInfo))
	server := c.New("server", jsonExampleShortHelp, jsonExampleLongHelp, getServerHandler(c, siteInfo))
	server.IntFlag("p", "port", 8080, "port for web server", false)
	all := c.New("all", jsonExampleShortHelp, jsonExampleLongHelp, getExecAllHandler(c, siteInfo))
	all.IntFlag("p", "port", 8080, "port for web server", false)
	return c
}

func main() {
	args := os.Args[1:]
	line := ""
	siteInfo, err := config.New("config.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "config.json reading error: %v\n", err)
		args = append(args, "-h")
	}

	c := createCommandTree(siteInfo)

	if len(args) == 0{
		args = append(args, "-h")
	}

	for _, s := range args {
		line += s + " "
	}

	_, err = c.Execute(line)
	if err != nil {
		fmt.Fprintf(os.Stderr, "command exited with error: %v\n", err)
		os.Exit(1)
	}
}
