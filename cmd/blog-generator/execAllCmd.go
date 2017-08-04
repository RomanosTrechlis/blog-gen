package main

import (
	"flag"

	"net/http"
	"strconv"

	"github.com/RomanosTrechlis/blog-generator/config"
	"github.com/RomanosTrechlis/blog-generator/datasource"
	"github.com/RomanosTrechlis/blog-generator/generator"
	"github.com/RomanosTrechlis/blog-generator/util/fs"
)

const execAllShortHelp = `Sequential execution of all commands`
const execAllLongHelp = `Sequential execution of all commands.

Downloads posts, and theme, then generates the blog and runs a web server.
`

type execAllCmd struct {
	siteInfo *config.SiteInformation
	port     int
}

func (cmd *execAllCmd) Name() string      { return "exec-all" }
func (cmd *execAllCmd) Args() string      { return "" }
func (cmd *execAllCmd) ShortHelp() string { return execAllShortHelp }
func (cmd *execAllCmd) LongHelp() string  { return execAllLongHelp }
func (cmd *execAllCmd) Hidden() bool      { return false }

func (cmd *execAllCmd) Register(fs *flag.FlagSet) {
	fs.IntVar(&cmd.port, "p", 8080, "port for web server")
}

func (cmd *execAllCmd) Run(ctx *ctx, args []string) (err error) {
	// download posts
	ds, err := datasource.New(cmd.siteInfo.DataSource.Type)
	if err != nil {
		ctx.Err.Fatal("please provide a datasource in the configuration file:", err)
		return
	}

	_, err = ds.Fetch(cmd.siteInfo.DataSource.Repository, cmd.siteInfo.DestFolder)
	if err != nil {
		ctx.Err.Fatal("failure to fetch posts:", err)
		return
	}

	// download theme
	ds, err = datasource.New(cmd.siteInfo.Theme.Type)
	if err != nil {
		ctx.Err.Fatal("please provide a datasource in the configuration file:", err)
		return
	}

	_, err = ds.Fetch(cmd.siteInfo.Theme.Repository, cmd.siteInfo.ThemeFolder)
	if err != nil {
		ctx.Err.Fatal("failure to fetch posts:", err)
		return
	}

	// generate blog
	dirs, err := fs.GetContentFolders(cmd.siteInfo.TempFolder)
	if err != nil {
		return
	}
	g := generator.NewSiteGenerator(&generator.SiteConfig{
		Sources:  dirs,
		SiteInfo: cmd.siteInfo,
	})
	err = g.Generate()
	if err != nil {
		ctx.Err.Fatal("failed to generate blog:", err)
		return
	}

	// run web server
	blog := http.FileServer(http.Dir(cmd.siteInfo.DestFolder))
	http.Handle("/", blog)

	serverPort := strconv.Itoa(cmd.port)
	ctx.Out.Println("Listening @ localhost:" + serverPort + "/")
	err = http.ListenAndServe(":"+serverPort, nil)
	ctx.Err.Fatal(err)
	return err
}
